package lib

import (
	"encoding/base64"
	"fmt"
	"image"
	"os"

	"github.com/ArtificialLegacy/imgscal/pkg/byteseeker"
	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	imageutil "github.com/ArtificialLegacy/imgscal/pkg/image_util"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	golua "github.com/yuin/gopher-lua"
)

const LIB_BASE64 = "base64"

/// @lib Base64
/// @import base64
/// @desc
/// Utility library for encoding and decoding base64 data.

func RegisterBase64(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_BASE64, r, r.State, lg)

	/// @func encode(data, url?, raw?) -> string
	/// @arg data {string}
	/// @arg? url {bool} - If true, use URL encoding.
	/// @arg? raw {bool} - If true, use raw encoding.
	/// @returns {string} - The base64 encoded data.
	lib.CreateFunction(tab, "encode",
		[]lua.Arg{
			{Type: lua.STRING, Name: "data"},
			{Type: lua.BOOL, Name: "url", Optional: true},
			{Type: lua.BOOL, Name: "raw", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			data := args["data"].(string)
			raw := args["raw"].(bool)

			var out string
			if !args["url"].(bool) {
				if raw {
					out = base64.RawStdEncoding.EncodeToString([]byte(data))
				} else {
					out = base64.StdEncoding.EncodeToString([]byte(data))
				}
			} else {
				if raw {
					out = base64.RawURLEncoding.EncodeToString([]byte(data))
				} else {
					out = base64.URLEncoding.EncodeToString([]byte(data))
				}
			}

			state.Push(golua.LString(out))
			return 1
		})

	/// @func decode(data, url?, raw?) -> string
	/// @arg data {string}
	// @arg? url {bool} - If true, use URL encoding.
	/// @arg? raw {bool} - If true, use raw encoding.
	/// @returns {string} - The base64 decoded data.
	lib.CreateFunction(tab, "decode",
		[]lua.Arg{
			{Type: lua.STRING, Name: "data"},
			{Type: lua.BOOL, Name: "url", Optional: true},
			{Type: lua.BOOL, Name: "raw", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			data := args["data"].(string)
			raw := args["raw"].(bool)

			var out []byte
			var err error
			if !args["url"].(bool) {
				if raw {
					out, err = base64.RawStdEncoding.DecodeString(data)
				} else {
					out, err = base64.StdEncoding.DecodeString(data)
				}
			} else {
				if raw {
					out, err = base64.RawURLEncoding.DecodeString(data)
				} else {
					out, err = base64.URLEncoding.DecodeString(data)
				}
			}

			if err != nil {
				lua.Error(state, lg.Appendf("failed to decode base64 data: %s", log.LEVEL_ERROR, err))
			}

			state.Push(golua.LString(out))
			return 1
		})

	/// @func encode_image(id, raw?) -> string
	/// @arg id {int<collection.IMAGE>} - The image to encode.
	/// @arg? raw {bool} - If true, use raw encoding.
	/// @returns {string} - The base64 encoded data.
	/// @blocking
	lib.CreateFunction(tab, "encode_image",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.BOOL, Name: "raw", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			id := args["id"].(int)
			raw := args["raw"].(bool)

			var img image.Image
			var encoding imageutil.ImageEncoding

			<-r.IC.Schedule(id, &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					img = i.Self.Image
					encoding = i.Self.Encoding
				},
			})

			wb := byteseeker.NewByteSeeker(20000, 1000)
			err := imageutil.Encode(wb, img, encoding)
			if err != nil {
				lua.Error(state, lg.Appendf("failed to encode image: %s", log.LEVEL_ERROR, err))
			}

			var out string
			if raw {
				out = base64.RawStdEncoding.EncodeToString(wb.Bytes())
			} else {
				out = base64.StdEncoding.EncodeToString(wb.Bytes())
			}

			state.Push(golua.LString(out))
			return 1
		})

	/// @func encode_image_to_file(id, path, raw?)
	/// @arg id {int<collection.IMAGE>} - The image to encode.
	/// @arg path {string} - The path to save the image.
	/// @arg? raw {bool} - If true, use raw encoding.
	lib.CreateFunction(tab, "encode_image_to_file",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "path"},
			{Type: lua.BOOL, Name: "raw", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			id := args["id"].(int)
			raw := args["raw"].(bool)
			pth := args["path"].(string)

			r.IC.Schedule(id, &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					f, err := os.OpenFile(pth, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
					if err != nil {
						lua.Error(state, i.Lg.Appendf("failed to open file: %s", log.LEVEL_ERROR, err))
					}
					defer f.Close()

					wb := byteseeker.NewByteSeeker(20000, 1000)
					err = imageutil.Encode(wb, i.Self.Image, i.Self.Encoding)
					if err != nil {
						lua.Error(state, i.Lg.Appendf("failed to encode image: %s", log.LEVEL_ERROR, err))
					}

					var out string
					if raw {
						out = base64.RawStdEncoding.EncodeToString(wb.Bytes())
					} else {
						out = base64.StdEncoding.EncodeToString(wb.Bytes())
					}

					_, err = f.WriteString(out)
					if err != nil {
						lua.Error(state, i.Lg.Appendf("failed to write to file: %s", log.LEVEL_ERROR, err))
					}
				},
			})

			return 0
		})

	/// @func decode_image(data, name, encoding, model?, raw?) -> int<collection.IMAGE>
	/// @arg data {string} - The base64 encoded image data.
	/// @arg name {string} - The name of the image.
	/// @arg encoding {int<image.ENCODING>} - The image encoding must match the encoding of the data.
	/// @arg? model {int<image.MODEL>} - Used only to specify default when there is an unsupported color model.
	/// @arg? raw {bool} - If true, use raw encoding.
	/// @returns {int<collection.IMAGE>} - The decoded image.
	lib.CreateFunction(tab, "decode_image",
		[]lua.Arg{
			{Type: lua.STRING, Name: "data"},
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "encoding"},
			{Type: lua.INT, Name: "model", Optional: true},
			{Type: lua.BOOL, Name: "raw", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			data := args["data"].(string)
			name := args["name"].(string)
			encoding := lua.ParseEnum(args["encoding"].(int), imageutil.EncodingList, lib)
			model := lua.ParseEnum(args["model"].(int), imageutil.ModelList, lib)
			raw := args["raw"].(bool)

			chLog := log.NewLogger(fmt.Sprintf("image_%s", name), lg)
			lg.Append(fmt.Sprintf("child log created: image_%s", name), log.LEVEL_INFO)

			id := r.IC.AddItem(&chLog)

			r.IC.Schedule(id, &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					var dd []byte
					var err error

					if raw {
						dd, err = base64.RawStdEncoding.DecodeString(data)
					} else {
						dd, err = base64.StdEncoding.DecodeString(data)
					}
					if err != nil {
						lua.Error(state, i.Lg.Appendf("failed to decode base64 data: %s", log.LEVEL_ERROR, err))
					}

					bs := byteseeker.NewByteSeekerFromBytes(dd, 1000, true)
					img, err := imageutil.Decode(bs, encoding)
					if err != nil {
						lua.Error(state, i.Lg.Appendf("failed to decode image: %s", log.LEVEL_ERROR, err))
					}

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

	/// @func decode_image_from_file(path, name, encoding, model?, raw?) -> int<collection.IMAGE>
	/// @arg path {string} - The path to the base64 encoded image data.
	/// @arg name {string} - The name of the image.
	/// @arg encoding {int<image.ENCODING>} - The image encoding must match the encoding of the data.
	/// @arg? model {int<image.MODEL>} - Used only to specify default when there is an unsupported color model.
	/// @arg? raw {bool} - If true, use raw encoding.
	/// @returns {int<collection.IMAGE>} - The decoded image.
	lib.CreateFunction(tab, "decode_image_from_file",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "encoding"},
			{Type: lua.INT, Name: "model", Optional: true},
			{Type: lua.BOOL, Name: "raw", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			pth := args["path"].(string)
			name := args["name"].(string)
			encoding := lua.ParseEnum(args["encoding"].(int), imageutil.EncodingList, lib)
			model := lua.ParseEnum(args["model"].(int), imageutil.ModelList, lib)
			raw := args["raw"].(bool)

			chLog := log.NewLogger(fmt.Sprintf("image_%s", name), lg)
			lg.Append(fmt.Sprintf("child log created: image_%s", name), log.LEVEL_INFO)

			id := r.IC.AddItem(&chLog)

			r.IC.Schedule(id, &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					data, err := os.ReadFile(pth)
					if err != nil {
						lua.Error(state, i.Lg.Appendf("failed to read file: %s", log.LEVEL_ERROR, err))
					}

					var dd []byte

					if raw {
						dd, err = base64.RawStdEncoding.DecodeString(string(data))
					} else {
						dd, err = base64.StdEncoding.DecodeString(string(data))
					}
					if err != nil {
						lua.Error(state, i.Lg.Appendf("failed to decode base64 data: %s", log.LEVEL_ERROR, err))
					}

					bs := byteseeker.NewByteSeekerFromBytes(dd, 1000, true)
					img, err := imageutil.Decode(bs, encoding)
					if err != nil {
						lua.Error(state, i.Lg.Appendf("failed to decode image: %s", log.LEVEL_ERROR, err))
					}

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
}
