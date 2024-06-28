package lib

import (
	"fmt"
	"image"
	"strconv"

	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	imageutil "github.com/ArtificialLegacy/imgscal/pkg/image_util"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
)

const LIB_SPRITESHEET = "spritesheet"

func RegisterSpritesheet(r *lua.Runner, lg *log.Logger) {
	lib := lua.NewLib(LIB_SPRITESHEET, r.State, lg)

	/// @func to_frames()
	/// @arg id
	/// @arg name - will be prefixed with the img index as `I_name`
	/// @arg count
	/// @arg width
	/// @arg height
	/// @arg perRow
	/// @arg offsets - {hpixel, vpixel, hcell, vcell}
	/// @arg hsep
	/// @arg vsep
	/// @returns array of new images
	lib.CreateFunction("to_frames",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "count"},
			{Type: lua.INT, Name: "width"},
			{Type: lua.INT, Name: "height"},
			{Type: lua.INT, Name: "perRow"},
			{Type: lua.TABLE, Name: "offsets", Optional: true, Table: &[]lua.Arg{
				{Type: lua.INT, Name: "hpixel"},
				{Type: lua.INT, Name: "vpixel"},
				{Type: lua.INT, Name: "hcell"},
				{Type: lua.INT, Name: "vcell"},
			}},
			{Type: lua.INT, Name: "hsep", Optional: true},
			{Type: lua.INT, Name: "vsep", Optional: true},
		},
		func(d lua.TaskData, args map[string]any) int {
			count := args["count"].(int)
			frameSimg := []chan image.Image{}
			var encoding imageutil.ImageEncoding
			var model imageutil.ColorModel

			for i := 0; i < count; i++ {
				frameSimg = append(frameSimg, make(chan image.Image, 2))
			}

			r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					encoding = i.Self.Encoding
					model = i.Self.Model

					width := args["width"].(int)
					height := args["height"].(int)

					offsets := args["offsets"].(map[string]any)
					offsetx := offsets["hpixel"].(int) + (offsets["hcell"].(int) * width)
					offsety := offsets["vpixel"].(int) + (offsets["vcell"].(int) * height)

					topx := offsetx
					topy := offsety
					bottomx := topx + width
					bottomy := topy + height

					for ind := 0; ind < count; ind++ {
						simg := imageutil.SubImage(i.Self.Image, topx, topy, bottomx, bottomy, true)
						frameSimg[ind] <- simg

						if (ind+1)%args["perRow"].(int) == 0 {
							topx = offsetx
							bottomx = topx + width

							topy += height + args["vsep"].(int)
							bottomy = topy + height
						} else {
							topx += width + args["hsep"].(int)
							bottomx = topx + width
						}
					}
				},
				Fail: func(i *collection.Item[collection.ItemImage]) {
					for ind := 0; ind < count; ind++ {
						frameSimg[ind] <- nil
					}
				},
			})

			frames := []int{}

			for ind := 0; ind < count; ind++ {
				name := fmt.Sprintf("%d_", ind) + args["name"].(string)

				chLog := log.NewLogger(fmt.Sprintf("image_%s", name))
				chLog.Parent = lg
				lg.Append(fmt.Sprintf("child log created: image_%s", name), log.LEVEL_INFO)

				id := r.IC.AddItem(&chLog)
				frames = append(frames, id)

				r.IC.Schedule(id, &collection.Task[collection.ItemImage]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemImage]) {
						simg := <-frameSimg[ind]
						i.Self = &collection.ItemImage{
							Image:    simg,
							Encoding: encoding,
							Name:     name,
							Model:    model,
						}
					},
				})
			}

			r.State.NewTable()
			for i, f := range frames {
				r.State.PushInteger(i + 1)
				r.State.PushInteger(f)
				r.State.SetTable(-3)
			}
			return 1
		})

	/// @func from_frames()
	/// @arg ids - array of image ids
	/// @arg name
	/// @arg width
	/// @arg height
	/// @arg model
	/// @arg encoding
	/// @returns new image
	lib.CreateFunction("from_frames",
		[]lua.Arg{
			lua.ArgArray("ids", lua.ArrayType{Type: lua.INT}, false),
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "width"},
			{Type: lua.INT, Name: "height"},
			{Type: lua.INT, Name: "model"},
			{Type: lua.INT, Name: "encoding"},
		},
		func(d lua.TaskData, args map[string]any) int {
			imgs := args["ids"].(map[string]any)
			simg := make(chan *imgData, len(imgs)+1)

			width := args["width"].(int)
			height := args["height"].(int)

			for ind, v := range imgs {
				id := v.(int)
				indHere64, _ := strconv.ParseInt(ind, 10, 64)
				indHere := int(indHere64) - 1

				r.IC.Schedule(id, &collection.Task[collection.ItemImage]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemImage]) {
						finish := make(chan struct{}, 2)
						simg <- &imgData{
							Img:    i.Self.Image,
							Index:  indHere,
							Finish: finish,
						}

						<-finish
					},
				})
			}

			name := args["name"].(string)

			chLog := log.NewLogger(fmt.Sprintf("image_%s", name))
			chLog.Parent = lg
			lg.Append(fmt.Sprintf("child log created: image_%s", name), log.LEVEL_INFO)

			id := r.IC.AddItem(&chLog)

			r.IC.Schedule(id, &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					model := lua.ParseEnum(args["model"].(int), imageutil.ModelList, lib)
					encoding := lua.ParseEnum(args["encoding"].(int), imageutil.EncodingList, lib)
					i.Self = &collection.ItemImage{
						Image:    imageutil.NewImage(width*len(imgs), height, model),
						Name:     name,
						Encoding: encoding,
						Model:    model,
					}

					for range imgs {
						si := <-simg
						imageutil.Draw(i.Self.Image, si.Img, args["width"].(int)*si.Index, 0, args["width"].(int), args["height"].(int))

						si.Finish <- struct{}{}
					}
				},
				Fail: func(i *collection.Item[collection.ItemImage]) {
					for range len(simg) {
						si := <-simg
						si.Finish <- struct{}{}
					}
				},
			})

			r.State.PushInteger(id)
			return 1
		})
}

type imgData struct {
	Img    image.Image
	Index  int
	Finish chan struct{}
}
