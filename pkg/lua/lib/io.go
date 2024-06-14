package lib

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path"
	"path/filepath"

	img "github.com/ArtificialLegacy/imgscal/pkg/image"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	golua "github.com/Shopify/go-lua"
)

const LIB_IO = "io"

func RegisterIO(r *lua.Runner, lg *log.Logger) {
	lib := lua.NewLib(LIB_IO, r.State, lg)

	/// @func load_image()
	/// @arg path - the path to grab the image from
	/// @returns int - the image id
	lib.CreateFunction("load_image",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
		},
		func(state *golua.State, args map[string]any) int {
			file, err := os.Stat(args["path"].(string))
			if err != nil {
				state.PushString(lg.Append("invalid image path provided to io.load_image", log.LEVEL_ERROR))
				state.Error()
			}
			if file.IsDir() {
				state.PushString(lg.Append("cannot load a directory as an image", log.LEVEL_ERROR))
				state.Error()
			}

			id := r.IC.AddImage(file.Name())

			r.IC.Schedule(id, &img.ImageTask{
				Lib:  LIB_IO,
				Name: "load_image",
				Fn: func(i *img.Image) {
					f, err := os.Open(args["path"].(string))
					if err != nil {
						state.PushString(lg.Append("cannot open provided file", log.LEVEL_ERROR))
						state.Error()
					}
					defer f.Close()

					image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
					image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
					image.RegisterFormat("gif", "gif", gif.Decode, gif.DecodeConfig)
					image, _, err := image.Decode(f)
					if err != nil {
						state.PushString(lg.Append("provided file is an invalid image", log.LEVEL_ERROR))
						state.Error()
					}

					i.Img = image
				},
			})

			state.PushInteger(id)
			return 1
		})

	/// @func out()
	/// @arg image_id - the image id to encode and save to file.
	/// @arg path - the directory path to save the file to.
	lib.CreateFunction("out",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "path"},
		},
		func(state *golua.State, args map[string]any) int {
			_, err := os.Stat(args["path"].(string))
			if err != nil {
				os.MkdirAll(args["path"].(string), 0o666)
			}

			r.IC.Schedule(args["id"].(int), &img.ImageTask{
				Lib:  LIB_IO,
				Name: "out",
				Fn: func(i *img.Image) {
					f, err := os.OpenFile(path.Join(args["path"].(string), i.Name), os.O_CREATE, 0o666)
					if err != nil {
						state.PushString(lg.Append("cannot open provided file", log.LEVEL_ERROR))
						state.Error()
					}
					defer f.Close()

					ext := filepath.Ext(i.Name)

					switch ext {
					case ".png":
						lg.Append("image encoded as png", log.LEVEL_INFO)
						png.Encode(f, i.Img)
					case ".jpg":
						lg.Append("image encoded as jpg", log.LEVEL_INFO)
						jpeg.Encode(f, i.Img, &jpeg.Options{Quality: 100})
					case ".gif":
						lg.Append("image encoded as gif", log.LEVEL_INFO)
						gif.Encode(f, i.Img, &gif.Options{})
					default:
						state.PushString(lg.Append(fmt.Sprintf("unknown encoding used: %s", ext), log.LEVEL_ERROR))
						state.Error()
					}
				},
			})
			return 0
		})

	/// @func dir_img()
	/// @arg path - the directory path to scan for images.
	/// @returns array containing strings of each valid image in the directory.
	lib.CreateFunction("dir_img",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
		},
		func(state *golua.State, args map[string]any) int {
			f, err := os.Stat(args["path"].(string))
			if err != nil {
				state.PushString(lg.Append("invalid dir path provided to io.dir_img", log.LEVEL_ERROR))
				state.Error()
			}
			if !f.IsDir() {
				state.PushString(lg.Append("dir provided is not a directory", log.LEVEL_ERROR))
				state.Error()
			}

			files, err := os.ReadDir(args["path"].(string))
			if err != nil {
				state.PushString(lg.Append("failed to open dir", log.LEVEL_ERROR))
				state.Error()
			}

			r.State.NewTable()

			i := 1
			for _, file := range files {
				ext := filepath.Ext(file.Name())
				if ext != ".png" && ext != ".jpg" && ext != ".gif" {
					continue
				}

				lg.Append(fmt.Sprintf("found file %s with dir_img", file.Name()), log.LEVEL_INFO)

				pth := path.Join(args["path"].(string), file.Name())
				r.State.PushInteger(i)
				r.State.PushString(pth)
				r.State.SetTable(-3)
				i++
			}
			return 1
		})
}
