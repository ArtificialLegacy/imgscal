package lib

import (
	"fmt"
	"image"
	"strconv"
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
/// Modeled after the 'To Frames' functionality in GameMaker.

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
		[]lua.Arg{
			{Type: lua.INT, Name: "count"},
			{Type: lua.INT, Name: "width"},
			{Type: lua.INT, Name: "height"},
			{Type: lua.INT, Name: "perRow"},
			{Type: lua.TABLE, Name: "offsets", Optional: true, Table: &[]lua.Arg{
				{Type: lua.INT, Name: "hpixel", Optional: true},
				{Type: lua.INT, Name: "vpixel", Optional: true},
				{Type: lua.INT, Name: "hcell", Optional: true},
				{Type: lua.INT, Name: "vcell", Optional: true},
				{Type: lua.INT, Name: "index", Optional: true},
			}},
			{Type: lua.INT, Name: "hsep", Optional: true},
			{Type: lua.INT, Name: "vsep", Optional: true},
		},
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
		[]lua.Arg{
			{Type: lua.INT, Name: "hpixel"},
			{Type: lua.INT, Name: "vpixel"},
			{Type: lua.INT, Name: "hcell"},
			{Type: lua.INT, Name: "vcell"},
			{Type: lua.INT, Name: "index"},
		},
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
			t := state.NewTable()

			t.RawSetString("hpixel", golua.LNumber(args["hpixel"].(int)))
			t.RawSetString("vpixel", golua.LNumber(args["vpixel"].(int)))
			t.RawSetString("hcell", golua.LNumber(0))
			t.RawSetString("vcell", golua.LNumber(0))
			t.RawSetString("index", golua.LNumber(0))

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
			t := state.NewTable()

			t.RawSetString("hpixel", golua.LNumber(0))
			t.RawSetString("vpixel", golua.LNumber(0))
			t.RawSetString("hcell", golua.LNumber(args["hcell"].(int)))
			t.RawSetString("vcell", golua.LNumber(args["vcell"].(int)))
			t.RawSetString("index", golua.LNumber(0))

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
			t := state.NewTable()

			t.RawSetString("hpixel", golua.LNumber(0))
			t.RawSetString("vpixel", golua.LNumber(0))
			t.RawSetString("hcell", golua.LNumber(0))
			t.RawSetString("vcell", golua.LNumber(0))
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
			{Type: lua.TABLE, Name: "sheet", Table: &[]lua.Arg{
				{Type: lua.INT, Name: "count"},
				{Type: lua.INT, Name: "width"},
				{Type: lua.INT, Name: "height"},
				{Type: lua.INT, Name: "perRow"},
				{Type: lua.TABLE, Name: "offsets", Optional: true, Table: &[]lua.Arg{
					{Type: lua.INT, Name: "hpixel", Optional: true},
					{Type: lua.INT, Name: "vpixel", Optional: true},
					{Type: lua.INT, Name: "hcell", Optional: true},
					{Type: lua.INT, Name: "vcell", Optional: true},
					{Type: lua.INT, Name: "index", Optional: true},
				}},
				{Type: lua.INT, Name: "hsep", Optional: true},
				{Type: lua.INT, Name: "vsep", Optional: true},
			}},
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

				r.IC.Schedule(id, &collection.Task[collection.ItemImage]{
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

			r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					width := sheet["width"].(int)
					height := sheet["height"].(int)

					perRow := sheet["perRow"].(int)

					offsets := sheet["offsets"].(map[string]any)
					hpixel := offsets["hpixel"].(int)
					vpixel := offsets["vpixel"].(int)
					hcell := offsets["hcell"].(int)
					vcell := offsets["vcell"].(int)
					index := offsets["index"].(int)

					hsep := sheet["hsep"].(int)
					vsep := sheet["vsep"].(int)

					imgs := imageutil.SpritesheetToFrames(i.Self.Image, !args["nocopy"].(bool), count, width, height, perRow, hpixel, vpixel, hcell, vcell, index, hsep, vsep)

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
			{Type: lua.TABLE, Name: "sheet", Table: &[]lua.Arg{
				{Type: lua.INT, Name: "count"},
				{Type: lua.INT, Name: "width"},
				{Type: lua.INT, Name: "height"},
				{Type: lua.INT, Name: "perRow"},
				{Type: lua.TABLE, Name: "offsets", Optional: true, Table: &[]lua.Arg{
					{Type: lua.INT, Name: "hpixel", Optional: true},
					{Type: lua.INT, Name: "vpixel", Optional: true},
					{Type: lua.INT, Name: "hcell", Optional: true},
					{Type: lua.INT, Name: "vcell", Optional: true},
					{Type: lua.INT, Name: "index", Optional: true},
				}},
				{Type: lua.INT, Name: "hsep", Optional: true},
				{Type: lua.INT, Name: "vsep", Optional: true},
			}},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			imgs := args["ids"].(map[string]any)
			wg := sync.WaitGroup{}
			finish := make(chan struct{})

			sheet := args["sheet"].(map[string]any)

			imgList := make([]image.Image, len(imgs))

			count := sheet["count"].(int)
			width := sheet["width"].(int)
			height := sheet["height"].(int)

			wg.Add(len(imgs))
			for ind := range len(imgs) {
				id := imgs[strconv.Itoa(ind+1)].(int)

				r.IC.Schedule(id, &collection.Task[collection.ItemImage]{
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

			r.IC.Schedule(id, &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					model := lua.ParseEnum(args["model"].(int), imageutil.ModelList, lib)
					encoding := lua.ParseEnum(args["encoding"].(int), imageutil.EncodingList, lib)

					perRow := sheet["perRow"].(int)

					offsets := sheet["offsets"].(map[string]any)
					hpixel := offsets["hpixel"].(int)
					vpixel := offsets["vpixel"].(int)
					hcell := offsets["hcell"].(int)
					vcell := offsets["vcell"].(int)
					index := offsets["index"].(int)

					hsep := sheet["hsep"].(int)
					vsep := sheet["vsep"].(int)

					wg.Wait()

					img := imageutil.FramesToSpritesheet(imgList, model, count, width, height, perRow, hpixel, vpixel, hcell, vcell, index, hsep, vsep)
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

	/// @func extract(id, name, spritesheet_in, spritesheet_out) -> int<collection.IMAGE>
	/// @arg id {int<collection.IMAGE>}
	/// @arg name {string}
	/// @arg spritesheet_in {struct<spritesheet.Spritesheet>} - The spritesheet related to the source image.
	/// @arg spritesheet_out {struct<spritesheet.Spritesheet>} - The spritesheet related to the returned image.
	/// @returns {int<collection.IMAGE>}
	/// @desc
	/// Note it is more efficient to exclude frames using index and count
	/// from spritesheet_in than from spritesheet_out. This prevents them from being sub-imaged completely.
	lib.CreateFunction(tab, "extract",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "name"},
			{Type: lua.TABLE, Name: "sheetin", Table: &[]lua.Arg{
				{Type: lua.INT, Name: "count"},
				{Type: lua.INT, Name: "width"},
				{Type: lua.INT, Name: "height"},
				{Type: lua.INT, Name: "perRow"},
				{Type: lua.TABLE, Name: "offsets", Optional: true, Table: &[]lua.Arg{
					{Type: lua.INT, Name: "hpixel", Optional: true},
					{Type: lua.INT, Name: "vpixel", Optional: true},
					{Type: lua.INT, Name: "hcell", Optional: true},
					{Type: lua.INT, Name: "vcell", Optional: true},
					{Type: lua.INT, Name: "index", Optional: true},
				}},
				{Type: lua.INT, Name: "hsep", Optional: true},
				{Type: lua.INT, Name: "vsep", Optional: true},
			}},
			{Type: lua.TABLE, Name: "sheetout", Table: &[]lua.Arg{
				{Type: lua.INT, Name: "count"},
				{Type: lua.INT, Name: "width"},
				{Type: lua.INT, Name: "height"},
				{Type: lua.INT, Name: "perRow"},
				{Type: lua.TABLE, Name: "offsets", Optional: true, Table: &[]lua.Arg{
					{Type: lua.INT, Name: "hpixel", Optional: true},
					{Type: lua.INT, Name: "vpixel", Optional: true},
					{Type: lua.INT, Name: "hcell", Optional: true},
					{Type: lua.INT, Name: "vcell", Optional: true},
					{Type: lua.INT, Name: "index", Optional: true},
				}},
				{Type: lua.INT, Name: "hsep", Optional: true},
				{Type: lua.INT, Name: "vsep", Optional: true},
			}},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var encoding imageutil.ImageEncoding
			var model imageutil.ColorModel

			var imgs []image.Image
			ready := make(chan struct{})
			finish := make(chan struct{})

			r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					sheet := args["sheetin"].(map[string]any)

					count := sheet["count"].(int)
					width := sheet["width"].(int)
					height := sheet["height"].(int)
					perRow := sheet["perRow"].(int)

					offsets := sheet["offsets"].(map[string]any)
					hpixel := offsets["hpixel"].(int)
					vpixel := offsets["vpixel"].(int)
					hcell := offsets["hcell"].(int)
					vcell := offsets["vcell"].(int)
					index := offsets["index"].(int)

					hsep := sheet["hsep"].(int)
					vsep := sheet["vsep"].(int)

					imgs = imageutil.SpritesheetToFrames(i.Self.Image, false, count, width, height, perRow, hpixel, vpixel, hcell, vcell, index, hsep, vsep)

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

			r.IC.Schedule(id, &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					sheet := args["sheetout"].(map[string]any)

					count := sheet["count"].(int)
					width := sheet["width"].(int)
					height := sheet["height"].(int)
					perRow := sheet["perRow"].(int)

					offsets := sheet["offsets"].(map[string]any)
					hpixel := offsets["hpixel"].(int)
					vpixel := offsets["vpixel"].(int)
					hcell := offsets["hcell"].(int)
					vcell := offsets["vcell"].(int)
					index := offsets["index"].(int)

					hsep := sheet["hsep"].(int)
					vsep := sheet["vsep"].(int)

					<-ready

					img := imageutil.FramesToSpritesheet(imgs, model, count, width, height, perRow, hpixel, vpixel, hcell, vcell, index, hsep, vsep)
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
		},
	)
}
