package lib

import (
	"fmt"
	"image"
	"os"
	"sync"

	"github.com/ArtificialLegacy/gm-proj-tool/yyp"
	"github.com/ArtificialLegacy/imgscal/pkg/byteseeker"
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

	/// @func folder_add(id, name, folderpath) -> struct<gamemaker.ResourceNode>, string
	/// @arg id {int<collection.CRATE_GAMEMAKER>} - ID for the loaded Gamemaker project.
	/// @arg name {string} - Name of the folder asset.
	/// @arg? folderpath {string} - Parent path for the folder asset, use an empty string for the root folder.
	/// @arg? tags {[]string} - List of tags to assign to the folder.
	/// @returns {struct<gamemaker.ResourceNode>}
	/// @returns {string} - Path to the folder asset.
	lib.CreateFunction(tab, "folder_add",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "name"},
			{Type: lua.STRING, Name: "path", Optional: true},
			lua.ArgVariadic("tags", lua.ArrayType{Type: lua.STRING}, true),
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			proj, err := r.CR_GMP.Item(args["id"].(int))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to find project: %d, %s", args["id"], err), log.LEVEL_ERROR)), 0)
			}

			name := args["name"].(string)
			folderpath := args["path"].(string)

			folder := yyp.NewFolder(name, folderpath)

			tags := args["tags"].([]any)
			if len(tags) > 0 {
				tagList := make([]string, len(tags))
				for i, v := range tags {
					tagList[i] = v.(string)
				}
				folder.Resource.Tags = tagList
			}

			err = proj.FolderSave(folder)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to save folder: %s", err), log.LEVEL_ERROR)), 0)
			}

			node := folder.AsParent()
			t := resourceNodeTable(state, node.Name, node.Path)

			newPath := folder.FolderPath()

			state.Push(t)
			state.Push(golua.LString(newPath))
			return 2
		})

	/// @func folder_get(id, folderpath) -> struct<gamemaker.ResourceNode>
	/// @arg id {int<collection.CRATE_GAMEMAKER>} - ID for the loaded Gamemaker project.
	/// @arg folderpath {string} - Path to the folder asset.
	/// @returns {struct<gamemaker.ResourceNode>}
	lib.CreateFunction(tab, "folder_get",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "path"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			proj, err := r.CR_GMP.Item(args["id"].(int))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to find project: %d, %s", args["id"], err), log.LEVEL_ERROR)), 0)
			}

			folderpath := args["path"].(string)

			folder, err := proj.FolderLoad(folderpath)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to load folder: %s", err), log.LEVEL_ERROR)), 0)
			}

			node := folder.AsParent()
			t := resourceNodeTable(state, node.Name, node.Path)

			state.Push(t)
			return 1
		})

	/// @func folder_delete(id, folderpath)
	/// @arg id {int<collection.CRATE_GAMEMAKER>} - ID for the loaded Gamemaker project.
	/// @arg folderpath {string} - Path to the folder asset.
	lib.CreateFunction(tab, "folder_delete",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "path"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			proj, err := r.CR_GMP.Item(args["id"].(int))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to find project: %d, %s", args["id"], err), log.LEVEL_ERROR)), 0)
			}

			folderpath := args["path"].(string)

			err = proj.FolderDelete(folderpath)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to delete folder: %s", err), log.LEVEL_ERROR)), 0)
			}

			return 0
		})

	/// @func folder_exists(id, folderpath) -> bool
	/// @arg id {int<collection.CRATE_GAMEMAKER>} - ID for the loaded Gamemaker project.
	/// @arg folderpath {string} - Path to the folder asset.
	/// @returns {bool}
	lib.CreateFunction(tab, "folder_exists",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "path"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			proj, err := r.CR_GMP.Item(args["id"].(int))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to find project: %d, %s", args["id"], err), log.LEVEL_ERROR)), 0)
			}

			folderpath := args["path"].(string)
			exists := proj.FolderExists(folderpath)

			state.Push(golua.LBool(exists))
			return 1
		})

	/// @func sprite(name, width, height, parent, texgroup) -> struct<gamemaker.Sprite>
	/// @arg name {string} - Name of the sprite asset.
	/// @arg width {int}
	/// @arg height {int}
	/// @arg parent {struct<gamemaker.ResourceNode>}
	/// @arg texgroup {struct<gamemaker.ResourceNode>}
	/// @returns {struct<gamemaker.Sprite>}
	lib.CreateFunction(tab, "sprite",
		[]lua.Arg{
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "width"},
			{Type: lua.INT, Name: "height"},
			{Type: lua.RAW_TABLE, Name: "parent"},
			{Type: lua.RAW_TABLE, Name: "texgroup"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := spriteTable(lib, state,
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
			{Type: lua.RAW_TABLE, Name: "sprite"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			proj, err := r.CR_GMP.Item(args["id"].(int))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to find project: %d, %s", args["id"], err), log.LEVEL_ERROR)), 0)
			}

			sprite, err := spriteBuild(args["sprite"].(*golua.LTable), r, state)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to build sprite: %s", err), log.LEVEL_ERROR)), 0)
			}

			err = proj.ImportResource(sprite)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to import sprite: %s", err), log.LEVEL_ERROR)), 0)
			}

			return 0
		})

	/// @func sprite_delete(id, name)
	/// @arg id {int<collection.CRATE_GAMEMAKER>} - ID for the loaded Gamemaker project.
	/// @arg name {string} - Name of the sprite asset.
	lib.CreateFunction(tab, "sprite_delete",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "name"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			proj, err := r.CR_GMP.Item(args["id"].(int))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to find project: %d, %s", args["id"], err), log.LEVEL_ERROR)), 0)
			}

			name := args["name"].(string)
			err = proj.SpriteDelete(name)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to delete sprite: %s", err), log.LEVEL_ERROR)), 0)
			}

			return 0
		})

	/// @func sprite_exists(id, name) -> bool
	/// @arg id {int<collection.CRATE_GAMEMAKER>} - ID for the loaded Gamemaker project.
	/// @arg name {string} - Name of the sprite asset.
	/// @returns {bool}
	lib.CreateFunction(tab, "sprite_exists",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "name"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			proj, err := r.CR_GMP.Item(args["id"].(int))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to find project: %d, %s", args["id"], err), log.LEVEL_ERROR)), 0)
			}

			name := args["name"].(string)
			exists := proj.SpriteExists(name)

			state.Push(golua.LBool(exists))
			return 1
		})

	/// @func sprite_load(id, name, encoding) -> struct<gamemaker.Sprite>
	/// @arg id {int<collection.CRATE_GAMEMAKER>} - ID for the loaded Gamemaker project.
	/// @arg name {string} - Name of the sprite asset.
	/// @arg encoding {int<image.Encoding}
	/// @returns {struct<gamemaker.Sprite>}
	lib.CreateFunction(tab, "sprite_load",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "encoding"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			proj, err := r.CR_GMP.Item(args["id"].(int))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to find project: %d, %s", args["id"], err), log.LEVEL_ERROR)), 0)
			}

			name := args["name"].(string)

			sprite, err := proj.SpriteLoad(name)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to load sprite: %s", err), log.LEVEL_ERROR)), 0)
			}

			parent := resourceNodeTable(state, sprite.Resource.Parent.Name, sprite.Resource.Parent.Path)
			texgroup := resourceNodeTable(state, sprite.Resource.TexGroupID.Name, sprite.Resource.TexGroupID.Path)
			t := spriteTable(lib, state, sprite.Name, sprite.Resource.Width, sprite.Resource.Height, parent, texgroup, lg)

			layerParent := layersTable(state, t)
			layers := state.NewTable()
			frames := state.NewTable()
			encoding := args["encoding"].(int)

			for _, layer := range sprite.Layers {
				l := layersParse(r, lg, state, layerParent, layer, frames, imageutil.ImageEncoding(encoding))
				layers.Append(l)
			}
			t.RawSetString("__layers", layers)
			t.RawSetString("__frames", frames)

			if sprite.Resource.Tags != nil {
				tags := state.NewTable()
				for i, tag := range sprite.Resource.Tags {
					tags.RawSetInt(i+1, golua.LString(tag))
				}
				t.RawSetString("__tags", tags)
			}

			if sprite.Resource.HTile {
				t.RawSetString("__htile", golua.LBool(sprite.Resource.HTile))
			}
			if sprite.Resource.VTile {
				t.RawSetString("__vtile", golua.LBool(sprite.Resource.VTile))
			}

			t.RawSetString("__origin", golua.LNumber(sprite.Resource.Origin))
			if sprite.Resource.Origin == yyp.SPRITEORIGIN_CUSTOM {
				t.RawSetString("__xorigin", golua.LNumber(sprite.Resource.Sequence.XOrigin))
				t.RawSetString("__yorigin", golua.LNumber(sprite.Resource.Sequence.YOrigin))
			}

			t.RawSetString("__collmask", golua.LNumber(sprite.Resource.CollisionKind))
			t.RawSetString("__bboxmode", golua.LNumber(sprite.Resource.BBOX_Mode))
			if sprite.Resource.BBOX_Mode == yyp.BBOXMODE_MANUAL {
				t.RawSetString("__bbox", bboxTable(state, sprite.Resource.BBOX_Top, sprite.Resource.BBOX_Left, sprite.Resource.BBOX_Bottom, sprite.Resource.BBOX_Right))
			}
			t.RawSetString("__tolerance", golua.LNumber(sprite.Resource.CollisionTolerance))

			t.RawSetString("__preAlpha", golua.LBool(sprite.Resource.PreMultiplyAlpha))
			t.RawSetString("__edgeFiltering", golua.LBool(sprite.Resource.EdgeFiltering))
			t.RawSetString("__dynamicTexturePage", golua.LBool(sprite.Resource.DynamicTexturePage))

			if sprite.Resource.NineSlice != nil {
				t.RawSetString("__nineslice", golua.LBool(sprite.Resource.NineSlice.Enabled))

				t.RawSetString("__ninesliceTop", golua.LNumber(sprite.Resource.NineSlice.Top))
				t.RawSetString("__ninesliceLeft", golua.LNumber(sprite.Resource.NineSlice.Left))
				t.RawSetString("__ninesliceBottom", golua.LNumber(sprite.Resource.NineSlice.Bottom))
				t.RawSetString("__ninesliceRight", golua.LNumber(sprite.Resource.NineSlice.Right))

				tiles := state.NewTable()

				for i, tm := range sprite.Resource.NineSlice.TileMode {
					if tm != yyp.NINESLICETILE_STRETCH {
						tile := state.NewTable()
						tile.RawSetString("__mode", golua.LNumber(tm))
						tile.RawSetString("__tile", golua.LNumber(i))
						tiles.Append(tile)
					}
				}

				if tiles.Len() > 0 {
					t.RawSetString("__ninesliceTiles", tiles)
				}
			}

			if len(sprite.Resource.Sequence.Events.Keyframes) > 0 {
				msgs := state.NewTable()

				for _, kf := range sprite.Resource.Sequence.Events.Keyframes {
					msg := state.NewTable()

					msg.RawSetString("__key", golua.LNumber(kf.Key))
					msg.RawSetString("__msg", golua.LString(kf.Channels["0"].Events[0]))

					msgs.Append(msg)
				}

				t.RawSetString("__messages", msgs)
			}

			t.RawSetString("__playbackSpeed", golua.LNumber(sprite.Resource.Sequence.PlaybackSpeed))
			t.RawSetString("__playbackUnits", golua.LNumber(sprite.Resource.Sequence.TimeUnits))

			state.Push(t)
			return 1
		})

	/// @func note(name, text, parent) -> struct<gamemaker.Note>
	/// @arg name {string} - Name of the note asset.
	/// @arg text {string}
	/// @arg parent {struct<gamemaker.ResourceNode>}
	/// @returns {struct<gamemaker.Note>}
	lib.CreateFunction(tab, "note",
		[]lua.Arg{
			{Type: lua.STRING, Name: "name"},
			{Type: lua.STRING, Name: "text"},
			{Type: lua.RAW_TABLE, Name: "parent"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := noteTable(lib, state, args["name"].(string), args["text"].(string), args["parent"].(*golua.LTable))

			state.Push(t)
			return 1
		})

	/// @func note_save(id, note)
	/// @arg id {int<collection.CRATE_GAMEMAKER>} - ID for the loaded Gamemaker project.
	/// @arg note {struct<gamemaker.Note>}
	lib.CreateFunction(tab, "note_save",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.RAW_TABLE, Name: "note"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			proj, err := r.CR_GMP.Item(args["id"].(int))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to find project: %d, %s", args["id"], err), log.LEVEL_ERROR)), 0)
			}

			note := noteBuild(args["note"].(*golua.LTable))

			err = proj.ImportResource(note)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to import note: %s", err), log.LEVEL_ERROR)), 0)
			}

			return 0
		})

	/// @func note_delete(id, name)
	/// @arg id {int<collection.CRATE_GAMEMAKER>} - ID for the loaded Gamemaker project.
	/// @arg name {string} - Name of the note asset.
	lib.CreateFunction(tab, "note_delete",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "name"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			proj, err := r.CR_GMP.Item(args["id"].(int))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to find project: %d, %s", args["id"], err), log.LEVEL_ERROR)), 0)
			}

			name := args["name"].(string)
			err = proj.NoteDelete(name)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to delete note: %s", err), log.LEVEL_ERROR)), 0)
			}

			return 0
		})

	/// @func note_exists(id, name) -> bool
	/// @arg id {int<collection.CRATE_GAMEMAKER>} - ID for the loaded Gamemaker project.
	/// @arg name {string} - Name of the note asset.
	/// @returns {bool}
	lib.CreateFunction(tab, "note_exists",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "name"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			proj, err := r.CR_GMP.Item(args["id"].(int))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to find project: %d, %s", args["id"], err), log.LEVEL_ERROR)), 0)
			}

			name := args["name"].(string)
			exists := proj.NoteExists(name)

			state.Push(golua.LBool(exists))
			return 1
		})

	/// @func note_load(id, name) -> struct<gamemaker.Note>
	/// @arg id {int<collection.CRATE_GAMEMAKER>} - ID for the loaded Gamemaker project.
	/// @arg name {string} - Name of the note asset.
	/// @returns {struct<gamemaker.Note>}
	lib.CreateFunction(tab, "note_load",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "name"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			proj, err := r.CR_GMP.Item(args["id"].(int))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to find project: %d, %s", args["id"], err), log.LEVEL_ERROR)), 0)
			}

			name := args["name"].(string)

			note, err := proj.NoteLoad(name)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to load note: %s", err), log.LEVEL_ERROR)), 0)
			}

			parent := resourceNodeTable(state, note.Resource.Parent.Name, note.Resource.Parent.Path)
			t := noteTable(lib, state, note.Name, note.Text, parent)

			if note.Resource.Tags != nil {
				tags := state.NewTable()
				for i, tag := range note.Resource.Tags {
					tags.RawSetInt(i+1, golua.LString(tag))
				}
				t.RawSetString("__tags", tags)
			}

			state.Push(t)
			return 1
		})

	/// @func script(name, code, parent) -> struct<gamemaker.Script>
	/// @arg name {string} - Name of the script asset.
	/// @arg code {string}
	/// @arg parent {struct<gamemaker.ResourceNode>}
	/// @returns {struct<gamemaker.Script>}
	lib.CreateFunction(tab, "script",
		[]lua.Arg{
			{Type: lua.STRING, Name: "name"},
			{Type: lua.STRING, Name: "code"},
			{Type: lua.RAW_TABLE, Name: "parent"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := scriptTable(lib, state, args["name"].(string), args["code"].(string), args["parent"].(*golua.LTable))

			state.Push(t)
			return 1
		})

	/// @func script_save(id, script)
	/// @arg id {int<collection.CRATE_GAMEMAKER>} - ID for the loaded Gamemaker project.
	/// @arg script {struct<gamemaker.Script>}
	lib.CreateFunction(tab, "script_save",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.RAW_TABLE, Name: "script"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			proj, err := r.CR_GMP.Item(args["id"].(int))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to find project: %d, %s", args["id"], err), log.LEVEL_ERROR)), 0)
			}

			script := scriptBuild(args["script"].(*golua.LTable))

			err = proj.ImportResource(script)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to import script: %s", err), log.LEVEL_ERROR)), 0)
			}

			return 0
		})

	/// @func script_delete(id, name)
	/// @arg id {int<collection.CRATE_GAMEMAKER>} - ID for the loaded Gamemaker project.
	/// @arg name {string} - Name of the script asset.
	lib.CreateFunction(tab, "script_delete",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "name"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			proj, err := r.CR_GMP.Item(args["id"].(int))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to find project: %d, %s", args["id"], err), log.LEVEL_ERROR)), 0)
			}

			name := args["name"].(string)
			err = proj.ScriptDelete(name)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to delete script: %s", err), log.LEVEL_ERROR)), 0)
			}

			return 0
		})

	/// @func script_exists(id, name) -> bool
	/// @arg id {int<collection.CRATE_GAMEMAKER>} - ID for the loaded Gamemaker project.
	/// @arg name {string} - Name of the script asset.
	/// @returns {bool}
	lib.CreateFunction(tab, "script_exists",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "name"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			proj, err := r.CR_GMP.Item(args["id"].(int))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to find project: %d, %s", args["id"], err), log.LEVEL_ERROR)), 0)
			}

			name := args["name"].(string)
			exists := proj.ScriptExists(name)

			state.Push(golua.LBool(exists))
			return 1
		})

	/// @func script_load(id, name) -> struct<gamemaker.Script>
	/// @arg id {int<collection.CRATE_GAMEMAKER>} - ID for the loaded Gamemaker project.
	/// @arg name {string} - Name of the script asset.
	/// @returns {struct<gamemaker.Script>}
	lib.CreateFunction(tab, "script_load",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "name"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			proj, err := r.CR_GMP.Item(args["id"].(int))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to find project: %d, %s", args["id"], err), log.LEVEL_ERROR)), 0)
			}

			name := args["name"].(string)

			script, err := proj.ScriptLoad(name)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to load script: %s", err), log.LEVEL_ERROR)), 0)
			}

			parent := resourceNodeTable(state, script.Resource.Parent.Name, script.Resource.Parent.Path)
			t := scriptTable(lib, state, script.Name, script.Code, parent)

			if script.Resource.Tags != nil {
				tags := state.NewTable()
				for i, tag := range script.Resource.Tags {
					tags.RawSetInt(i+1, golua.LString(tag))
				}
				t.RawSetString("__tags", tags)
			}

			state.Push(t)
			return 1
		})

	/// @func datafile_from_file(name, filepath, frompath) -> struct<gamemaker.DataFile>
	/// @arg name {string} - Name of the script asset.
	/// @arg filepath {string}
	/// @arg frompath {string} - File path to read file from. Note that this file is not read until the resource is saved.
	/// @returns {struct<gamemaker.DataFile>}
	lib.CreateFunction(tab, "datafile_from_file",
		[]lua.Arg{
			{Type: lua.STRING, Name: "name"},
			{Type: lua.STRING, Name: "filepath"},
			{Type: lua.STRING, Name: "frompath"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := datafileTable(state, args["name"].(string), args["filepath"].(string), golua.LString(args["frompath"].(string)), DATAFILE_FILE)

			state.Push(t)
			return 1
		})

	/// @func datafile_from_string(name, filepath, data) -> struct<gamemaker.DataFile>
	/// @arg name {string} - Name of the script asset.
	/// @arg filepath {string}
	/// @arg data {string}
	/// @returns {struct<gamemaker.DataFile>}
	lib.CreateFunction(tab, "datafile_from_string",
		[]lua.Arg{
			{Type: lua.STRING, Name: "name"},
			{Type: lua.STRING, Name: "filepath"},
			{Type: lua.STRING, Name: "data"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := datafileTable(state, args["name"].(string), args["filepath"].(string), golua.LString(args["data"].(string)), DATAFILE_STRING)

			state.Push(t)
			return 1
		})

	/// @func datafile_from_image(name, filepath, id) -> struct<gamemaker.DataFile>
	/// @arg name {string} - Name of the script asset.
	/// @arg filepath {string}
	/// @arg id {int<collection.IMAGE>}
	/// @returns {struct<gamemaker.DataFile>}
	lib.CreateFunction(tab, "datafile_from_image",
		[]lua.Arg{
			{Type: lua.STRING, Name: "name"},
			{Type: lua.STRING, Name: "filepath"},
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := datafileTable(state, args["name"].(string), args["filepath"].(string), golua.LNumber(args["id"].(int)), DATAFILE_IMAGE)

			state.Push(t)
			return 1
		})

	/// @func datafile_save(id, datafile)
	/// @arg id {int<collection.CRATE_GAMEMAKER>} - ID for the loaded Gamemaker project.
	/// @arg datafile {struct<gamemaker.DataFile>}
	lib.CreateFunction(tab, "datafile_save",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.RAW_TABLE, Name: "datafile"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			proj, err := r.CR_GMP.Item(args["id"].(int))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to find project: %d, %s", args["id"], err), log.LEVEL_ERROR)), 0)
			}

			datafile := datafileBuild(state, args["datafile"].(*golua.LTable), r, lg)

			err = proj.IncludedFileSave(datafile)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to import datafile: %s", err), log.LEVEL_ERROR)), 0)
			}

			return 0
		})

	/// @func datafile_delete(id, filepath, name)
	/// @arg id {int<collection.CRATE_GAMEMAKER>} - ID for the loaded Gamemaker project.
	/// @arg filepath {string}
	/// @arg name {string}
	lib.CreateFunction(tab, "datafile_delete",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "filepath"},
			{Type: lua.STRING, Name: "name"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			proj, err := r.CR_GMP.Item(args["id"].(int))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to find project: %d, %s", args["id"], err), log.LEVEL_ERROR)), 0)
			}

			err = proj.IncludedFileDelete(args["filepath"].(string), args["name"].(string))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to delete datafile: %s", err), log.LEVEL_ERROR)), 0)
			}

			return 0
		})

	/// @func datafile_exists(id, filepath, name) -> bool
	/// @arg id {int<collection.CRATE_GAMEMAKER>} - ID for the loaded Gamemaker project.
	/// @arg filepath {string}
	/// @arg name {string}
	/// @returns {bool}
	lib.CreateFunction(tab, "datafile_exists",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "filepath"},
			{Type: lua.STRING, Name: "name"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			proj, err := r.CR_GMP.Item(args["id"].(int))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to find project: %d, %s", args["id"], err), log.LEVEL_ERROR)), 0)
			}

			exists := proj.IncludedFileExists(args["filepath"].(string), args["name"].(string))

			state.Push(golua.LBool(exists))
			return 1
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
			t := bboxTable(state, args["top"].(int), args["left"].(int), args["bottom"].(int), args["right"].(int))

			state.Push(t)
			return 1
		})

	/// @constants SpriteOrigin {int}
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

	/// @constants CollMask {int}
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

	/// @constants BBOXMode {int}
	/// @const BBOXMODE_AUTO
	/// @const BBOXMODE_FULL
	/// @const BBOXMODE_MANUAL
	tab.RawSetString("BBOXMODE_AUTO", golua.LNumber(yyp.BBOXMODE_AUTO))
	tab.RawSetString("BBOXMODE_FULL", golua.LNumber(yyp.BBOXMODE_FULL))
	tab.RawSetString("BBOXMODE_MANUAL", golua.LNumber(yyp.BBOXMODE_MANUAL))

	/// @constants NinesliceTile {int}
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

	/// @constants NinesliceSlice {int}
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

	/// @constants PlaybackUnit {int}
	/// @const PLAYBACK_PERSECOND
	/// @const PLAYBACK_PERFRAME
	tab.RawSetString("PLAYBACK_PERSECOND", golua.LNumber(yyp.SEQUNITS_TIME))
	tab.RawSetString("PLAYBACK_PERFRAME", golua.LNumber(yyp.SEQUNITS_FRAME))

	/// @constants Directories {string}
	/// @const DIR_DATAFILES
	tab.RawSetString("DIR_DATAFILES", golua.LString(yyp.INCLUDEDFILE_DEFAULTPATH))

	/// @constants Color {int}
	/// @const COLOR_AQUA
	/// @const COLOR_BLACK
	/// @const COLOR_BLUE
	/// @const COLOR_DKGRAY
	/// @const COLOR_FUCHSIA
	/// @const COLOR_GRAY
	/// @const COLOR_GREEN
	/// @const COLOR_LIME
	/// @const COLOR_LTGRAY
	/// @const COLOR_MAROON
	/// @const COLOR_NAVY
	/// @const COLOR_OLIVE
	/// @const COLOR_ORANGE
	/// @const COLOR_PURPLE
	/// @const COLOR_RED
	/// @const COLOR_SILVER
	/// @const COLOR_TEAL
	/// @const COLOR_WHITE
	/// @const COLOR_YELLOW
	tab.RawSetString("COLOR_AQUA", golua.LNumber(16776960))
	tab.RawSetString("COLOR_BLACK", golua.LNumber(0))
	tab.RawSetString("COLOR_BLUE", golua.LNumber(16711680))
	tab.RawSetString("COLOR_DKGRAY", golua.LNumber(4210752))
	tab.RawSetString("COLOR_FUCHSIA", golua.LNumber(16711935))
	tab.RawSetString("COLOR_GRAY", golua.LNumber(8421504))
	tab.RawSetString("COLOR_GREEN", golua.LNumber(32768))
	tab.RawSetString("COLOR_LIME", golua.LNumber(65280))
	tab.RawSetString("COLOR_LTGRAY", golua.LNumber(12632256))
	tab.RawSetString("COLOR_MAROON", golua.LNumber(128))
	tab.RawSetString("COLOR_NAVY", golua.LNumber(8388608))
	tab.RawSetString("COLOR_OLIVE", golua.LNumber(32896))
	tab.RawSetString("COLOR_ORANGE", golua.LNumber(4235519))
	tab.RawSetString("COLOR_PURPLE", golua.LNumber(8388736))
	tab.RawSetString("COLOR_RED", golua.LNumber(255))
	tab.RawSetString("COLOR_SILVER", golua.LNumber(12632256))
	tab.RawSetString("COLOR_TEAL", golua.LNumber(8421376))
	tab.RawSetString("COLOR_WHITE", golua.LNumber(16777215))
	tab.RawSetString("COLOR_YELLOW", golua.LNumber(65535))
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

