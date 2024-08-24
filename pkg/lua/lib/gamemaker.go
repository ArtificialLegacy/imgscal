package lib

import (
	"fmt"
	"image"
	"sync"

	"github.com/ArtificialLegacy/gm-proj-tool/yyp"
	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	imageutil "github.com/ArtificialLegacy/imgscal/pkg/image_util"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	golua "github.com/yuin/gopher-lua"
)

const LIB_GAMEMAKER = "gamemaker"

/// @lib Gamemaker
/// @import gamemaker
/// @desc
/// Library for working with Gamemaker projects.

func RegisterGamemaker(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_GAMEMAKER, r, r.State, lg)

	/// @func project_load(path) -> int<collection.CRATE_GAMEMAKER>
	/// @arg path {string} -> Path to the directory containing a Gamemaker project.
	/// @returns {int<collection.CRATE_GAMEMAKER>}
	lib.CreateFunction(tab, "project_load",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			proj, err := yyp.NewProject(args["path"].(string))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to load project: %s", err), log.LEVEL_ERROR)), 0)
			}

			id := r.CR_GMP.Add(proj)
			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func project_save(id)
	/// @arg id {int<collection.CRATE_GAMEMAKER>} - ID for the loaded Gamemaker project.
	lib.CreateFunction(tab, "project_save",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			proj, err := r.CR_GMP.Item(args["id"].(int))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to find project: %d, %s", args["id"], err), log.LEVEL_ERROR)), 0)
			}

			err = proj.DataSave()
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to save project data: %d, %s", args["id"], err), log.LEVEL_ERROR)), 0)
			}

			return 0
		})

	/// @func sprite(name, width, height) -> struct<gamemaker.Sprite>
	/// @arg name {string} - Name of the sprite asset.
	/// @arg width {int}
	/// @arg height {int}
	/// @returns {struct<gamemaker.Sprite>}
	lib.CreateFunction(tab, "sprite",
		[]lua.Arg{
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "width"},
			{Type: lua.INT, Name: "height"},
			{Type: lua.ANY, Name: "parent"},
			{Type: lua.ANY, Name: "texgroup"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := spriteTable(state,
				args["name"].(string),
				args["width"].(int), args["height"].(int),
				args["parent"].(*golua.LTable),
				args["texgroup"].(*golua.LTable),
				lg,
			)

			state.Push(t)
			return 1
		})

	/// @func sprite_save(id, sprite)
	/// @arg id {int<collection.CRATE_GAMEMAKER>} - ID for the loaded Gamemaker project.
	/// @arg sprite {struct<gamemaker.Sprite>}
	lib.CreateFunction(tab, "sprite_save",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.ANY, Name: "sprite"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			proj, err := r.CR_GMP.Item(args["id"].(int))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to find project: %d, %s", args["id"], err), log.LEVEL_ERROR)), 0)
			}

			sprite, err := spriteBuild(args["sprite"].(*golua.LTable), r)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to build sprite: %s", err), log.LEVEL_ERROR)), 0)
			}

			err = proj.ImportResource(sprite)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to import sprite: %s", err), log.LEVEL_ERROR)), 0)
			}

			return 0
		})

	/// @func project_as_parent(id) -> struct<gamemaker.ResourceNode>
	/// @arg id {int<collection.CRATE_GAMEMAKER>} - ID for the loaded Gamemaker project.
	/// @returns {struct<gamemaker.ResourceNode>}
	lib.CreateFunction(tab, "project_as_parent",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			proj, err := r.CR_GMP.Item(args["id"].(int))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to find project: %d, %s", args["id"], err), log.LEVEL_ERROR)), 0)
			}

			node := proj.AsParent()
			t := resourceNodeTable(state, node.Name, node.Path)

			state.Push(t)
			return 1
		})

	/// @func texgroup_default() -> struct<gamemaker.ResourceNode>
	/// @returns {struct<gamemaker.ResourceNode>}
	lib.CreateFunction(tab, "texgroup_default",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			node := yyp.ProjectTextureGroupDefaultID()
			t := resourceNodeTable(state, node.Name, node.Path)

			state.Push(t)
			return 1
		})
}

const (
	LAYER_IMAGE  string = "image"
	LAYER_FOLDER string = "folder"
)

