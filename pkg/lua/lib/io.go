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

	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
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
		func(d lua.TaskData, args map[string]any) int {
			file, err := os.Stat(args["path"].(string))
			if err != nil {
				r.State.PushString(lg.Append("invalid image path provided to io.load_image", log.LEVEL_ERROR))
				r.State.Error()
			}
			if file.IsDir() {
				r.State.PushString(lg.Append("cannot load a directory as an image", log.LEVEL_ERROR))
				r.State.Error()
			}

			chLog := log.NewLogger(fmt.Sprintf("image_%s", file.Name()))
			chLog.Parent = lg
			lg.Append(fmt.Sprintf("child log created: image_%s", file.Name()), log.LEVEL_INFO)

			id := r.IC.AddItem(file.Name(), &chLog)

			r.IC.Schedule(id, &collection.Task[image.Image]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[image.Image]) {
					f, err := os.Open(args["path"].(string))
					if err != nil {
						r.State.PushString(lg.Append("cannot open provided file", log.LEVEL_ERROR))
						r.State.Error()
					}
					defer f.Close()

					image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
					image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
					image.RegisterFormat("gif", "gif", gif.Decode, gif.DecodeConfig)
					image, _, err := image.Decode(f)
					if err != nil {
						r.State.PushString(lg.Append("provided file is an invalid image", log.LEVEL_ERROR))
						r.State.Error()
					}

					i.Self = &image
				},
			})

			r.State.PushInteger(id)
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
		func(d lua.TaskData, args map[string]any) int {
			_, err := os.Stat(args["path"].(string))
			if err != nil {
				os.MkdirAll(args["path"].(string), 0o666)
			}

			r.IC.Schedule(args["id"].(int), &collection.Task[image.Image]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[image.Image]) {
					f, err := os.OpenFile(path.Join(args["path"].(string), i.Name), os.O_CREATE, 0o666)
					if err != nil {
						r.State.PushString(lg.Append("cannot open provided file", log.LEVEL_ERROR))
						r.State.Error()
					}
					defer f.Close()

					ext := filepath.Ext(i.Name)

					switch ext {
					case ".png":
						lg.Append("image encoded as png", log.LEVEL_INFO)
						png.Encode(f, *i.Self)
					case ".jpg":
						lg.Append("image encoded as jpg", log.LEVEL_INFO)
						jpeg.Encode(f, *i.Self, &jpeg.Options{Quality: 100})
					case ".gif":
						lg.Append("image encoded as gif", log.LEVEL_INFO)
						gif.Encode(f, *i.Self, &gif.Options{})
					default:
						r.State.PushString(lg.Append(fmt.Sprintf("unknown encoding used: %s", ext), log.LEVEL_ERROR))
						r.State.Error()
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
		func(d lua.TaskData, args map[string]any) int {
			f, err := os.Stat(args["path"].(string))
			if err != nil {
				r.State.PushString(lg.Append("invalid dir path provided to io.dir_img", log.LEVEL_ERROR))
				r.State.Error()
			}
			if !f.IsDir() {
				r.State.PushString(lg.Append("dir provided is not a directory", log.LEVEL_ERROR))
				r.State.Error()
			}

			files, err := os.ReadDir(args["path"].(string))
			if err != nil {
				r.State.PushString(lg.Append("failed to open dir", log.LEVEL_ERROR))
				r.State.Error()
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
