package lib

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	imageutil "github.com/ArtificialLegacy/imgscal/pkg/image_util"
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

			id := r.IC.AddItem(&chLog)

			r.IC.Schedule(id, &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					f, err := os.Open(args["path"].(string))
					if err != nil {
						r.State.PushString(i.Lg.Append("cannot open provided file", log.LEVEL_ERROR))
						r.State.Error()
					}
					defer f.Close()

					encoding := imageutil.ExtensionEncoding(path.Ext(file.Name()))
					img, err := imageutil.Decode(f, encoding)
					if err != nil {
						r.State.PushString(i.Lg.Append(fmt.Sprintf("provided file is an invalid image: %s", err), log.LEVEL_ERROR))
						r.State.Error()
					}

					i.Self = &collection.ItemImage{
						Name:     strings.TrimSuffix(file.Name(), path.Ext(file.Name())),
						Image:    img,
						Encoding: encoding,
					}
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

			r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					ext := imageutil.EncodingExtension(i.Self.Encoding)
					f, err := os.OpenFile(path.Join(args["path"].(string), i.Self.Name+ext), os.O_CREATE, 0o666)
					if err != nil {
						r.State.PushString(i.Lg.Append("cannot open provided file", log.LEVEL_ERROR))
						r.State.Error()
					}
					defer f.Close()

					lg.Append(fmt.Sprintf("encoding using %d", i.Self.Encoding), log.LEVEL_INFO)
					err = imageutil.Encode(f, i.Self.Image, i.Self.Encoding)
					if err != nil {
						r.State.PushString(i.Lg.Append(fmt.Sprintf("cannot write image to file: %s", err), log.LEVEL_ERROR))
						r.State.Error()
					}
				},
			})
			return 0
		})

	/// @func dir()
	/// @arg path
	/// @returns array containing all file paths in directory.
	lib.CreateFunction("dir",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
		},
		func(d lua.TaskData, args map[string]any) int {
			f, err := os.Stat(args["path"].(string))
			if err != nil {
				r.State.PushString(lg.Append("invalid dir path provided to io.dir", log.LEVEL_ERROR))
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
				lg.Append(fmt.Sprintf("found file %s with dir", file.Name()), log.LEVEL_INFO)

				pth := path.Join(args["path"].(string), file.Name())
				r.State.PushInteger(i)
				r.State.PushString(pth)
				r.State.SetTable(-3)
				i++
			}

			return 1
		})

	/// @func dir_img()
	/// @arg path - the directory path to scan for images.
	/// @returns array containing strings of each valid image in the directory.
	lib.CreateFunction("dir_img",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
		},
		func(d lua.TaskData, args map[string]any) int {
			parseDir("io.dir_img", args["path"].(string), []string{".png", ".jpg", ".gif"}, lib)
			return 1
		})

	// @func dir_txt()
	/// @arg path - the directory path to scan for txt.
	/// @returns array containing strings of each valid txt in the directory.
	lib.CreateFunction("dir_txt",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
		},
		func(d lua.TaskData, args map[string]any) int {
			parseDir("io.dir_txt", args["path"].(string), []string{".txt"}, lib)
			return 1
		})

	// @func dir_json()
	/// @arg path - the directory path to scan for json.
	/// @returns array containing strings of each valid json in the directory.
	lib.CreateFunction("dir_json",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
		},
		func(d lua.TaskData, args map[string]any) int {
			parseDir("io.dir_json", args["path"].(string), []string{".json"}, lib)
			return 1
		})

	/// @func mkdir()
	/// @arg path
	/// @arg all - if to create all directories going to the given path
	lib.CreateFunction("mkdir",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
			{Type: lua.BOOL, Name: "all", Optional: true},
		},
		func(d lua.TaskData, args map[string]any) int {
			if args["all"].(bool) {
				os.MkdirAll(args["path"].(string), 0o666)
			} else {
				os.Mkdir(args["path"].(string), 0o666)
			}
			return 0
		})
}

func parseDir(fn string, pathstr string, filter []string, lib *lua.Lib) {
	f, err := os.Stat(pathstr)
	if err != nil {
		lib.State.PushString(lib.Lg.Append(fmt.Sprintf("invalid dir path provided to %s", fn), log.LEVEL_ERROR))
		lib.State.Error()
	}
	if !f.IsDir() {
		lib.State.PushString(lib.Lg.Append("dir provided is not a directory", log.LEVEL_ERROR))
		lib.State.Error()
	}

	files, err := os.ReadDir(pathstr)
	if err != nil {
		lib.State.PushString(lib.Lg.Append("failed to open dir", log.LEVEL_ERROR))
		lib.State.Error()
	}

	lib.State.NewTable()

	i := 1
	for _, file := range files {
		ext := filepath.Ext(file.Name())
		if !slices.Contains(filter, ext) {
			continue
		}

		lib.Lg.Append(fmt.Sprintf("found file %s with %s", file.Name(), fn), log.LEVEL_INFO)

		pth := path.Join(pathstr, file.Name())
		lib.State.PushInteger(i)
		lib.State.PushString(pth)
		lib.State.SetTable(-3)
		i++
	}
}