func assignSpriteLayers(state *golua.LState, t *golua.LTable, layers *golua.LTable) {
	t.RawSetString("__layers", layers)

	tableBuilderFunc(state, t, "image", func(state *golua.LState, t *golua.LTable) {
		name := state.CheckString(-1)

		lt := t.RawGetString("__layers").(*golua.LTable)
		lt.Append(layerImageTable(state, name))

		base := t.RawGetString("__base").(*golua.LTable)
		count := base.RawGetString("__layerCount").(golua.LNumber)
		base.RawSetString("__layerCount", count+1)
	})

	tableBuilderFunc(state, t, "default", func(state *golua.LState, t *golua.LTable) {
		lt := t.RawGetString("__layers").(*golua.LTable)
		lt.Append(layerImageTable(state, yyp.SPRITELAYER_DEFAULTNAME))

		base := t.RawGetString("__base").(*golua.LTable)
		count := base.RawGetString("__layerCount").(golua.LNumber)
		base.RawSetString("__layerCount", count+1)
	})

	t.RawSetString("folder", state.NewFunction(func(state *golua.LState) int {
		t := state.CheckTable(-2)
		name := state.CheckString(-1)

		lt := t.RawGetString("__layers").(*golua.LTable)
		ft := layerFolderTable(state, name, t)
		lt.Append(ft)

		state.Push(ft)
		return 1
	}))
}

func assignSpriteNesting(state *golua.LState, t *golua.LTable, parent *golua.LTable, base *golua.LTable) {
	t.RawSetString("__parent", parent)
	t.RawSetString("__base", base)

	t.RawSetString("back", state.NewFunction(func(state *golua.LState) int {
		t := state.CheckTable(-1)

		parent := t.RawGetString("__parent").(*golua.LTable)
		state.Push(parent)
		return 1
	}))
}

func resourceNodeTable(state *golua.LState, name, pth string) *golua.LTable {
	/// @struct ResourceNode
	/// @prop name {string}
	/// @prop filepath {string}

	t := state.NewTable()
	t.RawSetString("name", golua.LString(name))
	t.RawSetString("filepath", golua.LString(pth))

	return t
}

func resourceNodeBuild(t *golua.LTable) yyp.ProjectResourceNode {
	name := string(t.RawGetString("name").(golua.LString))
	filepath := string(t.RawGetString("filepath").(golua.LString))

	return yyp.ProjectResourceNode{
		Name: name,
		Path: filepath,
	}
}

func spriteTable(state *golua.LState, name string, width, height int, parent, texgroup *golua.LTable, lg *log.Logger) *golua.LTable {
	t := state.NewTable()
	t.RawSetString("name", golua.LString(name))
	t.RawSetString("width", golua.LNumber(width))
	t.RawSetString("height", golua.LNumber(height))
	t.RawSetString("parent", parent)
	t.RawSetString("texgroup", texgroup)

	t.RawSetString("__layerCount", golua.LNumber(0))
	t.RawSetString("__layers", state.NewTable())

	t.RawSetString("__frames", state.NewTable())

	t.RawSetString("layers", state.NewFunction(func(state *golua.LState) int {
		t := state.CheckTable(-1)

		state.Push(layersTable(state, t))
		return 1
	}))

	t.RawSetString("frames", state.NewFunction(func(state *golua.LState) int {
		t := state.CheckTable(-1)

		state.Push(framesTable(state, t, lg))
		return 1
	}))

	return t
}

