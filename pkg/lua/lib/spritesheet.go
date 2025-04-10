package lib

import (
	"fmt"
	"image"
	"sync"

	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	imageutil "github.com/ArtificialLegacy/imgscal/pkg/image_util"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	golua "github.com/yuin/gopher-lua"
)

const LIB_SPRITESHEET = "spritesheet"

/// @lib SpriteSheet
/// @import spritesheet
/// @desc
/// Library to provide support for creating and splitting spritesheets.
/// @section
/// Modeled after the 'To Frames' functionality in GameMaker.

var offsets = []lua.Arg{
	{Type: lua.INT, Name: "hpixel", Optional: true},
	{Type: lua.INT, Name: "vpixel", Optional: true},
	{Type: lua.INT, Name: "hcell", Optional: true},
	{Type: lua.INT, Name: "vcell", Optional: true},
	{Type: lua.INT, Name: "index", Optional: true},
}

var sheet = []lua.Arg{
	{Type: lua.INT, Name: "count"},
	{Type: lua.INT, Name: "width"},
	{Type: lua.INT, Name: "height"},
	{Type: lua.INT, Name: "perRow"},
	{Type: lua.TABLE, Name: "offsets", Optional: true, Table: &offsets},
	{Type: lua.INT, Name: "hsep", Optional: true},
	{Type: lua.INT, Name: "vsep", Optional: true},
}

