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

	/// @func bbox() -> struct<gamemaker.BBOX>
	/// @arg top {int}
	/// @arg left {int}
	/// @arg bottom {int}
	/// @arg right {int}
	/// @returns {struct<gamemaker.BBOX>}
	lib.CreateFunction(tab, "bbox",
		[]lua.Arg{
			{Type: lua.INT, Name: "top"},
			{Type: lua.INT, Name: "left"},
			{Type: lua.INT, Name: "bottom"},
			{Type: lua.INT, Name: "right"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct BBOX
			/// @prop top {int}
			/// @prop left {int}
			/// @prop bottom {int}
			/// @prop right {int}

			t := state.NewTable()
			t.RawSetString("top", golua.LNumber(args["top"].(int)))
			t.RawSetString("left", golua.LNumber(args["left"].(int)))
			t.RawSetString("bottom", golua.LNumber(args["bottom"].(int)))
			t.RawSetString("right", golua.LNumber(args["right"].(int)))

			state.Push(t)
			return 1
		})

	/// @constants Sprite Origin
	/// @const SPRITEORIGIN_TOPLEFT
	/// @const SPRITEORIGIN_TOPCENTER
	/// @const SPRITEORIGIN_TOPRIGHT
	/// @const SPRITEORIGIN_MIDDLELEFT
	/// @const SPRITEORIGIN_MIDDLECENTER
	/// @const SPRITEORIGIN_MIDDLERIGHT
	/// @const SPRITEORIGIN_BOTTOMLEFT
	/// @const SPRITEORIGIN_BOTTOMCENTER
	/// @const SPRITEORIGIN_BOTTOMRIGHT
	/// @const SPRITEORIGIN_CUSTOM
	tab.RawSetString("SPRITEORIGIN_TOPLEFT", golua.LNumber(yyp.SPRITEORIGIN_TOPLEFT))
	tab.RawSetString("SPRITEORIGIN_TOPCENTER", golua.LNumber(yyp.SPRITEORIGIN_TOPCENTER))
	tab.RawSetString("SPRITEORIGIN_TOPRIGHT", golua.LNumber(yyp.SPRITEORIGIN_TOPRIGHT))
	tab.RawSetString("SPRITEORIGIN_MIDDLELEFT", golua.LNumber(yyp.SPRITEORIGIN_MIDDLELEFT))
	tab.RawSetString("SPRITEORIGIN_MIDDLECENTER", golua.LNumber(yyp.SPRITEORIGIN_MIDDLECENTER))
	tab.RawSetString("SPRITEORIGIN_MIDDLERIGHT", golua.LNumber(yyp.SPRITEORIGIN_MIDDLERIGHT))
	tab.RawSetString("SPRITEORIGIN_BOTTOMLEFT", golua.LNumber(yyp.SPRITEORIGIN_BOTTOMLEFT))
	tab.RawSetString("SPRITEORIGIN_BOTTOMCENTER", golua.LNumber(yyp.SPRITEORIGIN_BOTTOMCENTER))
	tab.RawSetString("SPRITEORIGIN_BOTTOMRIGHT", golua.LNumber(yyp.SPRITEORIGIN_BOTTOMRIGHT))
	tab.RawSetString("SPRITEORIGIN_CUSTOM", golua.LNumber(yyp.SPRITEORIGIN_CUSTOM))

	/// @constants Collision Masks
	/// @const COLLMASK_PRECISE
	/// @const COLLMASK_RECT
	/// @const COLLMASK_ELLIPSE
	/// @const COLLMASK_DIAMOND
	/// @const COLLMASK_PRECISEFRAME
	/// @const COLLMASK_RECTROT
	/// @const COLLMASK_SPINE
	tab.RawSetString("COLLMASK_PRECISE", golua.LNumber(yyp.COLLMASK_PRECISE))
	tab.RawSetString("COLLMASK_RECT", golua.LNumber(yyp.COLLMASK_RECT))
	tab.RawSetString("COLLMASK_ELLIPSE", golua.LNumber(yyp.COLLMASK_ELLIPSE))
	tab.RawSetString("COLLMASK_DIAMOND", golua.LNumber(yyp.COLLMASK_DIAMOND))
	tab.RawSetString("COLLMASK_PRECISEFRAME", golua.LNumber(yyp.COLLMASK_PRECISEFRAME))
	tab.RawSetString("COLLMASK_RECTROT", golua.LNumber(yyp.COLLMASK_RECTROT))
	tab.RawSetString("COLLMASK_SPINE", golua.LNumber(yyp.COLLMASK_SPINE))

	/// @constants BBOX Modes
	/// @const BBOXMODE_AUTO
	/// @const BBOXMODE_FULL
	/// @const BBOXMODE_MANUAL
	tab.RawSetString("BBOXMODE_AUTO", golua.LNumber(yyp.BBOXMODE_AUTO))
	tab.RawSetString("BBOXMODE_FULL", golua.LNumber(yyp.BBOXMODE_FULL))
	tab.RawSetString("BBOXMODE_MANUAL", golua.LNumber(yyp.BBOXMODE_MANUAL))

	/// @constants Nineslice Tile Modes
	/// @const NINESLICETILE_STRETCH
	/// @const NINESLICETILE_REPEAT
	/// @const NINESLICETILE_MIRROR
	/// @const NINESLICETILE_BLANKREPEAT
	/// @const NINESLICETILE_HIDE
	tab.RawSetString("NINESLICETILE_STRETCH", golua.LNumber(yyp.NINESLICETILE_STRETCH))
	tab.RawSetString("NINESLICETILE_REPEAT", golua.LNumber(yyp.NINESLICETILE_REPEAT))
	tab.RawSetString("NINESLICETILE_MIRROR", golua.LNumber(yyp.NINESLICETILE_MIRROR))
	tab.RawSetString("NINESLICETILE_BLANKREPEAT", golua.LNumber(yyp.NINESLICETILE_BLANKREPEAT))
	tab.RawSetString("NINESLICETILE_HIDE", golua.LNumber(yyp.NINESLICETILE_HIDE))

	/// @constants Nineslice Slices
	/// @const NINESLICESLICE_LEFT
	/// @const NINESLICESLICE_TOP
	/// @const NINESLICESLICE_RIGHT
	/// @const NINESLICESLICE_BOTTOM
	/// @const NINESLICESLICE_CENTER
	tab.RawSetString("NINESLICESLICE_LEFT", golua.LNumber(yyp.NINESLICESLICE_LEFT))
	tab.RawSetString("NINESLICESLICE_TOP", golua.LNumber(yyp.NINESLICESLICE_TOP))
	tab.RawSetString("NINESLICESLICE_RIGHT", golua.LNumber(yyp.NINESLICESLICE_RIGHT))
	tab.RawSetString("NINESLICESLICE_BOTTOM", golua.LNumber(yyp.NINESLICESLICE_BOTTOM))
	tab.RawSetString("NINESLICESLICE_CENTER", golua.LNumber(yyp.NINESLICESLICE_CENTER))

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
	/// @struct Sprite
	/// @prop name {string}
	/// @prop width {int}
	/// @prop height {int}
	/// @prop parent {struct<gamemaker.ResourceNode>}
	/// @prop texgroup {struct<gamemaker.ResourceNode>}
	/// @method layers() -> struct<gamemaker.SpriteLayers>
	/// @method frames() -> struct<gamemaker.SpriteFrames>
	/// @method tags([]string) -> self
	/// @method tile(htile bool, vtile bool) -> self
	/// @method origin(int<gamemaker.SpriteOrigin>, xorigin int?, yorigin int?) -> self
	/// @method collision(int<gamemaker.BBOXMode>, int<gamemaker.CollMask>, struct<gamemaker.BBOX>?, tolerance int?) -> self
	/// @method premultiply_alpha(bool) -> self
	/// @method edge_filtering(bool) -> self
	/// @method dynamic_texturepage(bool) -> self
	/// @method nineslice(top int, left int, bottom int, right int) -> self
	/// @method nineslice_tilemode(int<gamemaker.NineSliceSlice>, int<gamemaker.NineSliceTile>) -> self
	/// @method broadcast_message(frame int, msg string) -> self

	t := state.NewTable()
	t.RawSetString("name", golua.LString(name))
	t.RawSetString("width", golua.LNumber(width))
	t.RawSetString("height", golua.LNumber(height))
	t.RawSetString("parent", parent)
	t.RawSetString("texgroup", texgroup)

	t.RawSetString("__tags", golua.LNil)
	t.RawSetString("__htile", golua.LNil)
	t.RawSetString("__vtile", golua.LNil)
	t.RawSetString("__origin", golua.LNil)
	t.RawSetString("__xorigin", golua.LNil)
	t.RawSetString("__yorigin", golua.LNil)
	t.RawSetString("__bboxMode", golua.LNil)
	t.RawSetString("__mask", golua.LNil)
	t.RawSetString("__bbox", golua.LNil)
	t.RawSetString("__tolerance", golua.LNil)
	t.RawSetString("__preAlpha", golua.LNil)
	t.RawSetString("__edgeFiltering", golua.LNil)
	t.RawSetString("__dynamicTexturePage", golua.LNil)
	t.RawSetString("__nineslice", golua.LFalse)
	t.RawSetString("__ninesliceTop", golua.LNumber(0))
	t.RawSetString("__ninesliceLeft", golua.LNumber(0))
	t.RawSetString("__ninesliceBottom", golua.LNumber(0))
	t.RawSetString("__ninesliceRight", golua.LNumber(0))
	t.RawSetString("__ninesliceTiles", golua.LNil)
	t.RawSetString("__messages", golua.LNil)

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

	tableBuilderFunc(state, t, "tags", func(state *golua.LState, t *golua.LTable) {
		st := state.CheckTable(-1)
		t.RawSetString("__tags", st)
	})

	tableBuilderFunc(state, t, "tile", func(state *golua.LState, t *golua.LTable) {
		htile := state.CheckBool(-2)
		vtile := state.CheckBool(-1)
		t.RawSetString("__htile", golua.LBool(htile))
		t.RawSetString("__vtile", golua.LBool(vtile))
	})

	tableBuilderFunc(state, t, "origin", func(state *golua.LState, t *golua.LTable) {
		origin := state.CheckNumber(2)
		t.RawSetString("__origin", origin)

		if yyp.SpriteOrigin(origin) == yyp.SPRITEORIGIN_CUSTOM {
			xorigin := state.CheckNumber(3)
			yorigin := state.CheckNumber(4)

			t.RawSetString("__xorigin", xorigin)
			t.RawSetString("__yorigin", yorigin)
		}
	})

	tableBuilderFunc(state, t, "collision", func(state *golua.LState, t *golua.LTable) {
		bboxMode := state.CheckNumber(2)
		mask := state.CheckNumber(3)

		t.RawSetString("__bboxMode", bboxMode)
		t.RawSetString("__mask", mask)

		bbox := state.Get(4)
		tolerance := state.Get(5)

		if bbox.Type() == golua.LTTable {
			t.RawSetString("__bbox", bbox)
		}

		if tolerance.Type() == golua.LTNumber {
			t.RawSetString("__tolerance", tolerance)
		}
	})

	tableBuilderFunc(state, t, "premultiply_alpha", func(state *golua.LState, t *golua.LTable) {
		pa := state.CheckBool(2)
		t.RawSetString("__preAlpha", golua.LBool(pa))
	})

	tableBuilderFunc(state, t, "edge_filtering", func(state *golua.LState, t *golua.LTable) {
		ef := state.CheckBool(2)
		t.RawSetString("__edgeFiltering", golua.LBool(ef))
	})

	tableBuilderFunc(state, t, "dynamic_texturepage", func(state *golua.LState, t *golua.LTable) {
		dtp := state.CheckBool(2)
		t.RawSetString("__dynamicTexturePage", golua.LBool(dtp))
	})

	tableBuilderFunc(state, t, "nineslice", func(state *golua.LState, t *golua.LTable) {
		t.RawSetString("__nineslice", golua.LTrue)

		top := state.CheckNumber(2)
		left := state.CheckNumber(3)
		bottom := state.CheckNumber(4)
		right := state.CheckNumber(5)

		t.RawSetString("__ninesliceTop", top)
		t.RawSetString("__ninesliceLeft", left)
		t.RawSetString("__ninesliceBottom", bottom)
		t.RawSetString("__ninesliceRight", right)
	})

	tableBuilderFunc(state, t, "nineslice_tilemode", func(state *golua.LState, t *golua.LTable) {
		tile := state.CheckNumber(2)
		mode := state.CheckNumber(3)

		tt := state.NewTable()
		tt.RawSetString("__tile", tile)
		tt.RawSetString("__mode", mode)

		tiles := t.RawGetString("__ninesliceTiles")
		if tiles.Type() == golua.LTNil {
			tiles = state.NewTable()
		}

		tiles.(*golua.LTable).Append(tt)
		t.RawSetString("__ninesliceTiles", tiles)
	})

	tableBuilderFunc(state, t, "broadcast_message", func(state *golua.LState, t *golua.LTable) {
		key := state.CheckNumber(2)
		msg := state.CheckString(3)

		mt := state.NewTable()
		mt.RawSetString("__key", key)
		mt.RawSetString("__msg", golua.LString(msg))

		messages := t.RawGetString("__messages")
		if messages.Type() == golua.LTNil {
			messages = state.NewTable()
		}

		messages.(*golua.LTable).Append(mt)
		t.RawSetString("__messages", messages)
	})

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

	sprite, err := yyp.NewSprite(name, parent, texgroup, width, height, layers)
	if err != nil {
		return nil, err
	}

	tags := t.RawGetString("__tags")
	if tags.Type() == golua.LTTable {
		tgs := tags.(*golua.LTable)
		tagList := make([]string, tgs.Len())
		for i := range tgs.Len() {
			tagList[i] = string(tgs.RawGetInt(i + 1).(golua.LString))
		}
		sprite.Resource.Tags = tagList
	}

	htile := t.RawGetString("__htile")
	if htile.Type() == golua.LTBool {
		sprite.Resource.HTile = bool(htile.(golua.LBool))
	}

	vtile := t.RawGetString("__vtile")
	if vtile.Type() == golua.LTBool {
		sprite.Resource.VTile = bool(vtile.(golua.LBool))
	}

	origin := t.RawGetString("__origin")
	if origin.Type() == golua.LTNumber {
		xorg := 0
		yorg := 0

		sorg := yyp.SpriteOrigin(origin.(golua.LNumber))
		switch sorg {
		case yyp.SPRITEORIGIN_TOPLEFT:
			xorg = 0
			yorg = 0
		case yyp.SPRITEORIGIN_TOPCENTER:
			xorg = width / 2
			yorg = 0
		case yyp.SPRITEORIGIN_TOPRIGHT:
			xorg = width
			yorg = 0
		case yyp.SPRITEORIGIN_MIDDLELEFT:
			xorg = 0
			yorg = height / 2
		case yyp.SPRITEORIGIN_MIDDLECENTER:
			xorg = width / 2
			yorg = height / 2
		case yyp.SPRITEORIGIN_MIDDLERIGHT:
			xorg = width
			yorg = height / 2
		case yyp.SPRITEORIGIN_BOTTOMLEFT:
			xorg = 0
			yorg = height
		case yyp.SPRITEORIGIN_BOTTOMCENTER:
			xorg = width / 2
			yorg = height
		case yyp.SPRITEORIGIN_BOTTOMRIGHT:
			xorg = width
			yorg = height
		case yyp.SPRITEORIGIN_CUSTOM:
			xorg = int(t.RawGetString("__xorigin").(golua.LNumber))
			yorg = int(t.RawGetString("__yorigin").(golua.LNumber))
		default:
			return nil, fmt.Errorf("unknown sprite origin: %d", origin)
		}

		sprite.Resource.Origin = sorg
		sprite.Resource.Sequence.XOrigin = xorg
		sprite.Resource.Sequence.YOrigin = yorg
	}

	bboxMode := t.RawGetString("__bboxMode")
	if bboxMode.Type() == golua.LTNumber {
		sprite.Resource.BBOX_Mode = yyp.BBOXMode(bboxMode.(golua.LNumber))
	}

	mask := t.RawGetString("__mask")
	if mask.Type() == golua.LTNumber {
		sprite.Resource.CollisionKind = yyp.CollMask(mask.(golua.LNumber))
	}

	bbox := t.RawGetString("__bbox")
	if bbox.Type() == golua.LTTable {
		b := bbox.(*golua.LTable)
		top := b.RawGetString("top").(golua.LNumber)
		left := b.RawGetString("left").(golua.LNumber)
		bottom := b.RawGetString("bottom").(golua.LNumber)
		right := b.RawGetString("right").(golua.LNumber)

		sprite.Resource.BBOX_Top = int(top)
		sprite.Resource.BBOX_Left = int(left)
		sprite.Resource.BBOX_Bottom = int(bottom)
		sprite.Resource.BBOX_Right = int(right)
	}

	tolerance := t.RawGetString("__tolerance")
	if tolerance.Type() == golua.LTNumber {
		sprite.Resource.CollisionTolerance = int(tolerance.(golua.LNumber))
	}

	preAlpha := t.RawGetString("__preAlpha")
	if preAlpha.Type() == golua.LTBool {
		sprite.Resource.PreMultiplyAlpha = bool(preAlpha.(golua.LBool))
	}

	edgeFiltering := t.RawGetString("__edgeFiltering")
	if edgeFiltering.Type() == golua.LTBool {
		sprite.Resource.EdgeFiltering = bool(edgeFiltering.(golua.LBool))
	}

	dtp := t.RawGetString("__dynamicTexturePage")
	if dtp.Type() == golua.LTBool {
		sprite.Resource.DynamicTexturePage = bool(dtp.(golua.LBool))
	}

	nineslice := t.RawGetString("__nineslice")
	if nineslice.Type() == golua.LTBool && bool(nineslice.(golua.LBool)) {
		top := int(t.RawGetString("__ninesliceTop").(golua.LNumber))
		left := int(t.RawGetString("__ninesliceLeft").(golua.LNumber))
		bottom := int(t.RawGetString("__ninesliceBottom").(golua.LNumber))
		right := int(t.RawGetString("__ninesliceRight").(golua.LNumber))

		sprite.Resource.NineSlice = yyp.NewResourceNineSlice()
		sprite.Resource.NineSlice.Top = top
		sprite.Resource.NineSlice.Left = left
		sprite.Resource.NineSlice.Bottom = bottom
		sprite.Resource.NineSlice.Right = right

		tiles := t.RawGetString("__ninesliceTiles")
		if tiles.Type() == golua.LTTable {
			tt := tiles.(*golua.LTable)
			for i := range tt.Len() {
				tile := tt.RawGetInt(i + 1).(*golua.LTable)
				tileIndex := tile.RawGetString("__tile").(golua.LNumber)
				mode := tile.RawGetString("__mode").(golua.LNumber)

				sprite.Resource.NineSlice.TileMode[int(tileIndex)] = yyp.NineSliceTile(mode)
			}
		}
	}

	messages := t.RawGetString("__messages")
	if messages.Type() == golua.LTTable {
		mt := messages.(*golua.LTable)
		for i := range mt.Len() {
			msg := mt.RawGetInt(i + 1).(*golua.LTable)
			key := float64(msg.RawGetString("__key").(golua.LNumber))
			message := string(msg.RawGetString("__msg").(golua.LString))

			sprite.Resource.Sequence.Events.Keyframes = append(sprite.Resource.Sequence.Events.Keyframes, yyp.NewResourceSpriteSequenceEventKeyframe([][]string{{message}}, key))
		}
	}

	return sprite, nil
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
	/// @struct SpriteLayers
	/// @method image(name) -> self
	/// @method default() -> self
	/// @method folder(name) -> struct<gamemaker.SpriteLayerFolder>
	/// @method back() -> struct<gamemaker.Sprite>

	t := state.NewTable()
	layers := parent.RawGetString("__layers").(*golua.LTable)
	assignSpriteLayers(state, t, layers)
	assignSpriteNesting(state, t, parent, parent)
	return t
}