func spriteBuild(t *golua.LTable, r *lua.Runner) (*yyp.Sprite, error) {
	name := string(t.RawGetString("name").(golua.LString))
	width := int(t.RawGetString("width").(golua.LNumber))
	height := int(t.RawGetString("height").(golua.LNumber))
	parent := resourceNodeBuild(t.RawGetString("parent").(*golua.LTable))
	texgroup := resourceNodeBuild(t.RawGetString("texgroup").(*golua.LTable))

	frames := [][]*image.NRGBA{}
	frameList := t.RawGetString("__frames").(*golua.LTable)
	if frameList.Len() == 0 {
		return nil, fmt.Errorf("sprite must have at least 1 frame to import")
	}

	layerCount := int(t.RawGetString("__layerCount").(golua.LNumber))

	wg := sync.WaitGroup{}
	wg.Add(frameList.Len() * layerCount)

	for ind := range frameList.Len() {
		frame := frameList.RawGetInt(ind + 1).(*golua.LTable)
		fLayers := frameBuild(frame)

		frames = append(frames, make([]*image.NRGBA, len(fLayers)))

		for id, img := range fLayers {
			r.IC.Schedule(img, &collection.Task[collection.ItemImage]{
				Lib:  LIB_GAMEMAKER,
				Name: "sprite_save",
				Fn: func(i *collection.Item[collection.ItemImage]) {
					frames[ind][id] = imageutil.CopyImage(i.Self.Image, imageutil.MODEL_NRGBA).(*image.NRGBA)
					wg.Done()
				},
				Fail: func(i *collection.Item[collection.ItemImage]) {
					wg.Done()
				},
			})
		}
	}

	wg.Wait()

	layers := []yyp.SpriteLayer{}
	layerIndex := 0
	layersBuild(t, &layerIndex, &layers, &frames)

	return yyp.NewSprite(name, parent, texgroup, width, height, layers)
}

func layersBuild(t *golua.LTable, layerIndex *int, layers *[]yyp.SpriteLayer, frames *[][]*image.NRGBA) {
	layerList := t.RawGetString("__layers").(*golua.LTable)

	for i := range layerList.Len() {
		layer := layerList.RawGetInt(i + 1).(*golua.LTable)

		typ := layer.RawGetString("type").(golua.LString)
		name := layer.RawGetString("name").(golua.LString)
		switch string(typ) {
		case LAYER_IMAGE:
			frameList := []*image.NRGBA{}
			for _, f := range *frames {
				frameList = append(frameList, f[*layerIndex])
			}

			lt := yyp.SpriteLayer{Name: string(name), Frames: frameList}
			*layers = append(*layers, lt)
			*layerIndex++

		case LAYER_FOLDER:
			lt := yyp.SpriteLayer{Name: string(name)}
			layersBuild(layer, layerIndex, &lt.Layers, frames)
			*layers = append(*layers, lt)
		}
	}
}

func layersTable(state *golua.LState, parent *golua.LTable) *golua.LTable {
	t := state.NewTable()
	layers := parent.RawGetString("__layers").(*golua.LTable)
	assignSpriteLayers(state, t, layers)
	assignSpriteNesting(state, t, parent, parent)
	return t
}

func layerImageTable(state *golua.LState, name string) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("type", golua.LString(LAYER_IMAGE))
	t.RawSetString("name", golua.LString(name))

	return t
}

func layerFolderTable(state *golua.LState, name string, parent *golua.LTable) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("type", golua.LString(LAYER_FOLDER))
	t.RawSetString("name", golua.LString(name))

	base := parent.RawGetString("__base").(*golua.LTable)

	assignSpriteLayers(state, t, state.NewTable())
	assignSpriteNesting(state, t, parent, base)

	return t
}

func framesTable(state *golua.LState, parent *golua.LTable, lg *log.Logger) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("__frames", parent.RawGetString("__frames"))

	layerCount := parent.RawGetString("__layerCount").(golua.LNumber)
	if int(layerCount) == 0 {
		state.Error(golua.LString(lg.Append("sprite must have at least 1 layer before frames can be added", log.LEVEL_ERROR)), 0)
	}
	t.RawSetString("__layerCount", layerCount)

	tableBuilderFunc(state, t, "add", func(state *golua.LState, t *golua.LTable) {
		imgs := state.CheckTable(-1)

		layerCount := t.RawGetString("__layerCount").(golua.LNumber)
		if imgs.Len() != int(layerCount) {
			state.Error(golua.LString(lg.Append(fmt.Sprintf("frame count: %d must match layer count: %d", imgs.Len(), layerCount), log.LEVEL_ERROR)), 0)
		}

		lt := t.RawGetString("__frames").(*golua.LTable)
		lt.Append(imgs)
	})

	assignSpriteNesting(state, t, parent, parent)
	return t
}

func frameBuild(t *golua.LTable) []int {
	frameImgs := make([]int, t.Len())

	for i := range t.Len() {
		frameImgs[i] = int(t.RawGetInt(i + 1).(golua.LNumber))
	}

	return frameImgs
}