func RegisterSpritesheet(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_SPRITESHEET, r, r.State, lg)

	/// @func sheet(count, width, height, perRow, offsets?, hsep?, vsep?) -> struct<spritesheet.Spritesheet>
	/// @arg count {int}
	/// @arg width {int}
	/// @arg height {int}
	/// @arg perRow {int}
	/// @arg? offsets {struct<spritesheet.Offset>}
	/// @arg? hsep {int}
	/// @arg? vsep {int}
	/// @returns {struct<spritesheet.Spritesheet>}
	lib.CreateFunction(tab, "sheet",
		sheet,
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct Spritesheet
			/// @prop count {int}
			/// @prop width {int}
			/// @prop height {int}
			/// @prop perRow {int}
			/// @prop offsets {struct<spritesheet.Offset>}
			/// @prop hsep {int}
			/// @prop vsep {int}

			t := state.NewTable()

			t.RawSetString("count", golua.LNumber(args["count"].(int)))
			t.RawSetString("width", golua.LNumber(args["width"].(int)))
			t.RawSetString("height", golua.LNumber(args["height"].(int)))
			t.RawSetString("perRow", golua.LNumber(args["perRow"].(int)))

			offsets := args["offsets"].(map[string]any)
			ot := state.NewTable()
			ot.RawSetString("hpixel", golua.LNumber(offsets["hpixel"].(int)))
			ot.RawSetString("vpixel", golua.LNumber(offsets["vpixel"].(int)))
			ot.RawSetString("hcell", golua.LNumber(offsets["hcell"].(int)))
			ot.RawSetString("vcell", golua.LNumber(offsets["vcell"].(int)))
			ot.RawSetString("index", golua.LNumber(offsets["index"].(int)))
			t.RawSetString("offsets", ot)

			t.RawSetString("hsep", golua.LNumber(args["hsep"].(int)))
			t.RawSetString("vsep", golua.LNumber(args["vsep"].(int)))

			state.Push(t)
			return 1
		})

	/// @func offset(hpixel, vpixel, hcell, vcell, index) -> struct<spritesheet.Offset>
	/// @arg hpixel {int}
	/// @arg vpixel {int}
	/// @arg hcell {int}
	/// @arg vcell {int}
	/// @arg index {int}
	/// @returns {struct<spritesheet.Offset>}
	lib.CreateFunction(tab, "offset",
		offsets,
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct Offset
			/// @prop hpixel {int}
			/// @prop vpixel {int}
			/// @prop hcell {int}
			/// @prop vcell {int}
			/// @prop index {int}

			t := state.NewTable()

			t.RawSetString("hpixel", golua.LNumber(args["hpixel"].(int)))
			t.RawSetString("vpixel", golua.LNumber(args["vpixel"].(int)))
			t.RawSetString("hcell", golua.LNumber(args["hcell"].(int)))
			t.RawSetString("vcell", golua.LNumber(args["vcell"].(int)))
			t.RawSetString("index", golua.LNumber(args["index"].(int)))

			state.Push(t)
			return 1
		})

	/// @func offset_pixel(hpixel, vpixel) -> struct<spritesheet.Offset>
	/// @arg hpixel {int}
	/// @arg vpixel {int}
	/// @returns {struct<spritesheet.Offset>}
	lib.CreateFunction(tab, "offset_pixel",
		[]lua.Arg{
			{Type: lua.INT, Name: "hpixel"},
			{Type: lua.INT, Name: "vpixel"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := spritesheetOffsetTable(state)

			t.RawSetString("hpixel", golua.LNumber(args["hpixel"].(int)))
			t.RawSetString("vpixel", golua.LNumber(args["vpixel"].(int)))

			state.Push(t)
			return 1
		})

	/// @func offset_cell(hcell, vcell) -> struct<spritesheet.Offset>
	/// @arg hcell {int}
	/// @arg vcell {int}
	/// @returns {struct<spritesheet.Offset>}
	lib.CreateFunction(tab, "offset_cell",
		[]lua.Arg{
			{Type: lua.INT, Name: "hcell"},
			{Type: lua.INT, Name: "vcell"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := spritesheetOffsetTable(state)

			t.RawSetString("hcell", golua.LNumber(args["hcell"].(int)))
			t.RawSetString("vcell", golua.LNumber(args["vcell"].(int)))

			state.Push(t)
			return 1
		})

	/// @func offset_index(index) -> struct<spritesheet.Offset>
	/// @arg index {int}
	/// @returns {struct<spritesheet.Offset>}
	lib.CreateFunction(tab, "offset_index",
		[]lua.Arg{
			{Type: lua.INT, Name: "index"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := spritesheetOffsetTable(state)

			t.RawSetString("index", golua.LNumber(args["index"].(int)))

			state.Push(t)
			return 1
		})

	/// @func to_frames(id, name, spritesheet, nocopy?) -> []int<collection.IMAGE>
	/// @arg id {int<collection.IMAGE>}
	/// @arg name {string} - Will be prefixed using the frame index as '%d_name'.
	/// @arg spritesheet {struct<spritesheet.Spritesheet>}
	/// @arg? nocopy {bool}
	/// @returns {[]int<collection.IMAGE>}
	lib.CreateFunction(tab, "to_frames",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "name"},
			{Type: lua.TABLE, Name: "sheet", Table: &sheet},
			{Type: lua.BOOL, Name: "nocopy", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			name := args["name"].(string)

			sheet := args["sheet"].(map[string]any)
			count := sheet["count"].(int)
			frames := make([]int, count)
			frameChannels := make([]chan image.Image, count)
			var encoding imageutil.ImageEncoding
			var model imageutil.ColorModel

			for ind := range count {
				frameName := fmt.Sprintf("%d_", ind) + name

				chLog := log.NewLogger(fmt.Sprintf("image_%s", frameName), lg)
				lg.Append(fmt.Sprintf("child log created: image_%s", frameName), log.LEVEL_INFO)

				id := r.IC.AddItem(&chLog)
				frames[ind] = id
				frameChannels[ind] = make(chan image.Image)

				r.IC.Schedule(state, id, &collection.Task[collection.ItemImage]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemImage]) {
						img := <-frameChannels[ind]
						i.Self = &collection.ItemImage{
							Image:    img,
							Encoding: encoding,
							Name:     frameName,
							Model:    model,
						}
					},
				})
			}

			r.IC.Schedule(state, args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					imgs := imageutil.SpritesheetToFramesTable(i.Self.Image, !args["nocopy"].(bool), sheet)

					encoding = i.Self.Encoding
					model = i.Self.Model

					for fi, img := range imgs {
						frameChannels[fi] <- img
					}
				},
			})

			t := state.NewTable()
			for i, f := range frames {
				t.RawSetInt(i+1, golua.LNumber(f))
			}

			state.Push(t)
			return 1
		})

	/// @func to_frames_cached(id, spritesheet, nocopy?) -> []int<collection.CRATE_CACHEDIMAGE>
	/// @arg id {int<collection.IMAGE>}
	/// @arg spritesheet {struct<spritesheet.Spritesheet>}
	/// @arg? nocopy {bool}
	/// @returns {[]int<collection.CRATE_CACHEDIMAGE>}
	/// @blocking
	lib.CreateFunction(tab, "to_frames_cached",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.TABLE, Name: "sheet", Table: &sheet},
			{Type: lua.BOOL, Name: "nocopy", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			ids := state.NewTable()

			<-r.IC.Schedule(state, args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					sheet := args["sheet"].(map[string]any)

					imgs := imageutil.SpritesheetToFramesTable(i.Self.Image, !args["nocopy"].(bool), sheet)

					for fi, img := range imgs {
						id := r.CR_CIM.Add(&collection.CachedImageItem{
							Model: i.Self.Model,
							Image: img,
						})

						ids.RawSetInt(fi+1, golua.LNumber(id))
					}
				},
			})

			state.Push(ids)
			return 1
		})

	/// @func to_frames_into_cached(id, ids, spritesheet, nocopy?)
	/// @arg id {int<collection.IMAGE>}
	/// @arg ids {[]int<collection.CRATE_CACHEDIMAGE>} - Must be the same length as the amount of frames in the spritesheet.
	/// @arg spritesheet {struct<spritesheet.Spritesheet>}
	/// @arg? nocopy {bool}
	/// @returns {[]int<collection.CRATE_CACHEDIMAGE>}
	/// @blocking
	lib.CreateFunction(tab, "to_frames_into_cached",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			lua.ArgArray("ids", lua.ArrayType{Type: lua.INT}, false),
			{Type: lua.TABLE, Name: "sheet", Table: &sheet},
			{Type: lua.BOOL, Name: "nocopy", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			<-r.IC.Schedule(state, args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					sheet := args["sheet"].(map[string]any)

					imgs := imageutil.SpritesheetToFramesTable(i.Self.Image, !args["nocopy"].(bool), sheet)

					ids := args["ids"].([]any)
					if len(imgs) > len(ids) {
						lua.Error(state, lg.Appendf("not enough cache ids: expected=%d, got=%d", log.LEVEL_ERROR, len(imgs), len(ids)))
					}

					for fi, img := range imgs {
						id := ids[fi].(int)

						citem, err := r.CR_CIM.Item(id)
						if err != nil {
							lua.Error(state, lg.Appendf("failed to get cached image: %s", log.LEVEL_ERROR, err))
						}

						citem.Image = img
						citem.Model = i.Self.Model
					}
				},
			})

			return 0
		})

	/// @func from_frames(ids, name, model, encoding, spritesheet) -> int<collection.IMAGE>
	/// @arg ids {[]int<collection.IMAGE>}
	/// @arg name {string}
	/// @arg model {int<image.ColorModel>}
	/// @arg encoding {int<image.Encoding>}
	/// @arg spritesheet {struct<spritesheet.Spritesheet>}
	/// @returns {int<collection.IMAGE>}
	lib.CreateFunction(tab, "from_frames",
		[]lua.Arg{
			lua.ArgArray("ids", lua.ArrayType{Type: lua.INT}, false),
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "model"},
			{Type: lua.INT, Name: "encoding"},
			{Type: lua.TABLE, Name: "sheet", Table: &sheet},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			imgs := args["ids"].([]any)
			wg := sync.WaitGroup{}
			finish := make(chan struct{})

			sheet := args["sheet"].(map[string]any)

			imgList := make([]image.Image, len(imgs))

			wg.Add(len(imgs))
			for ind := range len(imgs) {
				id := imgs[ind].(int)

				r.IC.Schedule(state, id, &collection.Task[collection.ItemImage]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemImage]) {
						imgList[ind] = i.Self.Image
						wg.Done()
						<-finish
					},
				})
			}

			name := args["name"].(string)

			chLog := log.NewLogger(fmt.Sprintf("image_%s", name), lg)
			lg.Append(fmt.Sprintf("child log created: image_%s", name), log.LEVEL_INFO)

			id := r.IC.AddItem(&chLog)

			r.IC.Schedule(state, id, &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					model := lua.ParseEnum(args["model"].(int), imageutil.ModelList, lib)
					encoding := lua.ParseEnum(args["encoding"].(int), imageutil.EncodingList, lib)

					wg.Wait()

					img := imageutil.FramesToSpritesheetTable(imgList, model, sheet)
					i.Self = &collection.ItemImage{
						Image:    img,
						Name:     name,
						Encoding: encoding,
						Model:    model,
					}

					close(finish)
				},
				Fail: func(i *collection.Item[collection.ItemImage]) {
					close(finish)
				},
			})

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func from_frames_into(ids, id, model, spritesheet)
	/// @arg ids {[]int<collection.IMAGE>}
	/// @arg id {int<collection.IMAGE>}
	/// @arg model {int<image.ColorModel>}
	/// @arg spritesheet {struct<spritesheet.Spritesheet>}
	lib.CreateFunction(tab, "from_frames_into",
		[]lua.Arg{
			lua.ArgArray("ids", lua.ArrayType{Type: lua.INT}, false),
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "model"},
			{Type: lua.TABLE, Name: "sheet", Table: &sheet},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			imgs := args["ids"].([]any)
			wg := sync.WaitGroup{}
			finish := make(chan struct{})

			sheet := args["sheet"].(map[string]any)

			imgList := make([]image.Image, len(imgs))

			wg.Add(len(imgs))
			for ind := range len(imgs) {
				id := imgs[ind].(int)

				r.IC.Schedule(state, id, &collection.Task[collection.ItemImage]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemImage]) {
						imgList[ind] = i.Self.Image
						wg.Done()
						<-finish
					},
				})
			}

			r.IC.Schedule(state, args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					model := lua.ParseEnum(args["model"].(int), imageutil.ModelList, lib)

					wg.Wait()

					img := imageutil.FramesToSpritesheetTable(imgList, model, sheet)
					i.Self = &collection.ItemImage{
						Image:    img,
						Name:     i.Self.Name,
						Encoding: i.Self.Encoding,
						Model:    model,
					}

					close(finish)
				},
				Fail: func(i *collection.Item[collection.ItemImage]) {
					close(finish)
				},
			})

			return 0
		})

	/// @func extract(id, name, spritesheet_in, spritesheet_out, nocopy?) -> int<collection.IMAGE>
	/// @arg id {int<collection.IMAGE>}
	/// @arg name {string}
	/// @arg spritesheet_in {struct<spritesheet.Spritesheet>} - The spritesheet related to the source image.
	/// @arg spritesheet_out {struct<spritesheet.Spritesheet>} - The spritesheet related to the returned image.
	/// @arg? nocopy {bool}
	/// @returns {int<collection.IMAGE>}
	/// @desc
	/// Note it is more efficient to exclude frames using index and count
	/// from spritesheet_in than from spritesheet_out. This prevents them from being sub-imaged completely.
	lib.CreateFunction(tab, "extract",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "name"},
			{Type: lua.TABLE, Name: "sheetin", Table: &sheet},
			{Type: lua.TABLE, Name: "sheetout", Table: &sheet},
			{Type: lua.BOOL, Name: "nocopy", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var encoding imageutil.ImageEncoding
			var model imageutil.ColorModel

			var imgs []image.Image
			ready := make(chan struct{})
			finish := make(chan struct{})

			r.IC.Schedule(state, args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					sheet := args["sheetin"].(map[string]any)

					imgs = imageutil.SpritesheetToFramesTable(i.Self.Image, args["nocopy"].(bool), sheet)

					encoding = i.Self.Encoding
					model = i.Self.Model

					ready <- struct{}{}
					<-finish
				},
				Fail: func(i *collection.Item[collection.ItemImage]) {
					ready <- struct{}{}
				},
			})

			name := args["name"].(string)

			chLog := log.NewLogger(fmt.Sprintf("image_%s", name), lg)
			lg.Append(fmt.Sprintf("child log created: image_%s", name), log.LEVEL_INFO)

			id := r.IC.AddItem(&chLog)

			r.IC.Schedule(state, id, &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					sheet := args["sheetout"].(map[string]any)

					<-ready

					img := imageutil.FramesToSpritesheetTable(imgs, model, sheet)
					i.Self = &collection.ItemImage{
						Image:    img,
						Name:     name,
						Encoding: encoding,
						Model:    model,
					}

					finish <- struct{}{}
				},
				Fail: func(i *collection.Item[collection.ItemImage]) {
					finish <- struct{}{}
				},
			})

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func extract_into(id, idinto, spritesheet_in, spritesheet_out, nocopy?)
	/// @arg id {int<collection.IMAGE>}
	/// @arg idinto {int<collection.IMAGE>}
	/// @arg spritesheet_in {struct<spritesheet.Spritesheet>} - The spritesheet related to the source image.
	/// @arg spritesheet_out {struct<spritesheet.Spritesheet>} - The spritesheet related to the returned image.
	/// @arg? nocopy {bool}
	/// @desc
	/// Note it is more efficient to exclude frames using index and count
	/// from spritesheet_in than from spritesheet_out. This prevents them from being sub-imaged completely.
	lib.CreateFunction(tab, "extract_into",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "idinto"},
			{Type: lua.TABLE, Name: "sheetin", Table: &sheet},
			{Type: lua.TABLE, Name: "sheetout", Table: &sheet},
			{Type: lua.BOOL, Name: "nocopy", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var model imageutil.ColorModel
			var imgs []image.Image

			r.IC.SchedulePipe(state, args["id"].(int), args["idinto"].(int),
				&collection.Task[collection.ItemImage]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemImage]) {
						sheet := args["sheetin"].(map[string]any)

						imgs = imageutil.SpritesheetToFramesTable(i.Self.Image, args["nocopy"].(bool), sheet)

						model = i.Self.Model
					},
				},
				&collection.Task[collection.ItemImage]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemImage]) {
						sheet := args["sheetout"].(map[string]any)

						img := imageutil.FramesToSpritesheetTable(imgs, model, sheet)
						i.Self = &collection.ItemImage{
							Image:    img,
							Name:     i.Self.Name,
							Encoding: i.Self.Encoding,
							Model:    model,
						}
					},
				})

			return 0
		})
}

func spritesheetOffsetTable(state *golua.LState) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("hpixel", golua.LNumber(0))
	t.RawSetString("vpixel", golua.LNumber(0))
	t.RawSetString("hcell", golua.LNumber(0))
	t.RawSetString("vcell", golua.LNumber(0))
	t.RawSetString("index", golua.LNumber(0))

	return t
}
