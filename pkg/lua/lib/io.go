package lib

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	imageutil "github.com/ArtificialLegacy/imgscal/pkg/image_util"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	golua "github.com/yuin/gopher-lua"
)

const LIB_IO = "io"

func RegisterIO(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_IO, r, r.State, lg)

	/// @func load_image()
	/// @arg path - the path to grab the image from
	/// @arg? model - used only to specify default of unsupported color models
	/// @returns int - the image id
	lib.CreateFunction(tab, "load_image",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
			{Type: lua.INT, Name: "model", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			file, err := os.Stat(args["path"].(string))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("invalid image path provided to io.load_image: %s", args["path"]), log.LEVEL_ERROR)), 0)
			}
			if file.IsDir() {
				state.Error(golua.LString(lg.Append("cannot load a directory as an image", log.LEVEL_ERROR)), 0)
			}

			chLog := log.NewLogger(fmt.Sprintf("image_%s", file.Name()))
			chLog.Parent(lg)
			lg.Append(fmt.Sprintf("child log created: image_%s", file.Name()), log.LEVEL_INFO)

			id := r.IC.AddItem(&chLog)

			r.IC.Schedule(id, &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					f, err := os.Open(args["path"].(string))
					if err != nil {
						state.Error(golua.LString(i.Lg.Append("cannot open provided file", log.LEVEL_ERROR)), 0)
					}
					defer f.Close()

					encoding := imageutil.ExtensionEncoding(path.Ext(file.Name()))
					img, err := imageutil.Decode(f, encoding)
					if err != nil {
						state.Error(golua.LString(i.Lg.Append(fmt.Sprintf("provided file is an invalid image: %s", err), log.LEVEL_ERROR)), 0)
					}

					model := lua.ParseEnum(args["model"].(int), imageutil.ModelList, lib)
					img, model = imageutil.Limit(img, model)

					i.Self = &collection.ItemImage{
						Name:     strings.TrimSuffix(file.Name(), path.Ext(file.Name())),
						Image:    img,
						Encoding: encoding,
						Model:    model,
					}
				},
			})

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func out()
	/// @arg image_id - the image id to encode and save to file.
	/// @arg path - the directory path to save the file to.
	lib.CreateFunction(tab, "out",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "path"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			_, err := os.Stat(args["path"].(string))
			if err != nil {
				os.MkdirAll(args["path"].(string), 0o777)
			}

			r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					ext := imageutil.EncodingExtension(i.Self.Encoding)
					f, err := os.OpenFile(path.Join(args["path"].(string), i.Self.Name+ext), os.O_CREATE|os.O_RDWR, 0o666)
					if err != nil {
						state.Error(golua.LString(i.Lg.Append("cannot open provided file", log.LEVEL_ERROR)), 0)
					}
					defer f.Close()

					i.Lg.Append(fmt.Sprintf("encoding using %d", i.Self.Encoding), log.LEVEL_INFO)
					err = imageutil.Encode(f, i.Self.Image, i.Self.Encoding)
					if err != nil {
						state.Error(golua.LString(i.Lg.Append(fmt.Sprintf("cannot write image to file: %s", err), log.LEVEL_ERROR)), 0)
					}
				},
			})
			return 0
		})

	/// @func dir()
	/// @arg path
	/// @returns array containing all file paths in directory.
	lib.CreateFunction(tab, "dir",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			f, err := os.Stat(args["path"].(string))
			if err != nil {
				state.Error(golua.LString(lg.Append("invalid dir path provided to io.dir", log.LEVEL_ERROR)), 0)
			}
			if !f.IsDir() {
				state.Error(golua.LString(lg.Append("dir provided is not a directory", log.LEVEL_ERROR)), 0)
			}

			files, err := os.ReadDir(args["path"].(string))
			if err != nil {
				state.Error(golua.LString(lg.Append("failed to open dir", log.LEVEL_ERROR)), 0)
			}

			t := state.NewTable()

			i := 1
			for _, file := range files {
				lg.Append(fmt.Sprintf("found file %s with dir", file.Name()), log.LEVEL_INFO)

				pth := path.Join(args["path"].(string), file.Name())
				state.SetTable(t, golua.LNumber(i), golua.LString(pth))
				i++
			}

			state.Push(t)
			return 1
		})

	/// @func dir_img()
	/// @arg path - the directory path to scan for images.
	/// @returns array containing strings of each valid image in the directory.
	lib.CreateFunction(tab, "dir_img",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			parseDir("io.dir_img", args["path"].(string), []string{".png", ".jpg", ".gif"}, lib)
			return 1
		})

	// @func dir_txt()
	/// @arg path - the directory path to scan for txt.
	/// @returns array containing strings of each valid txt in the directory.
	lib.CreateFunction(tab, "dir_txt",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			parseDir("io.dir_txt", args["path"].(string), []string{".txt"}, lib)
			return 1
		})

	// @func dir_json()
	/// @arg path - the directory path to scan for json.
	/// @returns array containing strings of each valid json in the directory.
	lib.CreateFunction(tab, "dir_json",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			parseDir("io.dir_json", args["path"].(string), []string{".json"}, lib)
			return 1
		})

	// @func dir_dir()
	/// @arg path - the directory path to scan for directories.
	/// @returns array containing strings of each valid dir in the directory.
	lib.CreateFunction(tab, "dir_dir",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			parseDirDir("io.dir_dir", args["path"].(string), lib)
			return 1
		})

	/// @func dir_filter()
	/// @arg path - the directory path to scan
	/// @arg filter - array of file paths to include
	/// @returns array containing all files that match the filter
	lib.CreateFunction(tab, "dir_filter",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
			lua.ArgArray("filter", lua.ArrayType{Type: lua.STRING}, false),
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			filter := []string{}
			for _, v := range args["filter"].(map[string]any) {
				filter = append(filter, v.(string))
			}

			parseDir("dir_filter", args["path"].(string), filter, lib)
			return 1
		})

	/// @func mkdir()
	/// @arg path
	/// @arg? all - if to create all directories going to the given path
	lib.CreateFunction(tab, "mkdir",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
			{Type: lua.BOOL, Name: "all", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			if args["all"].(bool) {
				os.MkdirAll(args["path"].(string), 0o777)
			} else {
				os.Mkdir(args["path"].(string), 0o777)
			}
			return 0
		})

	/// @func path_join()
	/// @arg []string
	/// @returns path
	lib.CreateFunction(tab, "path_join",
		[]lua.Arg{
			lua.ArgArray("paths", lua.ArrayType{Type: lua.STRING}, false),
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			strs := []string{}
			pths := args["paths"].(map[string]any)

			for i := range len(pths) {
				strs = append(strs, pths[strconv.Itoa(i+1)].(string))
			}

			pth := path.Join(strs...)

			state.Push(golua.LString(pth))
			return 1
		})

	/// @func wd()
	/// @returns string
	/// @desc
	/// returns the dir of the currently running workflow
	lib.CreateFunction(tab, "wd",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			state.Push(golua.LString(r.Dir))
			return 1
		})
}

