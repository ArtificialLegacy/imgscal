package lib

import (
	"bytes"
	"fmt"
	"image"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	goico "github.com/ArtificialLegacy/go-ico"
	"github.com/ArtificialLegacy/imgscal/pkg/assets"
	"github.com/ArtificialLegacy/imgscal/pkg/byteseeker"
	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	imageutil "github.com/ArtificialLegacy/imgscal/pkg/image_util"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	"github.com/crazy3lf/colorconv"
	golua "github.com/yuin/gopher-lua"
)

const LIB_IO = "io"

/// @lib IO
/// @import io
/// @desc
/// Library for handling io operations with the file system.

func RegisterIO(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_IO, r, r.State, lg)

	/// @func decode(path, model?) -> int<collection.IMAGE>
	/// @arg path {string} - The path to grab the image from.
	/// @arg? model {int<image.ColorModel>} - Used only to specify default when there is an unsupported color model.
	/// @returns {int<collection.IMAGE>}
	lib.CreateFunction(tab, "decode",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
			{Type: lua.INT, Name: "model", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			file, err := os.Stat(args["path"].(string))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("invalid image path provided to io.decode: %s", args["path"]), log.LEVEL_ERROR)), 0)
			}
			if file.IsDir() {
				state.Error(golua.LString(lg.Append("cannot load a directory as an image", log.LEVEL_ERROR)), 0)
			}

			name := strings.TrimSuffix(file.Name(), path.Ext(file.Name()))
			chLog := log.NewLogger(fmt.Sprintf("image_%s", name), lg)
			lg.Append(fmt.Sprintf("child log created: image_%s", name), log.LEVEL_INFO)

			id := r.IC.AddItem(&chLog)

			r.IC.Schedule(state, id, &collection.Task[collection.ItemImage]{
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
						Name:     name,
						Image:    img,
						Encoding: encoding,
						Model:    model,
					}
				},
			})

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func decode_string(name, encoding, data, model?) -> int<collection.IMAGE>
	/// @arg name {string}
	/// @arg encoding {int<image.Encoding>}
	/// @arg data {string}
	/// @arg? model {int<image.ColorModel>} - Used only to specify default when there is an unsupported color model.
	/// @returns {int<collection.IMAGE>}
	lib.CreateFunction(tab, "decode_string",
		[]lua.Arg{
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "encoding"},
			{Type: lua.STRING, Name: "data"},
			{Type: lua.INT, Name: "model", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			name := args["name"].(string)

			chLog := log.NewLogger(fmt.Sprintf("image_%s", name), lg)
			lg.Append(fmt.Sprintf("child log created: image_%s", name), log.LEVEL_INFO)

			id := r.IC.AddItem(&chLog)

			r.IC.Schedule(state, id, &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					encoding := lua.ParseEnum(args["encoding"].(int), imageutil.EncodingList, lib)

					img, err := imageutil.Decode(strings.NewReader(args["data"].(string)), encoding)
					if err != nil {
						state.Error(golua.LString(i.Lg.Append(fmt.Sprintf("provided data is an invalid image: %s", err), log.LEVEL_ERROR)), 0)
					}

					model := lua.ParseEnum(args["model"].(int), imageutil.ModelList, lib)
					img, model = imageutil.Limit(img, model)

					i.Self = &collection.ItemImage{
						Name:     name,
						Image:    img,
						Encoding: encoding,
						Model:    model,
					}
				},
			})

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func decode_config(path) -> int, int, bool
	/// @arg path {string} - The path to grab the image from.
	/// @returns {int} - The width of the image.
	/// @returns {int} - The height of the image.
	/// @returns {bool} - If the image can be decoded.
	lib.CreateFunction(tab, "decode_config",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			f, err := os.Open(args["path"].(string))
			if err != nil {
				state.Error(golua.LString(lg.Append("cannot open provided file", log.LEVEL_ERROR)), 0)
			}
			defer f.Close()

			encoding := imageutil.ExtensionEncoding(path.Ext(f.Name()))
			width, height, err := imageutil.DecodeConfig(f, encoding)

			state.Push(golua.LNumber(width))
			state.Push(golua.LNumber(height))
			if err != nil {
				state.Push(golua.LFalse)
				lg.Append(err.Error(), log.LEVEL_WARN)
			} else {
				state.Push(golua.LTrue)
			}
			return 0
		})

	/// @func decode_config_string(encoding, data) -> int, int, bool
	/// @arg encoding {int<image.Encoding>}
	/// @arg data {string}
	/// @returns {int} - The width of the image.
	/// @returns {int} - The height of the image.
	/// @returns {bool} - If the image can be decoded.
	lib.CreateFunction(tab, "decode_config_string",
		[]lua.Arg{
			{Type: lua.INT, Name: "encoding"},
			{Type: lua.STRING, Name: "data"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			encoding := lua.ParseEnum(args["encoding"].(int), imageutil.EncodingList, lib)
			width, height, err := imageutil.DecodeConfig(strings.NewReader(args["data"].(string)), encoding)

			state.Push(golua.LNumber(width))
			state.Push(golua.LNumber(height))
			if err != nil {
				state.Push(golua.LFalse)
				lg.Append(err.Error(), log.LEVEL_WARN)
			} else {
				state.Push(golua.LTrue)
			}
			return 0
		})

	/// @func decode_into(path, id, model?)
	/// @arg path {string} - The path to grab the image from.
	/// @arg id {int<collection.INT>} - Image ID to overwrite with decoded image.
	/// @arg? model {int<image.ColorModel>} - Used only to specify default when there is an unsupported color model.
	lib.CreateFunction(tab, "decode_into",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "model", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			file, err := os.Stat(args["path"].(string))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("invalid image path provided to io.decode_into: %s", args["path"]), log.LEVEL_ERROR)), 0)
			}
			if file.IsDir() {
				state.Error(golua.LString(lg.Append("cannot load a directory as an image", log.LEVEL_ERROR)), 0)
			}

			name := strings.TrimSuffix(file.Name(), path.Ext(file.Name()))
			id := args["id"].(int)

			r.IC.Schedule(state, id, &collection.Task[collection.ItemImage]{
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
						Name:     name,
						Image:    img,
						Encoding: encoding,
						Model:    model,
					}
				},
			})

			return 0
		})

	/// @func decode_into_string(name, encoding, data, id, model?)
	/// @arg name {string}
	/// @arg encoding {int<image.Encoding>}
	/// @arg data {string}
	/// @arg id {int<collection.INT>} - Image ID to overwrite with decoded image.
	/// @arg? model {int<image.ColorModel>} - Used only to specify default when there is an unsupported color model.
	lib.CreateFunction(tab, "decode_into_string",
		[]lua.Arg{
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "encoding"},
			{Type: lua.STRING, Name: "data"},
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "model", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			name := args["name"].(string)
			id := args["id"].(int)

			r.IC.Schedule(state, id, &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					encoding := lua.ParseEnum(args["encoding"].(int), imageutil.EncodingList, lib)
					img, err := imageutil.Decode(strings.NewReader(args["data"].(string)), encoding)
					if err != nil {
						state.Error(golua.LString(i.Lg.Append(fmt.Sprintf("provided data is an invalid image: %s", err), log.LEVEL_ERROR)), 0)
					}

					model := lua.ParseEnum(args["model"].(int), imageutil.ModelList, lib)
					img, model = imageutil.Limit(img, model)

					i.Self = &collection.ItemImage{
						Name:     name,
						Image:    img,
						Encoding: encoding,
						Model:    model,
					}
				},
			})

			return 0
		})

	/// @func decode_cached(path, model?) -> int<collection.CRATE_CACHEDIMAGE>
	/// @arg path {string} - The path to grab the image from.
	/// @arg? model {int<image.ColorModel>} - Used only to specify default when there is an unsupported color model.
	/// @returns {int<collection.CRATE_CACHEDIMAGE>}
	lib.CreateFunction(tab, "decode_cached",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
			{Type: lua.INT, Name: "model", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			pth := args["path"].(string)
			file, err := os.Stat(pth)
			if err != nil {
				lua.Error(state, lg.Appendf("invalid image path provided to io.decode_cached: %s", log.LEVEL_ERROR, pth))
			}
			if file.IsDir() {
				lua.Error(state, lg.Append("cannot load a directory as an image", log.LEVEL_ERROR))
			}

			f, err := os.Open(pth)
			if err != nil {
				lua.Error(state, lg.Append("cannot open provided file", log.LEVEL_ERROR))
			}
			defer f.Close()

			encoding := imageutil.ExtensionEncoding(path.Ext(file.Name()))
			img, err := imageutil.Decode(f, encoding)
			if err != nil {
				lua.Error(state, lg.Appendf("provided file is an invalid image: %s", log.LEVEL_ERROR, err))
			}

			model := lua.ParseEnum(args["model"].(int), imageutil.ModelList, lib)
			img, model = imageutil.Limit(img, model)

			id := r.CR_CIM.Add(&collection.CachedImageItem{
				Model: model,
				Image: img,
			})

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func decode_cached_string(encoding, data, model?) -> int<collection.CRATE_CACHEDIMAGE>
	/// @arg encoding {int<image.Encoding>}
	/// @arg data {string}
	/// @arg? model {int<image.ColorModel>} - Used only to specify default when there is an unsupported color model.
	/// @returns {int<collection.CRATE_CACHEDIMAGE>}
	lib.CreateFunction(tab, "decode_cached_string",
		[]lua.Arg{
			{Type: lua.INT, Name: "encoding"},
			{Type: lua.STRING, Name: "data"},
			{Type: lua.INT, Name: "model", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			encoding := lua.ParseEnum(args["encoding"].(int), imageutil.EncodingList, lib)
			img, err := imageutil.Decode(strings.NewReader(args["data"].(string)), encoding)
			if err != nil {
				lua.Error(state, lg.Appendf("provided data is an invalid image: %s", log.LEVEL_ERROR, err))
			}

			model := lua.ParseEnum(args["model"].(int), imageutil.ModelList, lib)
			img, model = imageutil.Limit(img, model)

			id := r.CR_CIM.Add(&collection.CachedImageItem{
				Model: model,
				Image: img,
			})

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func decode_png_data(path, model?) -> int<collection.IMAGE>, []struct<image.PNGDataChunk>
	/// @arg path {string} - The path to grab the image from.
	/// @arg? model {int<image.ColorModel>} - Used only to specify default when there is an unsupported color model.
	/// @returns {int<collection.IMAGE>}
	/// @returns {[]struct<image.PNGDataChunk>}
	lib.CreateFunction(tab, "decode_png_data",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
			{Type: lua.INT, Name: "model", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			f, err := os.Open(args["path"].(string))
			if err != nil {
				state.Error(golua.LString(lg.Append("cannot open provided file", log.LEVEL_ERROR)), 0)
			}
			defer f.Close()

			name := strings.TrimSuffix(f.Name(), path.Ext(f.Name()))
			chLog := log.NewLogger(fmt.Sprintf("image_%s", name), lg)
			lg.Append(fmt.Sprintf("child log created: image_%s", name), log.LEVEL_INFO)

			id := r.IC.AddItem(&chLog)

			img, chunks, err := imageutil.PNGDataChunkDecode(f)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("provided file is an invalid image: %s", err), log.LEVEL_ERROR)), 0)
			}

			r.IC.Schedule(state, id, &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
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

			ct := state.NewTable()
			for _, chunk := range chunks {
				t := imageutil.DataChunkToTable(chunk, state)
				ct.Append(t)
			}

			state.Push(golua.LNumber(id))
			state.Push(ct)
			return 2
		})

	/// @func decode_png_data_string(name, data, model?) -> int<collection.IMAGE>, []struct<image.PNGDataChunk>
	/// @arg name {string}
	/// @arg data {string}
	/// @arg? model {int<image.ColorModel>} - Used only to specify default when there is an unsupported color model.
	/// @returns {int<collection.IMAGE>}
	/// @returns {[]struct<image.PNGDataChunk>}
	lib.CreateFunction(tab, "decode_png_data_string",
		[]lua.Arg{
			{Type: lua.STRING, Name: "name"},
			{Type: lua.STRING, Name: "data"},
			{Type: lua.INT, Name: "model", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			name := args["name"].(string)

			chLog := log.NewLogger(fmt.Sprintf("image_%s", name), lg)
			lg.Append(fmt.Sprintf("child log created: image_%s", name), log.LEVEL_INFO)

			id := r.IC.AddItem(&chLog)

			img, chunks, err := imageutil.PNGDataChunkDecode(strings.NewReader(args["data"].(string)))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("provided data is an invalid image: %s", err), log.LEVEL_ERROR)), 0)
			}

			r.IC.Schedule(state, id, &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
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

			ct := state.NewTable()
			for _, chunk := range chunks {
				t := imageutil.DataChunkToTable(chunk, state)
				ct.Append(t)
			}

			state.Push(golua.LNumber(id))
			state.Push(ct)
			return 2
		})

	/// @func decode_favicon(path, encoding, model?) -> []int<collection.IMAGE>
	/// @arg path {string} - The path to grab the image from.
	/// @arg encoding {int<image.Encoding>} - The encoding to use for the extracted images.
	/// @arg? model {int<image.ColorModel>} - Used only to specify default when there is an unsupported color model.
	/// @returns {[]int<collection.IMAGE>}
	/// @desc
	/// Decodes an ICO type favicon file into an array of images.
	lib.CreateFunction(tab, "decode_favicon",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
			{Type: lua.INT, Name: "encoding"},
			{Type: lua.INT, Name: "model", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			file, err := os.Stat(args["path"].(string))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("invalid image path provided to io.decode_favicon: %s", args["path"]), log.LEVEL_ERROR)), 0)
			}
			if file.IsDir() {
				state.Error(golua.LString(lg.Append("cannot load a directory as an image", log.LEVEL_ERROR)), 0)
			}

			f, err := os.Open(args["path"].(string))
			if err != nil {
				state.Error(golua.LString(lg.Append("cannot open provided file", log.LEVEL_ERROR)), 0)
			}
			defer f.Close()

			cfg, imgs, err := goico.Decode(f)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("provided file is an invalid favicon: %s", err), log.LEVEL_ERROR)), 0)
			}
			if cfg.Type != goico.TYPE_ICO {
				state.Error(golua.LString(lg.Append("provided file is not an ICO type favicon", log.LEVEL_ERROR)), 0)
			}

			ids := make([]int, len(imgs))
			model := lua.ParseEnum(args["model"].(int), imageutil.ModelList, lib)

			for i, img := range imgs {
				name := fmt.Sprintf("%s_%dx%d", strings.TrimSuffix(file.Name(), path.Ext(file.Name())), img.Bounds().Dx(), img.Bounds().Dy())
				chLog := log.NewLogger("image_"+name, lg)
				lg.Append(fmt.Sprintf("child log created: %s", "image_"+name), log.LEVEL_INFO)

				id := r.IC.AddItem(&chLog)
				ids[i] = id

				r.IC.Schedule(state, id, &collection.Task[collection.ItemImage]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemImage]) {
						img, model = imageutil.Limit(img, model)

						i.Self = &collection.ItemImage{
							Name:     name,
							Image:    img,
							Encoding: lua.ParseEnum(args["encoding"].(int), imageutil.EncodingList, lib),
							Model:    model,
						}
					},
				})
			}

			t := state.NewTable()
			for i, id := range ids {
				t.RawSetInt(i+1, golua.LNumber(id))
			}
			state.Push(t)
			return 1
		})

	/// @func decode_favicon_string(name, data, encoding, model?) -> []int<collection.IMAGE>
	/// @arg name {string}
	/// @arg data{string}
	/// @arg encoding {int<image.Encoding>} - The encoding to use for the extracted images.
	/// @arg? model {int<image.ColorModel>} - Used only to specify default when there is an unsupported color model.
	/// @returns {[]int<collection.IMAGE>}
	/// @desc
	/// Decodes an ICO type favicon into an array of images.
	lib.CreateFunction(tab, "decode_favicon_string",
		[]lua.Arg{
			{Type: lua.STRING, Name: "name"},
			{Type: lua.STRING, Name: "data"},
			{Type: lua.INT, Name: "encoding"},
			{Type: lua.INT, Name: "model", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			cfg, imgs, err := goico.Decode(strings.NewReader(args["data"].(string)))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("provided data is an invalid favicon: %s", err), log.LEVEL_ERROR)), 0)
			}
			if cfg.Type != goico.TYPE_ICO {
				state.Error(golua.LString(lg.Append("provided data is not an ICO type favicon", log.LEVEL_ERROR)), 0)
			}

			ids := make([]int, len(imgs))
			model := lua.ParseEnum(args["model"].(int), imageutil.ModelList, lib)

			for i, img := range imgs {
				name := fmt.Sprintf("%s_%dx%d", args["name"].(string), img.Bounds().Dx(), img.Bounds().Dy())
				chLog := log.NewLogger("image_"+name, lg)
				lg.Append(fmt.Sprintf("child log created: %s", "image_"+name), log.LEVEL_INFO)

				id := r.IC.AddItem(&chLog)
				ids[i] = id

				r.IC.Schedule(state, id, &collection.Task[collection.ItemImage]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemImage]) {
						img, model = imageutil.Limit(img, model)

						i.Self = &collection.ItemImage{
							Name:     name,
							Image:    img,
							Encoding: lua.ParseEnum(args["encoding"].(int), imageutil.EncodingList, lib),
							Model:    model,
						}
					},
				})
			}

			t := state.NewTable()
			for i, id := range ids {
				t.RawSetInt(i+1, golua.LNumber(id))
			}
			state.Push(t)
			return 1
		})

	/// @func decode_favicon_cursor(path, encoding, model?) -> []int<collection.IMAGE>, []int
	/// @arg path {string} - The path to grab the image from.
	/// @arg encoding {int<image.Encoding>} - The encoding to use for the extracted images.
	/// @arg? model {int<image.ColorModel>} - Used only to specify default when there is an unsupported color model.
	/// @returns {[]int<collection.IMAGE>}
	/// @returns {[]int} - Pairs of ints representing the hotspot of each cursor. e.g. [x1, y1, x2, y2, ...]
	/// @desc
	/// Decodes a CUR type favicon file into an array of images and hotspots.
	lib.CreateFunction(tab, "decode_favicon_cursor",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
			{Type: lua.INT, Name: "encoding"},
			{Type: lua.INT, Name: "model", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			file, err := os.Stat(args["path"].(string))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("invalid image path provided to io.decode_favicon_cursor: %s", args["path"]), log.LEVEL_ERROR)), 0)
			}
			if file.IsDir() {
				state.Error(golua.LString(lg.Append("cannot load a directory as an image", log.LEVEL_ERROR)), 0)
			}

			f, err := os.Open(args["path"].(string))
			if err != nil {
				state.Error(golua.LString(lg.Append("cannot open provided file", log.LEVEL_ERROR)), 0)
			}
			defer f.Close()

			cfg, imgs, err := goico.Decode(f)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("provided file is an invalid favicon: %s", err), log.LEVEL_ERROR)), 0)
			}
			if cfg.Type != goico.TYPE_CUR {
				state.Error(golua.LString(lg.Append("provided file is not a CUR type favicon", log.LEVEL_ERROR)), 0)
			}

			ids := make([]int, len(imgs))
			model := lua.ParseEnum(args["model"].(int), imageutil.ModelList, lib)

			for i, img := range imgs {
				name := fmt.Sprintf("%s_%dx%d", strings.TrimSuffix(file.Name(), path.Ext(file.Name())), img.Bounds().Dx(), img.Bounds().Dy())
				chLog := log.NewLogger(name, lg)
				lg.Append(fmt.Sprintf("child log created: %s", name), log.LEVEL_INFO)

				id := r.IC.AddItem(&chLog)
				ids[i] = id

				r.IC.Schedule(state, id, &collection.Task[collection.ItemImage]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemImage]) {
						img, model = imageutil.Limit(img, model)

						i.Self = &collection.ItemImage{
							Name:     name,
							Image:    img,
							Encoding: lua.ParseEnum(args["encoding"].(int), imageutil.EncodingList, lib),
							Model:    model,
						}
					},
				})
			}

			t := state.NewTable()
			for _, id := range ids {
				t.Append(golua.LNumber(id))
			}

			ht := state.NewTable()
			for _, e := range cfg.Entries {
				ht.Append(golua.LNumber(e.Data1))
				ht.Append(golua.LNumber(e.Data2))
			}

			state.Push(t)
			state.Push(ht)
			return 2
		})

	/// @func decode_favicon_cursor_string(name, data, encoding, model?) -> []int<collection.IMAGE>, []int
	/// @arg name {string}
	/// @arg data {string}
	/// @arg encoding {int<image.Encoding>} - The encoding to use for the extracted images.
	/// @arg? model {int<image.ColorModel>} - Used only to specify default when there is an unsupported color model.
	/// @returns {[]int<collection.IMAGE>}
	/// @returns {[]int} - Pairs of ints representing the hotspot of each cursor. e.g. [x1, y1, x2, y2, ...]
	/// @desc
	/// Decodes a CUR type favicon into an array of images and hotspots.
	lib.CreateFunction(tab, "decode_favicon_cursor_string",
		[]lua.Arg{
			{Type: lua.STRING, Name: "name"},
			{Type: lua.STRING, Name: "data"},
			{Type: lua.INT, Name: "encoding"},
			{Type: lua.INT, Name: "model", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			cfg, imgs, err := goico.Decode(strings.NewReader(args["data"].(string)))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("provided data is an invalid favicon: %s", err), log.LEVEL_ERROR)), 0)
			}
			if cfg.Type != goico.TYPE_CUR {
				state.Error(golua.LString(lg.Append("provided data is not a CUR type favicon", log.LEVEL_ERROR)), 0)
			}

			ids := make([]int, len(imgs))
			model := lua.ParseEnum(args["model"].(int), imageutil.ModelList, lib)

			for i, img := range imgs {
				name := fmt.Sprintf("%s_%dx%d", args["name"].(string), img.Bounds().Dx(), img.Bounds().Dy())
				chLog := log.NewLogger(name, lg)
				lg.Append(fmt.Sprintf("child log created: %s", name), log.LEVEL_INFO)

				id := r.IC.AddItem(&chLog)
				ids[i] = id

				r.IC.Schedule(state, id, &collection.Task[collection.ItemImage]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemImage]) {
						img, model = imageutil.Limit(img, model)

						i.Self = &collection.ItemImage{
							Name:     name,
							Image:    img,
							Encoding: lua.ParseEnum(args["encoding"].(int), imageutil.EncodingList, lib),
							Model:    model,
						}
					},
				})
			}

			t := state.NewTable()
			for _, id := range ids {
				t.Append(golua.LNumber(id))
			}

			ht := state.NewTable()
			for _, e := range cfg.Entries {
				ht.Append(golua.LNumber(e.Data1))
				ht.Append(golua.LNumber(e.Data2))
			}

			state.Push(t)
			state.Push(ht)
			return 2
		})

	/// @func decode_favicon_config(path) -> struct<io.FaviconConfig>
	/// @arg path {string} - The path to grab the image from.
	/// @returns {struct<io.FaviconConfig>}
	lib.CreateFunction(tab, "decode_favicon_config",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct FaviconConfig
			/// @prop type {int<io.ICOType>} - The type of favicon.
			/// @prop count {int} - The number of images in the favicon.
			/// @prop largest {int} - The index of the largest image in the favicon.
			/// @prop entries {[]struct<io.ICOEntry>} - The entries of the favicon.

			/// @struct ICOEntry
			/// @prop width {int} - The width of the image.
			/// @prop height {int} - The height of the image.
			/// @prop data1 {int} - The x hotspot of the cursor; always 0 for icons.
			/// @prop data2 {int} - The y hotspot of the cursor; always 0 for icons.

			file, err := os.Stat(args["path"].(string))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("invalid image path provided to io.decode_favicon_config: %s", args["path"]), log.LEVEL_ERROR)), 0)
			}
			if file.IsDir() {
				state.Error(golua.LString(lg.Append("cannot load a directory as an image", log.LEVEL_ERROR)), 0)
			}

			f, err := os.Open(args["path"].(string))
			if err != nil {
				state.Error(golua.LString(lg.Append("cannot open provided file", log.LEVEL_ERROR)), 0)
			}
			defer f.Close()

			cfg, err := goico.DecodeConfig(f)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("provided file is an invalid favicon: %s", err), log.LEVEL_ERROR)), 0)
			}

			t := state.NewTable()
			t.RawSetString("type", golua.LNumber(cfg.Type))
			t.RawSetString("count", golua.LNumber(cfg.Count))
			t.RawSetString("largest", golua.LNumber(cfg.Largest))

			entries := state.NewTable()
			for _, e := range cfg.Entries {
				entry := state.NewTable()
				entry.RawSetString("width", golua.LNumber(e.Width))
				entry.RawSetString("height", golua.LNumber(e.Height))
				entry.RawSetString("data1", golua.LNumber(e.Data1))
				entry.RawSetString("data2", golua.LNumber(e.Data2))

				entries.Append(entry)
			}

			t.RawSetString("entries", entries)

			state.Push(t)
			return 1
		})

	/// @func decode_favicon_config_string(data) -> struct<io.FaviconConfig>
	/// @arg data {string}
	/// @returns {struct<io.FaviconConfig>}
	lib.CreateFunction(tab, "decode_favicon_config_string",
		[]lua.Arg{
			{Type: lua.STRING, Name: "data"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			cfg, err := goico.DecodeConfig(strings.NewReader(args["data"].(string)))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("provided data is an invalid favicon: %s", err), log.LEVEL_ERROR)), 0)
			}

			t := state.NewTable()
			t.RawSetString("type", golua.LNumber(cfg.Type))
			t.RawSetString("count", golua.LNumber(cfg.Count))
			t.RawSetString("largest", golua.LNumber(cfg.Largest))

			entries := state.NewTable()
			for _, e := range cfg.Entries {
				entry := state.NewTable()
				entry.RawSetString("width", golua.LNumber(e.Width))
				entry.RawSetString("height", golua.LNumber(e.Height))
				entry.RawSetString("data1", golua.LNumber(e.Data1))
				entry.RawSetString("data2", golua.LNumber(e.Data2))

				entries.Append(entry)
			}

			t.RawSetString("entries", entries)

			state.Push(t)
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

			r.IC.Schedule(state, id, &collection.Task[collection.ItemImage]{
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

	/// @func encode(id, path)
	/// @arg id {int<collection.IMAGE>} - The image id to encode and save to file.
	/// @arg path {string} - The directory path to save the file to.
	lib.CreateFunction(tab, "encode",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "path"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			_, err := os.Stat(args["path"].(string))
			if err != nil {
				os.MkdirAll(args["path"].(string), 0o777)
			}

			r.IC.Schedule(state, args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					ext := imageutil.EncodingExtension(i.Self.Encoding)
					f, err := os.OpenFile(path.Join(args["path"].(string), i.Self.Name+ext), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0o666)
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

	/// @func encode_string(id) -> string
	/// @arg id {int<collection.IMAGE>} - The image id to encode and save to file.
	/// @returns {string}
	/// @blocking
	lib.CreateFunction(tab, "encode_string",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var data string

			<-r.IC.Schedule(state, args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					strwriter := byteseeker.NewByteSeeker(20000, 1000)

					i.Lg.Append(fmt.Sprintf("encoding using %d", i.Self.Encoding), log.LEVEL_INFO)
					err := imageutil.Encode(strwriter, i.Self.Image, i.Self.Encoding)
					if err != nil {
						state.Error(golua.LString(i.Lg.Append(fmt.Sprintf("cannot write image to string: %s", err), log.LEVEL_ERROR)), 0)
					}

					data = string(strwriter.Bytes())
				},
			})

			state.Push(golua.LString(data))
			return 1
		})

	/// @func encode_png_data(id, chunks, path)
	/// @arg id {int<collection.IMAGE>} - The image id to encode and save to file.
	/// @arg chunks {[]struct<image.PNGDataChunk>} - The PNG data chunks to encode with the image.
	/// @arg path {string} - The directory path to save the file to.
	lib.CreateFunction(tab, "encode_png_data",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			lua.ArgArray("chunks", lua.ArrayType{Type: lua.RAW_TABLE}, false),
			{Type: lua.STRING, Name: "path"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			_, err := os.Stat(args["path"].(string))
			if err != nil {
				os.MkdirAll(args["path"].(string), 0o777)
			}

			r.IC.Schedule(state, args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					f, err := os.OpenFile(path.Join(args["path"].(string), i.Self.Name+".png"), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0o666)
					if err != nil {
						state.Error(golua.LString(i.Lg.Append("cannot open provided file", log.LEVEL_ERROR)), 0)
					}
					defer f.Close()

					data := []*imageutil.PNGDataChunk{}
					chunks := args["chunks"].([]any)

					for _, chunk := range chunks {
						data = append(data, imageutil.TableToDataChunk(chunk.(*golua.LTable)))
					}

					i.Lg.Append(fmt.Sprintf("encoding using %d, with data", imageutil.ENCODING_PNG), log.LEVEL_INFO)
					err = imageutil.PNGDataChunkEncode(f, i.Self.Image, data)
					if err != nil {
						state.Error(golua.LString(i.Lg.Append(fmt.Sprintf("cannot write image to file: %s", err), log.LEVEL_ERROR)), 0)
					}
				},
			})
			return 0
		})

	/// @func encode_png_data_string(id, chunks) -> string
	/// @arg id {int<collection.IMAGE>} - The image id to encode and save to file.
	/// @arg chunks {[]struct<image.PNGDataChunk>} - The PNG data chunks to encode with the image.
	/// @returns {string}
	/// @blocking
	lib.CreateFunction(tab, "encode_png_data_string",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			lua.ArgArray("chunks", lua.ArrayType{Type: lua.RAW_TABLE}, false),
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var data string

			<-r.IC.Schedule(state, args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					strwriter := byteseeker.NewByteSeeker(20000, 1000)

					pngdata := []*imageutil.PNGDataChunk{}
					chunks := args["chunks"].([]any)

					for _, chunk := range chunks {
						pngdata = append(pngdata, imageutil.TableToDataChunk(chunk.(*golua.LTable)))
					}

					i.Lg.Append(fmt.Sprintf("encoding using %d, with data", imageutil.ENCODING_PNG), log.LEVEL_INFO)
					err := imageutil.PNGDataChunkEncode(strwriter, i.Self.Image, pngdata)
					if err != nil {
						state.Error(golua.LString(i.Lg.Append(fmt.Sprintf("cannot write image to string: %s", err), log.LEVEL_ERROR)), 0)
					}

					data = string(strwriter.Bytes())
				},
			})

			state.Push(golua.LString(data))
			return 1
		})

	/// @func encode_favicon(name, ids, path)
	/// @arg name {string} - The name of the favicon file, without the extension.
	/// @arg ids {[]int<collection.IMAGE>} - List of image ids to encode and save into a favicon.
	/// @arg path {string} - The directory path to save the file to.
	/// @blocking
	lib.CreateFunction(tab, "encode_favicon",
		[]lua.Arg{
			{Type: lua.STRING, Name: "name"},
			lua.ArgArray("ids", lua.ArrayType{Type: lua.INT}, false),
			{Type: lua.STRING, Name: "path"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			_, err := os.Stat(args["path"].(string))
			if err != nil {
				os.MkdirAll(args["path"].(string), 0o777)
			}

			ids := args["ids"].([]any)
			imgs := make([]image.Image, len(ids))
			wg := sync.WaitGroup{}
			wg.Add(len(ids))

			for ind, id := range ids {
				r.IC.Schedule(state, id.(int), &collection.Task[collection.ItemImage]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemImage]) {
						imgs[ind] = i.Self.Image
						wg.Done()
					},
				})
			}

			wg.Wait()

			name := args["name"].(string)

			f, err := os.OpenFile(path.Join(args["path"].(string), name+".ico"), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0o666)
			if err != nil {
				lua.Error(state, lg.Append("cannot open provided file", log.LEVEL_ERROR))
			}
			defer f.Close()

			lg.Append("encoding as an ICO favicon", log.LEVEL_INFO)

			cfg, err := goico.NewICOConfig(imgs)
			if err != nil {
				lua.Error(state, lg.Appendf("failed to create ICO config: %s", log.LEVEL_ERROR, err))
			}
			err = goico.Encode(f, cfg, imgs)
			if err != nil {
				lua.Error(state, lg.Appendf("failed to encode ICO favicon: %s", log.LEVEL_ERROR, err))
			}

			return 0
		})

	/// @func encode_favicon_string(ids) -> string
	/// @arg ids {[]int<collection.IMAGE>} - List of image ids to encode and save into a favicon.
	/// @returns {string}
	/// @blocking
	lib.CreateFunction(tab, "encode_favicon_string",
		[]lua.Arg{
			lua.ArgArray("ids", lua.ArrayType{Type: lua.INT}, false),
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			ids := args["ids"].([]any)
			imgs := make([]image.Image, len(ids))
			wg := sync.WaitGroup{}
			wg.Add(len(ids))

			for ind, id := range ids {
				r.IC.Schedule(state, id.(int), &collection.Task[collection.ItemImage]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemImage]) {
						imgs[ind] = i.Self.Image
						wg.Done()
					},
				})
			}

			wg.Wait()

			strwriter := byteseeker.NewByteSeeker(20000, 1000)

			lg.Append("encoding as an ICO favicon", log.LEVEL_INFO)

			cfg, err := goico.NewICOConfig(imgs)
			if err != nil {
				lua.Error(state, lg.Appendf("failed to create ICO config: %s", log.LEVEL_ERROR, err))
			}
			err = goico.Encode(strwriter, cfg, imgs)
			if err != nil {
				lua.Error(state, lg.Appendf("failed to encode ICO favicon: %s", log.LEVEL_ERROR, err))
			}

			state.Push(golua.LString(strwriter.Bytes()))
			return 1
		})

	/// @func encode_favicon_cursor(name, ids, hotspots, path)
	/// @arg name {string} - The name of the favicon file, without the extension.
	/// @arg ids {[]int<collection.IMAGE>} - List of image ids to encode and save into a favicon.
	/// @arg hotspots {[]int} - Pairs of ints representing the hotspot of each cursor.
	/// @arg path {string} - The directory path to save the file to.
	/// @blocking
	lib.CreateFunction(tab, "encode_favicon_cursor",
		[]lua.Arg{
			{Type: lua.STRING, Name: "name"},
			lua.ArgArray("ids", lua.ArrayType{Type: lua.INT}, false),
			lua.ArgArray("hotspots", lua.ArrayType{Type: lua.INT}, false),
			{Type: lua.STRING, Name: "path"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			_, err := os.Stat(args["path"].(string))
			if err != nil {
				os.MkdirAll(args["path"].(string), 0o777)
			}

			ids := args["ids"].([]any)
			imgs := make([]image.Image, len(ids))
			wg := sync.WaitGroup{}
			wg.Add(len(ids))

			for ind, id := range ids {
				r.IC.Schedule(state, id.(int), &collection.Task[collection.ItemImage]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemImage]) {
						imgs[ind] = i.Self.Image
						wg.Done()
					},
				})
			}

			wg.Wait()

			name := args["name"].(string)

			f, err := os.OpenFile(path.Join(args["path"].(string), name+".cur"), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0o666)
			if err != nil {
				lua.Error(state, lg.Append("cannot open provided file", log.LEVEL_ERROR))
			}
			defer f.Close()

			lg.Append("encoding as a CUR favicon", log.LEVEL_INFO)

			hotspotsArg := args["hotspots"].([]any)
			hotspots := make([]int, len(hotspotsArg))
			for i, v := range hotspotsArg {
				hotspots[i] = v.(int)
			}

			cfg, err := goico.NewCURConfig(imgs, hotspots)
			if err != nil {
				lua.Error(state, lg.Appendf("failed to create CUR config: %s", log.LEVEL_ERROR, err))
			}
			err = goico.Encode(f, cfg, imgs)
			if err != nil {
				lua.Error(state, lg.Appendf("failed to encode CUR favicon: %s", log.LEVEL_ERROR, err))
			}

			return 0
		})

	/// @func encode_favicon_cursor_string(ids, hotspots) -> string
	/// @arg ids {[]int<collection.IMAGE>} - List of image ids to encode and save into a favicon.
	/// @arg hotspots {[]int} - Pairs of ints representing the hotspot of each cursor.
	/// @returns {string}
	/// @blocking
	lib.CreateFunction(tab, "encode_favicon_cursor_string",
		[]lua.Arg{
			lua.ArgArray("ids", lua.ArrayType{Type: lua.INT}, false),
			lua.ArgArray("hotspots", lua.ArrayType{Type: lua.INT}, false),
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			ids := args["ids"].([]any)
			imgs := make([]image.Image, len(ids))
			wg := sync.WaitGroup{}
			wg.Add(len(ids))

			for ind, id := range ids {
				r.IC.Schedule(state, id.(int), &collection.Task[collection.ItemImage]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemImage]) {
						imgs[ind] = i.Self.Image
						wg.Done()
					},
				})
			}

			wg.Wait()

			strwriter := byteseeker.NewByteSeeker(20000, 1000)

			lg.Append("encoding as a CUR favicon", log.LEVEL_INFO)

			hotspotsArg := args["hotspots"].([]any)
			hotspots := make([]int, len(hotspotsArg))
			for i, v := range hotspotsArg {
				hotspots[i] = v.(int)
			}

			cfg, err := goico.NewCURConfig(imgs, hotspots)
			if err != nil {
				lua.Error(state, lg.Appendf("failed to create CUR config: %s", log.LEVEL_ERROR, err))
			}
			err = goico.Encode(strwriter, cfg, imgs)
			if err != nil {
				lua.Error(state, lg.Appendf("failed to encode CUR favicon: %s", log.LEVEL_ERROR, err))
			}

			state.Push(golua.LString(strwriter.Bytes()))
			return 1
		})

	/// @func load_palette(path) -> []struct<image.Color>
	/// @arg path {string} - Path to a .hex file for the palette.
	/// @returns {[]struct<image.Color>}
	/// @desc
	/// Use to load a hex color palette file; For example, from lospec.
	lib.CreateFunction(tab, "load_palette",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			pth := args["path"].(string)
			b, err := os.ReadFile(pth)
			if err != nil {
				lua.Error(state, lg.Appendf("failed to read hex file: %s with error (%s)", log.LEVEL_ERROR, pth, err))
			}

			hexValues := strings.Split(string(b), "\n")
			colors := state.NewTable()

			for i, v := range hexValues {
				r, g, b, err := colorconv.HexToRGB(v)
				if err != nil {
					lua.Error(state, lg.Appendf("failed to parse hex color: %s with error (%s)", log.LEVEL_ERROR, v, err))
				}

				colors.RawSetInt(i+1, imageutil.RGBAToColorTable(state, int(r), int(g), int(b), 255))
			}

			state.Push(colors)
			return 1
		})

	/// @func save_palette(path, colors)
	/// @arg path {string} - File to save hex data to, filename should end in .hex.
	/// @arg colors {[]struct<image.Color>}
	/// @desc
	/// Discards alpha channels.
	lib.CreateFunction(tab, "save_palette",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
			lua.ArgArray("colors", lua.ArrayType{Type: lua.RAW_TABLE}, false),
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			pth := args["path"].(string)
			colors := args["colors"].([]any)

			fs, err := os.OpenFile(pth, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0o666)
			if err != nil {
				lua.Error(state, lg.Appendf("failed to open file: %s with error (%s)", log.LEVEL_ERROR, pth, err))
			}
			defer fs.Close()

			for _, v := range colors {
				col := v.(*golua.LTable)
				r, g, b, _ := imageutil.ColorTableToRGBA(col)

				fmt.Fprintf(fs, "%02x%02x%02x\n", r, g, b)
			}

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
			files := parseDir("io.dir_img", args["path"].(string), imageutil.EncodingExts, state, lg)

			state.Push(files)
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
			files := parseDir("io.dir_txt", args["path"].(string), []string{".txt"}, state, lg)

			state.Push(files)
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
			files := parseDir("io.dir_json", args["path"].(string), []string{".json"}, state, lg)

			state.Push(files)
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
			files := parseDirDir("io.dir_dir", args["path"].(string), state, lg)

			state.Push(files)
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

			files := parseDir("dir_filter", args["path"].(string), filter, state, lg)

			state.Push(files)
			return 1
		})

	/// @func dir_recursive(path) -> []string
	/// @arg path {string} - The directory path to scan for files, within all sub-directories.
	/// @returns {[]string} - Array containing paths to each valid file in the directory and sub-directories.
	lib.CreateFunction(tab, "dir_recursive",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			files := parseDirRecursive(args["path"].(string), state, lg)

			state.Push(files)
			return 1
		})

	/// @func dir_recursive_filter(path, filter) -> []string
	/// @arg path {string} - The directory path to scan for files, within all sub-directories.
	/// @arg filter {[]string} - Array of file extensions to include.
	/// @returns {[]string} - Array containing paths to each valid file in the directory and sub-directories.
	lib.CreateFunction(tab, "dir_recursive_filter",
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

			files := parseDirRecursiveFilter(args["path"].(string), filter, state, lg)

			state.Push(files)
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

	/// @func default_input() -> string
	/// @returns {string}
	/// @desc
	/// Returns the default input directory specified in the config file.
	lib.CreateFunction(tab, "default_input",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			if !r.UseDefaultInput {
				lua.Error(state, lg.Append("cannot use default_input, it has not been enabled within the init function", log.LEVEL_ERROR))
			}
			pth := path.Join(r.Config.InputDirectory, r.Entry)
			state.Push(golua.LString(pth))
			return 1
		})

	/// @func default_output() -> string
	/// @returns {string}
	/// @desc
	/// Returns the default output directory specified in the config file.
	lib.CreateFunction(tab, "default_output",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			if !r.UseDefaultOutput {
				lua.Error(state, lg.Append("cannot use default_output, it has not been enabled within the init function", log.LEVEL_ERROR))
			}
			pth := path.Join(r.Config.OutputDirectory, r.Entry)
			state.Push(golua.LString(pth))
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

	/// @constants Embedded {int}
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

	/// @constants ICOType {int}
	/// @const ICOTYPE_ICO
	/// @const ICOTYPE_CUR
	tab.RawSetString("ICOTYPE_ICO", golua.LNumber(goico.TYPE_ICO))
	tab.RawSetString("ICOTYPE_CUR", golua.LNumber(goico.TYPE_CUR))
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

func parseDir(fn string, pathstr string, filter []string, state *golua.LState, lg *log.Logger) *golua.LTable {
	f, err := os.Stat(pathstr)
	if err != nil {
		lua.Error(state, lg.Appendf("invalid dir path provided to %s: %s", log.LEVEL_ERROR, fn, pathstr))
	}
	if !f.IsDir() {
		lua.Error(state, lg.Appendf("dir provided to %s is not a directory: %s", log.LEVEL_ERROR, fn, pathstr))
	}

	files, err := os.ReadDir(pathstr)
	if err != nil {
		lua.Error(state, lg.Appendf("failed to open dir: %s with error (%s)", log.LEVEL_ERROR, pathstr, err))
	}

	t := state.NewTable()

	i := 1
	for _, file := range files {
		ext := filepath.Ext(file.Name())
		if !slices.Contains(filter, ext) {
			continue
		}

		pth := path.Join(pathstr, file.Name())
		t.RawSetInt(i, golua.LString(pth))
		i++
	}

	return t
}

func parseDirDir(fn string, pathstr string, state *golua.LState, lg *log.Logger) *golua.LTable {
	f, err := os.Stat(pathstr)
	if err != nil {
		lua.Error(state, lg.Appendf("invalid dir path provided to %s: %s", log.LEVEL_ERROR, fn, pathstr))
	}
	if !f.IsDir() {
		lua.Error(state, lg.Appendf("dir provided to %s is not a directory: %s", log.LEVEL_ERROR, fn, pathstr))
	}

	files, err := os.ReadDir(pathstr)
	if err != nil {
		lua.Error(state, lg.Appendf("failed to open dir: %s with error (%s)", log.LEVEL_ERROR, pathstr, err))
	}

	t := state.NewTable()

	i := 1
	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		pth := path.Join(pathstr, file.Name())
		t.RawSetInt(i, golua.LString(pth))
		i++
	}

	return t
}