func spriteTable(lib *lua.Lib, state *golua.LState, name string, width, height int, parent, texgroup *golua.LTable, lg *log.Logger) *golua.LTable {
	/// @struct Sprite
	/// @prop name {string}
	/// @prop width {int}
	/// @prop height {int}
	/// @prop parent {struct<gamemaker.ResourceNode>}
	/// @prop texgroup {struct<gamemaker.ResourceNode>}
	/// @method layers(self) -> struct<gamemaker.SpriteLayers>
	/// @method frames(self) -> struct<gamemaker.SpriteFrames>
	/// @method tags(self, string...) -> self
	/// @method tag_list(self, []string) -> self
	/// @method tile(self, htile bool, vtile bool) -> self
	/// @method origin(self, int<gamemaker.SpriteOrigin>, xorigin int?, yorigin int?) -> self
	/// @method collision(self, int<gamemaker.BBOXMode>, int<gamemaker.CollMask>, struct<gamemaker.BBOX>?, tolerance int?) -> self
	/// @method premultiply_alpha(self, bool) -> self
	/// @method edge_filtering(self, bool) -> self
	/// @method dynamic_texturepage(self, bool) -> self
	/// @method nineslice(self, top int, left int, bottom int, right int) -> self
	/// @method nineslice_tilemode(self, int<gamemaker.NineSliceSlice>, int<gamemaker.NineSliceTile>) -> self
	/// @method broadcast_message(self, frame int, msg string) -> self
	/// @method playback(self, speed int, units int<gamemaker.PlaybackUnit>?) -> self

	t := state.NewTable()
	t.RawSetString("name", golua.LString(name))
	t.RawSetString("width", golua.LNumber(width))
	t.RawSetString("height", golua.LNumber(height))
	t.RawSetString("parent", parent)
	t.RawSetString("texgroup", texgroup)

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
	t.RawSetString("__playbackSpeed", golua.LNil)
	t.RawSetString("__playbackUnits", golua.LNil)

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

	resourceTags(lib, state, t)

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

	tableBuilderFunc(state, t, "playback", func(state *golua.LState, t *golua.LTable) {
		speed := state.CheckNumber(2)
		units := state.OptNumber(3, -1)

		t.RawSetString("__playbackSpeed", speed)

		if units > -1 {
			t.RawSetString("__playbackUnits", units)
		}
	})

	return t
}

