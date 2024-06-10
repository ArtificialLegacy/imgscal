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

	"github.com/ArtificialLegacy/imgscal/pkg/cli"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/workflow"
	"github.com/Shopify/go-lua"
)

const LIB_IMGSCAL = "imgscal"

func RegisterImgscal(state *lua.State, data *workflow.WorkflowData, lg *log.Logger) {
	state.NewTable()

	state.PushGoFunction(func(state *lua.State) int {
		lg.Append("imgscal.name called", log.LEVEL_INFO)

		id, ok := state.ToInteger(-2)
		if !ok {
			state.PushString(lg.Append("invalid image id provided to name", log.LEVEL_ERROR))
			state.Error()
		}

		name, ok := state.ToString(-1)
		if !ok {
			state.PushString(lg.Append("invalid image name provided to name", log.LEVEL_ERROR))
			state.Error()
		}

		go func() {
			i, err := data.IC.Image(id)
			if err != nil {
				state.PushString(lg.Append("invalid image provided to name", log.LEVEL_ERROR))
				state.Error()
			}

			i.Name = name
			i.Mutex.Unlock()
		}()

		return 0
	})
	state.SetField(-2, "name")

	state.PushGoFunction(func(state *lua.State) int {
		lg.Append("imgscal.prompt called", log.LEVEL_INFO)

		question, ok := state.ToString(-1)
		if !ok {
			state.PushString(lg.Append("invalid question provided to prompt_file", log.LEVEL_ERROR))
			state.Error()
		}

		result, err := cli.Question(question, cli.QuestionOptions{})
		if err != nil {
			state.PushString(lg.Append("invalid answer provided to prompt", log.LEVEL_ERROR))
			state.Error()
		}

		state.PushString(result)
		return 1
	})
	state.SetField(-2, "prompt")

	state.PushGoFunction(func(state *lua.State) int {
		lg.Append("imgscal.load_image called", log.LEVEL_INFO)

		result, ok := state.ToString(-1)
		if !ok {
			state.PushString(lg.Append("invalid image path provided to load_image", log.LEVEL_ERROR))
			state.Error()
		}

		file, err := os.Stat(result)
		if err != nil {
			state.PushString(lg.Append("invalid image path provided to load_image", log.LEVEL_ERROR))
			state.Error()
		}

		i, id := data.IC.AddImage(file.Name())
		i.Mutex.Lock()

		go func() {
			if file.IsDir() {
				state.PushString(lg.Append("cannot load directory as image", log.LEVEL_ERROR))
				state.Error()
			}

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
			i.Mutex.Unlock()
		}()

		state.PushInteger(id)
		return 1
	})
	state.SetField(-2, "load_image")

	state.PushGoFunction(func(state *lua.State) int {
		lg.Append("imgscal.out called", log.LEVEL_INFO)

		id, ok := state.ToInteger(-2)
		if !ok {
			state.PushString(lg.Append("invalid image id provided to out", log.LEVEL_ERROR))
			state.Error()
		}

		outDir, ok := state.ToString(-1)
		if !ok {
			state.PushString(lg.Append("invalid outDir provided to out", log.LEVEL_ERROR))
			state.Error()
		}

		_, err := os.Stat(outDir)
		if err != nil {
			os.MkdirAll(outDir, 0o666)
		}

		go func() {
			i, err := data.IC.Image(id)
			if err != nil {
				state.PushString(lg.Append("invalid image provided to out", log.LEVEL_ERROR))
				state.Error()
			}

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

			i.Mutex.Unlock()
		}()

		return 0
	})
	state.SetField(-2, "out")

	state.SetGlobal(LIB_IMGSCAL)
}