func parseDir(fn string, pathstr string, filter []string, lib *lua.Lib) {
	f, err := os.Stat(pathstr)
	if err != nil {
		lib.State.Error(golua.LString(lib.Lg.Append(fmt.Sprintf("invalid dir path provided to %s", fn), log.LEVEL_ERROR)), 0)
	}
	if !f.IsDir() {
		lib.State.Error(golua.LString(lib.Lg.Append("dir provided is not a directory", log.LEVEL_ERROR)), 0)
	}

	files, err := os.ReadDir(pathstr)
	if err != nil {
		lib.State.Error(golua.LString(lib.Lg.Append("failed to open dir", log.LEVEL_ERROR)), 0)
	}

	t := lib.State.NewTable()

	i := 1
	for _, file := range files {
		ext := filepath.Ext(file.Name())
		if !slices.Contains(filter, ext) {
			continue
		}

		lib.Lg.Append(fmt.Sprintf("found file %s with %s", file.Name(), fn), log.LEVEL_INFO)

		pth := path.Join(pathstr, file.Name())
		lib.State.SetTable(t, golua.LNumber(i), golua.LString(pth))
		i++
	}

	lib.State.Push(t)
}

func parseDirDir(fn string, pathstr string, lib *lua.Lib) {
	f, err := os.Stat(pathstr)
	if err != nil {
		lib.State.Error(golua.LString(lib.Lg.Append(fmt.Sprintf("invalid dir path provided to %s", fn), log.LEVEL_ERROR)), 0)
	}
	if !f.IsDir() {
		lib.State.Error(golua.LString(lib.Lg.Append("dir provided is not a directory", log.LEVEL_ERROR)), 0)
	}

	files, err := os.ReadDir(pathstr)
	if err != nil {
		lib.State.Error(golua.LString(lib.Lg.Append("failed to open dir", log.LEVEL_ERROR)), 0)
	}

	t := lib.State.NewTable()

	i := 1
	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		lib.Lg.Append(fmt.Sprintf("found dir %s with %s", file.Name(), fn), log.LEVEL_INFO)

		pth := path.Join(pathstr, file.Name())
		lib.State.SetTable(t, golua.LNumber(i), golua.LString(pth))
		i++
	}

	lib.State.Push(t)
}
