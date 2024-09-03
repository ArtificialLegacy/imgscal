package lib

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/ArtificialLegacy/imgscal/pkg/assets"
	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	imageutil "github.com/ArtificialLegacy/imgscal/pkg/image_util"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	golua "github.com/yuin/gopher-lua"
)

const LIB_IO = "io"

/// @lib IO
/// @import io
/// @desc
/// Library for handling io operations with the file system.

func RegisterIO(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_IO, r, r.State, lg)

	/// @func load_image(path, model?) -> int<collection.IMAGE>
	/// @arg path {string} - The path to grab the image from.
	/// @arg? model {int<image.ColorModel>} - Used only to specify default when there is an unsupported color model.
	/// @returns {int<collection.IMAGE>}
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

			chLog := log.NewLogger(fmt.Sprintf("image_%s", file.Name()), lg)
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

	/// @func load_embedded(embedded, model?) -> int<collection.IMAGE>
	/// @arg embedded {int<io.Embedded>}
	/// @arg? model {int<image.ColorModel>} - Used only to specify default of unsupported color models.
	/// @returns {int<collection.IMAGE>}
	lib.CreateFunction(tab, "load_embedded",
		[]lua.Arg{
			{Type: lua.INT, Name: "embedded"},
			{Type: lua.INT, Name: "model", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var asset []byte
			var name string

			switch args["embedded"].(int) {
			case EMBEDDED_ICONCIRCLE_16x16:
				asset = assets.FAVICON_16x16_circle
				name = "embedded_favicon_16x16_circle"
			case EMBEDDED_ICONCIRCLE_32x32:
				asset = assets.FAVICON_32x32_circle
				name = "embedded_favicon_32x32_circle"
			case EMBEDDED_ICON_16x16:
				asset = assets.FAVICON_16x16
				name = "embedded_favicon_16x16"
			case EMBEDDED_ICON_32x32:
				asset = assets.FAVICON_32x32
				name = "embedded_favicon_32x32"
			case EMBEDDED_ICON_180x180:
				asset = assets.FAVICON_180x180
				name = "embedded_favicon_180x180"
			case EMBEDDED_ICON_192x192:
				asset = assets.FAVICON_192x192
				name = "embedded_favicon_192x192"
			case EMBEDDED_ICON_512x512:
				asset = assets.FAVICON_512x512
				name = "embedded_favicon_512x512"

			default:
				state.Error(golua.LString(lg.Append(fmt.Sprintf("invalid embedded asset selected: %d", args["embedded"]), log.LEVEL_ERROR)), 0)
			}

			chLog := log.NewLogger(fmt.Sprintf("image_%s", name), lg)
			lg.Append(fmt.Sprintf("child log created: image_%s", name), log.LEVEL_INFO)

			id := r.IC.AddItem(&chLog)

			r.IC.Schedule(id, &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					img, err := imageutil.Decode(bytes.NewReader(asset), imageutil.ENCODING_PNG)
					if err != nil {
						state.Error(golua.LString(i.Lg.Append(fmt.Sprintf("provided embedded is an invalid image: %s", err), log.LEVEL_ERROR)), 0)
					}

					model := lua.ParseEnum(args["model"].(int), imageutil.ModelList, lib)
					img, model = imageutil.Limit(img, model)

					i.Self = &collection.ItemImage{
						Name:     name,
						Image:    img,
						Encoding: imageutil.ENCODING_PNG,
						Model:    model,
					}
				},
			})

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func out(id, path)
	/// @arg id {int<collection.IMAGE>} - The image id to encode and save to file.
	/// @arg path {string} - The directory path to save the file to.
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

	/// @func remove(path, all?)
	/// @arg path {string}
	/// @arg? all {bool} - If to remove all directories going to the given path.
	lib.CreateFunction(tab, "remove",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
			{Type: lua.BOOL, Name: "all", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			if args["all"].(bool) {
				err := os.RemoveAll(args["path"].(string))
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to remove all directories: %s", err), log.LEVEL_ERROR)), 0)
				}
			} else {
				err := os.Remove(args["path"].(string))
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to remove file: %s", err), log.LEVEL_ERROR)), 0)
				}
			}

			return 0
		})

	/// @func exists(path) -> bool, bool
	/// @arg path {string}
	/// @returns {bool} - If the file exists.
	/// @returns {bool} - If the file is a directory.
	lib.CreateFunction(tab, "exists",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			fs, err := os.Stat(args["path"].(string))
			if err != nil {
				state.Push(golua.LFalse)
				state.Push(golua.LFalse)
			} else {
				state.Push(golua.LTrue)
				state.Push(golua.LBool(fs.IsDir()))
			}
			return 2
		})

	/// @func dir(path) -> []string
	/// @arg path {string}
	/// @returns {[]string} - Array containing all file paths in the directory.
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

	/// @func dir_img(path) -> []string
	/// @arg path {string} - The directory path to scan for images.
	/// @returns {[]string} - Array containing paths to each valid image in the directory.
	lib.CreateFunction(tab, "dir_img",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			parseDir("io.dir_img", args["path"].(string), []string{".png", ".jpg", ".gif"}, lib)
			return 1
		})

	/// @func dir_txt(path) -> []string
	/// @arg path {string} - The directory path to scan for txt.
	/// @returns {[]string} - Array containing paths to each valid txt in the directory.
	lib.CreateFunction(tab, "dir_txt",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			parseDir("io.dir_txt", args["path"].(string), []string{".txt"}, lib)
			return 1
		})

	/// @func dir_json(path) -> []string
	/// @arg path {string} - The directory path to scan for json.
	/// @returns {[]string} - Array containing paths to each valid json in the directory.
	lib.CreateFunction(tab, "dir_json",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			parseDir("io.dir_json", args["path"].(string), []string{".json"}, lib)
			return 1
		})

	/// @func dir_dir(path) -> []string
	/// @arg path {string} - The directory path to scan for directories.
	/// @returns {[]string} - Array containing paths to each valid dir in the directory.
	lib.CreateFunction(tab, "dir_dir",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			parseDirDir("io.dir_dir", args["path"].(string), lib)
			return 1
		})

	/// @func dir_filter(path, filter) -> []string
	/// @arg path {string} - The directory path to scan.
	/// @arg filter {[]string} - Array of file extensions to include.
	/// @returns {[]string} - Array containing paths to all files that match the filter.
	lib.CreateFunction(tab, "dir_filter",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
			lua.ArgArray("filter", lua.ArrayType{Type: lua.STRING}, false),
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			fv := args["filter"].([]any)
			filter := make([]string, len(fv))
			for i, v := range fv {
				filter[i] = v.(string)
			}

			parseDir("dir_filter", args["path"].(string), filter, lib)
			return 1
		})

	/// @func mkdir(path, all?)
	/// @arg path {string}
	/// @arg? all {bool} - If to create all directories going to the given path.
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

	/// @func path_join(paths...) -> string
	/// @arg paths {string...}
	/// @returns {string}
	lib.CreateFunction(tab, "path_join",
		[]lua.Arg{
			lua.ArgVariadic("paths", lua.ArrayType{Type: lua.STRING}, false),
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			pths := args["paths"].([]any)
			strs := make([]string, len(pths))

			for i, v := range pths {
				strs[i] = v.(string)
			}

			pth := path.Join(strs...)

			state.Push(golua.LString(pth))
			return 1
		})

	/// @func wd() -> string
	/// @returns {string}
	/// @desc
	/// Returns the dir of the currently running workflow.
	lib.CreateFunction(tab, "wd",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			state.Push(golua.LString(r.Dir))
			return 1
		})

	/// @func default_output() -> string
	/// @returns {string}
	/// @desc
	/// Returns the default output directory specified in the config file.
	lib.CreateFunction(tab, "default_output",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			state.Push(golua.LString(r.Config.OutputDirectory))
			return 1
		})

	/// @func base(pth) -> string
	/// @arg pth {string}
	/// @returns {string}
	/// @desc
	/// Returns the name of the file only, without the extension or trailing path.
	lib.CreateFunction(tab, "base",
		[]lua.Arg{
			{Type: lua.STRING, Name: "pth"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			pth := args["pth"].(string)
			pth = path.Base(pth)
			base := strings.TrimSuffix(pth, path.Ext(pth))

			state.Push(golua.LString(base))
			return 1
		})

	/// @func path_to(pth) -> string
	/// @arg pth {string}
	/// @returns {string}
	/// @desc
	/// Returns the path to a file.
	lib.CreateFunction(tab, "path_to",
		[]lua.Arg{
			{Type: lua.STRING, Name: "pth"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			pth := args["pth"].(string)
			dir := path.Dir(pth)

			state.Push(golua.LString(dir))
			return 1
		})

	/// @func ext(pth) -> string
	/// @arg pth {string}
	/// @returns {string}
	/// @desc
	/// Returns the ext of the file only.
	lib.CreateFunction(tab, "ext",
		[]lua.Arg{
			{Type: lua.STRING, Name: "pth"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			pth := args["pth"].(string)
			ext := path.Ext(pth)

			state.Push(golua.LString(ext))
			return 1
		})

	/// @constants Embedded Assets
	/// @const  EMBEDDED_ICONCIRCLE_16x16
	/// @const  EMBEDDED_ICONCIRCLE_32x32
	/// @const  EMBEDDED_ICON_16x16
	/// @const  EMBEDDED_ICON_32x32
	/// @const  EMBEDDED_ICON_180x180
	/// @const  EMBEDDED_ICON_192x192
	/// @const  EMBEDDED_ICON_512x512
	tab.RawSetString("EMBEDDED_ICONCIRCLE_16x16", golua.LNumber(EMBEDDED_ICONCIRCLE_16x16))
	tab.RawSetString("EMBEDDED_ICONCIRCLE_32x32", golua.LNumber(EMBEDDED_ICONCIRCLE_32x32))
	tab.RawSetString("EMBEDDED_ICON_16x16", golua.LNumber(EMBEDDED_ICON_16x16))
	tab.RawSetString("EMBEDDED_ICON_32x32", golua.LNumber(EMBEDDED_ICON_32x32))
	tab.RawSetString("EMBEDDED_ICON_180x180", golua.LNumber(EMBEDDED_ICON_180x180))
	tab.RawSetString("EMBEDDED_ICON_192x192", golua.LNumber(EMBEDDED_ICON_192x192))
	tab.RawSetString("EMBEDDED_ICON_512x512", golua.LNumber(EMBEDDED_ICON_512x512))
}

const (
	EMBEDDED_ICONCIRCLE_16x16 int = iota
	EMBEDDED_ICONCIRCLE_32x32
	EMBEDDED_ICON_16x16
	EMBEDDED_ICON_32x32
	EMBEDDED_ICON_180x180
	EMBEDDED_ICON_192x192
	EMBEDDED_ICON_512x512
)

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
