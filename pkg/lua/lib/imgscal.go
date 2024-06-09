package lib

import (
	"image"
	"image/png"
	"os"
	"path"
	"path/filepath"

	"github.com/ArtificialLegacy/imgscal/pkg/cli"
	"github.com/ArtificialLegacy/imgscal/pkg/workflow"
	"github.com/Shopify/go-lua"
)

const LIB_IMGSCAL = "imgscal"

func RegisterImgscal(state *lua.State, data *workflow.WorkflowData) {
	state.NewTable()

	state.PushGoFunction(func(state *lua.State) int {
		id, ok := state.ToInteger(-2)
		if !ok {
			state.PushString("invalid image id provided to name")
			state.Error()
		}

		name, ok := state.ToString(-1)
		if !ok {
			state.PushString("invalid image name provided to name")
			state.Error()
		}

		go func() {
			i, err := data.IC.Image(id)
			if err != nil {
				state.PushString("invalid image provided to name")
				state.Error()
			}

			i.Name = name
			i.Mutex.Unlock()
		}()

		return 0
	})
	state.SetField(-2, "name")

	state.PushGoFunction(func(state *lua.State) int {
		question, ok := state.ToString(-1)
		if !ok {
			state.PushString("invalid question provided to prompt_file")
			state.Error()
		}

		result, err := cli.Question(question, cli.QuestionOptions{})
		if err != nil {
			state.PushString("invalid answer provided to prompt")
			state.Error()
		}

		file, err := os.Stat(result)
		if err != nil {
			state.PushString("invalid answer provided to prompt")
			state.Error()
		}

		i, id := data.IC.AddImage(file.Name())

		if file.IsDir() {
			state.PushString("directory provided to file only prompt")
			state.Error()
		}

		f, err := os.Open(result)
		if err != nil {
			state.PushString("cannot open provided file")
			state.Error()
		}
		defer f.Close()

		image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
		image, _, err := image.Decode(f)
		if err != nil {
			state.PushString("provided file is an invalid image")
			state.Error()
		}

		i.Img = &image

		state.PushInteger(id)
		return 1
	})
	state.SetField(-2, "prompt_file")

	state.PushGoFunction(func(state *lua.State) int {
		id, ok := state.ToInteger(-2)
		if !ok {
			state.PushString("invalid image id provided to out")
			state.Error()
		}

		outDir, ok := state.ToString(-1)
		if !ok {
			state.PushString("invalid outDir provided to out")
			state.Error()
		}

		_, err := os.Stat(outDir)
		if err != nil {
			os.MkdirAll(outDir, 0o666)
		}

		go func() {
			i, err := data.IC.Image(id)
			if err != nil {
				state.PushString("invalid image provided to out")
				state.Error()
			}

			f, err := os.OpenFile(path.Join(outDir, i.Name), os.O_CREATE, 0o666)
			if err != nil {
				state.PushString("cannot open provided file")
				state.Error()
			}
			defer f.Close()

			ext := filepath.Ext(i.Name)

			switch ext {
			case ".png":
				png.Encode(f, *i.Img)
			}

			i.Mutex.Unlock()
		}()

		return 0
	})
	state.SetField(-2, "out")

	state.SetGlobal(LIB_IMGSCAL)
}