func spriteBuild(t *golua.LTable, r *lua.Runner, state *golua.LState) (*yyp.Sprite, error) {
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
			r.IC.Schedule(state, img, &collection.Task[collection.ItemImage]{
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

	playback := t.RawGetString("__playbackSpeed")
	if playback.Type() == golua.LTNumber {
		sprite.Resource.Sequence.PlaybackSpeed = float64(playback.(golua.LNumber))

		units := t.RawGetString("__playbackUnits")
		if units.Type() == golua.LTNumber {
			sprite.Resource.Sequence.TimeUnits = yyp.SeqTimeUnits(units.(golua.LNumber))
		}
	}

	return sprite, nil
}

func layersParse(r *lua.Runner, lg *log.Logger, state *golua.LState, parent *golua.LTable, layer yyp.SpriteLayer, frames *golua.LTable, encoding imageutil.ImageEncoding) *golua.LTable {
	if len(layer.Frames) > 0 {
		it := layerImageTable(state, layer.Name)
		base := parent.RawGetString("__base").(*golua.LTable)
		count := base.RawGetString("__layerCount").(golua.LNumber)
		base.RawSetString("__layerCount", count+1)

		for fi, f := range layer.Frames {
			name := fmt.Sprintf("%s_%d", layer.Name, fi)

			chLog := log.NewLogger(fmt.Sprintf("image_%s", name), lg)
			lg.Append(fmt.Sprintf("child log created: image_%s", name), log.LEVEL_INFO)

			id := r.IC.AddItem(&chLog)

			frame := frames.RawGetInt(fi + 1)
			if frame.Type() != golua.LTTable {
				frame = state.NewTable()
			}
			ft := frame.(*golua.LTable)
			ft.Append(golua.LNumber(id))
			frames.RawSetInt(fi+1, ft)

			r.IC.Schedule(state, id, &collection.Task[collection.ItemImage]{
				Lib:  LIB_GAMEMAKER,
				Name: "sprite_load",
				Fn: func(i *collection.Item[collection.ItemImage]) {
					i.Self = &collection.ItemImage{
						Image:    f,
						Encoding: encoding,
						Name:     name,
						Model:    imageutil.MODEL_NRGBA,
					}
				},
			})
		}

		return it
	}

	t := layerFolderTable(state, layer.Name, parent)
	layers := t.RawGetString("__layers").(*golua.LTable)

	for _, l := range layer.Layers {
		layers.Append(layersParse(r, lg, state, t, l, frames, encoding))
	}

	t.RawSetString("__layers", layers)
	return t
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
	/// @method image(self, name string) -> self
	/// @method default(self) -> self
	/// @method folder(self, name string) -> struct<gamemaker.SpriteLayerFolder>
	/// @method back(self) -> struct<gamemaker.Sprite>

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
	/// @method image(self, name string) -> self
	/// @method default(self) -> self
	/// @method folder(self, name string) -> struct<gamemaker.SpriteLayerFolder>
	/// @method back(self) -> struct<gamemaker.SpriteLayers>

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
	/// @method add(self, []int<collection.IMAGE>) -> self
	/// @method back(self) -> struct<gamemaker.Sprite>

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

func noteTable(lib *lua.Lib, state *golua.LState, name, text string, parent golua.LValue) *golua.LTable {
	/// @struct Note
	/// @prop name {string}
	/// @prop text {string}
	/// @prop parent {struct<gamemaker.ResourceNode>}
	/// @method tags(self, string...) -> self
	/// @method tag_list(self, []string) -> self

	t := state.NewTable()
	t.RawSetString("name", golua.LString(name))
	t.RawSetString("text", golua.LString(text))
	t.RawSetString("parent", parent)

	resourceTags(lib, state, t)

	return t
}

func noteBuild(t *golua.LTable) *yyp.Note {
	name := string(t.RawGetString("name").(golua.LString))
	text := string(t.RawGetString("text").(golua.LString))
	parent := resourceNodeBuild(t.RawGetString("parent").(*golua.LTable))

	note := yyp.NewNote(name, text, parent)

	tags := t.RawGetString("__tags")
	if tags.Type() == golua.LTTable {
		tgs := tags.(*golua.LTable)
		tagList := make([]string, tgs.Len())
		for i := range tgs.Len() {
			tagList[i] = string(tgs.RawGetInt(i + 1).(golua.LString))
		}
		note.Resource.Tags = tagList
	}

	return note
}

func scriptTable(lib *lua.Lib, state *golua.LState, name, code string, parent golua.LValue) *golua.LTable {
	/// @struct Script
	/// @prop name {string}
	/// @prop code {string}
	/// @prop parent {struct<gamemaker.ResourceNode>}
	/// @method tags(self, string...) -> self
	/// @method tag_list(self, []string) -> self

	t := state.NewTable()
	t.RawSetString("name", golua.LString(name))
	t.RawSetString("code", golua.LString(code))
	t.RawSetString("parent", parent)

	resourceTags(lib, state, t)

	return t
}

func scriptBuild(t *golua.LTable) *yyp.Script {
	name := string(t.RawGetString("name").(golua.LString))
	code := string(t.RawGetString("code").(golua.LString))
	parent := resourceNodeBuild(t.RawGetString("parent").(*golua.LTable))

	script := yyp.NewScript(name, code, parent)

	tags := t.RawGetString("__tags")
	if tags.Type() == golua.LTTable {
		tgs := tags.(*golua.LTable)
		tagList := make([]string, tgs.Len())
		for i := range tgs.Len() {
			tagList[i] = string(tgs.RawGetInt(i + 1).(golua.LString))
		}
		script.Resource.Tags = tagList
	}

	return script
}

func bboxTable(state *golua.LState, top, left, bottom, right int) *golua.LTable {
	/// @struct BBOX
	/// @prop top {int}
	/// @prop left {int}
	/// @prop bottom {int}
	/// @prop right {int}

	t := state.NewTable()
	t.RawSetString("top", golua.LNumber(top))
	t.RawSetString("left", golua.LNumber(left))
	t.RawSetString("bottom", golua.LNumber(bottom))
	t.RawSetString("right", golua.LNumber(right))

	return t
}

type DataFileType int

const (
	DATAFILE_STRING DataFileType = iota
	DATAFILE_FILE
	DATAFILE_IMAGE
)

func datafileTable(state *golua.LState, name, filepath string, data golua.LValue, dataType DataFileType) *golua.LTable {
	/// @struct DataFile
	/// @prop name {string}
	/// @prop filepath {string}
	/// @prop data {string | int<collection.IMAGE>}
	/// @prop datatype {int<gamemaker.DataFileType>}

	t := state.NewTable()
	t.RawSetString("name", golua.LString(name))
	t.RawSetString("filepath", golua.LString(filepath))
	t.RawSetString("data", data)
	t.RawSetString("datatype", golua.LNumber(dataType))

	return t
}

func datafileBuild(state *golua.LState, t *golua.LTable, r *lua.Runner, lg *log.Logger) *yyp.IncludedFile {
	name := string(t.RawGetString("name").(golua.LString))
	filepath := string(t.RawGetString("filepath").(golua.LString))
	data := t.RawGetString("data")
	dataType := DataFileType(t.RawGetString("datatype").(golua.LNumber))

	var fdata []byte
	var err error

	switch dataType {
	case DATAFILE_STRING:
		fdata = []byte(string(data.(golua.LString)))
	case DATAFILE_FILE:
		fdata, err = os.ReadFile(string(data.(golua.LString)))
		if err != nil {
			state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to read file: %s", data), log.LEVEL_ERROR)), 0)
			return nil
		}
	case DATAFILE_IMAGE:
		var img image.Image
		var encoding imageutil.ImageEncoding
		<-r.IC.Schedule(state, int(data.(golua.LNumber)), &collection.Task[collection.ItemImage]{
			Lib:  LIB_GAMEMAKER,
			Name: "datafile_save",
			Fn: func(i *collection.Item[collection.ItemImage]) {
				img = i.Self.Image
				encoding = i.Self.Encoding
			},
		})

		b := byteseeker.NewByteSeeker(1000, 500)
		err := imageutil.Encode(b, img, encoding)
		if err != nil {
			state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to encode image: %s", data), log.LEVEL_ERROR)), 0)
			return nil
		}

		fdata = b.Bytes()
	}

	return yyp.NewIncludedFile(filepath, name, &fdata)
}

func resourceTags(lib *lua.Lib, state *golua.LState, t *golua.LTable) {
	t.RawSetString("__tags", golua.LNil)

	lib.BuilderFunction(state, t, "tags",
		[]lua.Arg{
			lua.ArgVariadic("tags", lua.ArrayType{Type: lua.STRING}, false),
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			tv := args["tags"].([]any)
			tags := state.NewTable()
			for _, v := range tv {
				tags.Append(golua.LString(v.(string)))
			}

			t.RawSetString("__tags", tags)
		})

	lib.BuilderFunction(state, t, "tag_list",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "tags"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			t.RawSetString("__tags", args["tags"].(*golua.LTable))
		})
}