func parseDirRecursive(pathstr string, state *golua.LState, lg *log.Logger) *golua.LTable {
	f, err := os.Stat(pathstr)
	if err != nil {
		lua.Error(state, lg.Appendf("invalid dir path provided to io.dir_recursive: %s", log.LEVEL_ERROR, pathstr))
	}
	if !f.IsDir() {
		lua.Error(state, lg.Appendf("dir provided to io.dir_recursive is not a directory: %s", log.LEVEL_ERROR, pathstr))
	}

	t := state.NewTable()

	files, err := os.ReadDir(pathstr)
	if err != nil {
		lua.Error(state, lg.Appendf("failed to open dir: %s with error (%s)", log.LEVEL_ERROR, pathstr, err))
	}

	for i, file := range files {
		child := path.Join(pathstr, file.Name())

		if file.IsDir() {
			t.RawSetInt(i+1, parseDirRecursive(child, state, lg))
		} else {
			t.RawSetInt(i+1, golua.LString(child))
		}
	}

	return t
}

func parseDirRecursiveFilter(pathstr string, filter []string, state *golua.LState, lg *log.Logger) *golua.LTable {
	f, err := os.Stat(pathstr)
	if err != nil {
		lua.Error(state, lg.Appendf("invalid dir path provided to io.dir_recursive_filter: %s", log.LEVEL_ERROR, pathstr))
	}
	if !f.IsDir() {
		lua.Error(state, lg.Appendf("dir provided to io.dir_recursive_filter is not a directory: %s", log.LEVEL_ERROR, pathstr))
	}
	t := state.NewTable()

	files, err := os.ReadDir(pathstr)
	if err != nil {
		lua.Error(state, lg.Appendf("failed to open dir: %s with error (%s)", log.LEVEL_ERROR, pathstr, err))
	}

	i := 0
	for _, file := range files {
		child := path.Join(pathstr, file.Name())

		if file.IsDir() {
			t.RawSetInt(i+1, parseDirRecursive(child, state, lg))
		} else {
			ext := path.Ext(child)
			if slices.Contains(filter, ext) {
				continue
			}

			t.RawSetInt(i+1, golua.LString(child))
		}

		i++
	}

	return t
}
