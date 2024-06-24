package lib

import (
	"fmt"
	"image"

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

			for i := 0; i < count; i++ {
				frameSimg = append(frameSimg, make(chan image.Image, 2))
			}

			r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					encoding = i.Self.Encoding

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
						simg := imageutil.SubImage(i.Self.Image, topx, topy, bottomx, bottomy)
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

				indHere := ind

				r.IC.Schedule(id, &collection.Task[collection.ItemImage]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemImage]) {
						simg := <-frameSimg[indHere]
						i.Self = &collection.ItemImage{
							Image:    simg,
							Encoding: encoding,
							Name:     name,
						}
					},
				})
			}

			r.State.NewTable()
			for i, f := range frames {
				r.State.PushInteger(i)
				r.State.PushInteger(f)
				r.State.SetTable(-3)
			}
			return 1
		})
}
