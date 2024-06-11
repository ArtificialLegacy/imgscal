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
	r.State.NewTable()

	r.State.PushGoFunction(func(state *golua.State) int {
		lg.Append("io.load_image called", log.LEVEL_INFO)

		result, ok := state.ToString(-1)
		if !ok {
			state.PushString(lg.Append("invalid image path provided to io.load_image", log.LEVEL_ERROR))
			state.Error()
		}

		file, err := os.Stat(result)
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
			Fn: func(i *img.Image) {
				lg.Append("io.load_image task ran", log.LEVEL_INFO)

				f, err := os.Open(result)
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

				i.Img = &image

				lg.Append("io.load_image task finished", log.LEVEL_INFO)
			},
		})

		state.PushInteger(id)
		return 1
	})
	r.State.SetField(-2, "load_image")

	r.State.PushGoFunction(func(state *golua.State) int {
		lg.Append("io.out called", log.LEVEL_INFO)

		id, ok := state.ToInteger(-2)
		if !ok {
			state.PushString(lg.Append("invalid image id provided to io.out", log.LEVEL_ERROR))
			state.Error()
		}

		outDir, ok := state.ToString(-1)
		if !ok {
			state.PushString(lg.Append("invalid outDir provided to io.out", log.LEVEL_ERROR))
			state.Error()
		}

		_, err := os.Stat(outDir)
		if err != nil {
			os.MkdirAll(outDir, 0o666)
		}

		r.IC.Schedule(id, &img.ImageTask{
			Fn: func(i *img.Image) {
				lg.Append("io.out task ran", log.LEVEL_INFO)

				f, err := os.OpenFile(path.Join(outDir, i.Name), os.O_CREATE, 0o666)
				if err != nil {
					state.PushString(lg.Append("cannot open provided file", log.LEVEL_ERROR))
					state.Error()
				}
				defer f.Close()

				ext := filepath.Ext(i.Name)

				switch ext {
				case ".png":
					lg.Append("image encoded as png", log.LEVEL_INFO)
					png.Encode(f, *i.Img)
				case ".jpg":
					lg.Append("image encoded as jpg", log.LEVEL_INFO)
					jpeg.Encode(f, *i.Img, &jpeg.Options{Quality: 100})
				case ".gif":
					lg.Append("image encoded as gif", log.LEVEL_INFO)
					gif.Encode(f, *i.Img, &gif.Options{})
				default:
					state.PushString(lg.Append(fmt.Sprintf("unknown encoding used: %s", ext), log.LEVEL_ERROR))
					state.Error()
				}

				lg.Append("io.out task finished", log.LEVEL_INFO)
			},
		})

		return 0
	})
	r.State.SetField(-2, "out")

	r.State.PushGoFunction(func(state *golua.State) int {
		lg.Append("io.dir_img called", log.LEVEL_INFO)

		dir, ok := state.ToString(-1)
		if !ok {
			state.PushString(lg.Append("invalid dir provided to dir_img", log.LEVEL_ERROR))
			state.Error()
		}

		f, err := os.Stat(dir)
		if err != nil {
			state.PushString(lg.Append("invalid dir path provided to io.dir_img", log.LEVEL_ERROR))
			state.Error()
		}
		if !f.IsDir() {
			state.PushString(lg.Append("dir provided is not a directory", log.LEVEL_ERROR))
			state.Error()
		}

		files, err := os.ReadDir(dir)
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

			pth := path.Join(dir, file.Name())
			r.State.PushInteger(i)
			r.State.PushString(pth)
			r.State.SetTable(-3)
			i++
		}
		return 1
	})
	r.State.SetField(-2, "dir_img")

	r.State.SetGlobal(LIB_IO)
}