func layerImageTable(state *golua.LState, name string) *golua.LTable {
	/// @struct SpriteLayerImage
	/// @prop type {string<gamemaker.LayerType>}
	/// @prop name {string}

	t := state.NewTable()

	t.RawSetString("type", golua.LString(LAYER_IMAGE))
	t.RawSetString("name", golua.LString(name))

	return t
}

func layerFolderTable(state *golua.LState, name string, parent *golua.LTable) *golua.LTable {
	/// @struct SpriteLayerFolder
	/// @prop type {string<gamemaker.LayerType>}
	/// @prop name {string}
	/// @method image(name) -> self
	/// @method default() -> self
	/// @method folder(name) -> struct<gamemaker.SpriteLayerFolder>
	/// @method back() -> struct<gamemaker.SpriteLayers>

	t := state.NewTable()

	t.RawSetString("type", golua.LString(LAYER_FOLDER))
	t.RawSetString("name", golua.LString(name))

	base := parent.RawGetString("__base").(*golua.LTable)

	assignSpriteLayers(state, t, state.NewTable())
	assignSpriteNesting(state, t, parent, base)

	return t
}

func framesTable(state *golua.LState, parent *golua.LTable, lg *log.Logger) *golua.LTable {
	/// @struct SpriteFrames
	/// @method add([]int<collection.IMAGE>) -> self
	/// @method back() -> struct<gamemaker.Sprite>

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
