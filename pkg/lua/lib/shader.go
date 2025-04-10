package lib

import (
	"image"
	"os"

	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	imageutil "github.com/ArtificialLegacy/imgscal/pkg/image_util"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	golua "github.com/yuin/gopher-lua"
	"gitlab.com/microo8/blackcl"
)

const LIB_SHADER = "shader"

/// @lib Shader
/// @import shader
/// @desc
/// Library for running GPU shaders.

func RegisterShader(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_SHADER, r, r.State, lg)

	/// @func device() -> int<collection.CRATE_SHADER>
	/// @returns {int<collection.CRATE_SHADER>}
	lib.CreateFunction(tab, "device",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			devices, err := blackcl.GetDevices(blackcl.DeviceTypeDefault)
			if err != nil {
				lua.Error(state, lg.Appendf("failed to get shader device: %s", log.LEVEL_ERROR, err))
			}
			if len(devices) == 0 {
				lua.Error(state, lg.Append("no devices found", log.LEVEL_ERROR))
			}

			device := devices[0]
			id := r.CR_SHD.Add(&collection.ShaderItem{
				Device:  device,
				Kernels: []*blackcl.KernelCall{},

				BuffersImage:  []*blackcl.Image{},
				BuffersVector: []*blackcl.Vector{},
				BuffersBytes:  []*blackcl.Bytes{},
			})

			if len(devices) > 1 {
				for _, d := range devices[1:] {
					d.Release()
				}
			}

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func release(device)
	/// @arg id {int<collection.CRATE_SHADER>}
	/// @desc
	/// After a device is released, it and any associated buffers and kernels cannot be used.
	lib.CreateFunction(tab, "release",
		[]lua.Arg{
			{Type: lua.INT, Name: "device"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CR_SHD.Clean(args["device"].(int))
			return 0
		})

	/// @func add_program(device, program)
	/// @arg id {int<collection.CRATE_SHADER>}
	/// @arg program {string}
	lib.CreateFunction(tab, "add_program",
		[]lua.Arg{
			{Type: lua.INT, Name: "device"},
			{Type: lua.STRING, Name: "program"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			shd, err := r.CR_SHD.Item(args["device"].(int))
			if err != nil {
				lua.Error(state, lg.Append("could not retrieve shader device", log.LEVEL_ERROR))
			}

			shd.Device.AddProgram(args["program"].(string))

			return 0
		})

	/// @func add_program_file(device, path)
	/// @arg device {int<collection.CRATE_SHADER>}
	/// @arg path{string}
	lib.CreateFunction(tab, "add_program_file",
		[]lua.Arg{
			{Type: lua.INT, Name: "device"},
			{Type: lua.STRING, Name: "path"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			shd, err := r.CR_SHD.Item(args["device"].(int))
			if err != nil {
				lua.Error(state, lg.Append("could not retrieve shader device", log.LEVEL_ERROR))
			}

			b, err := os.ReadFile(args["path"].(string))
			if err != nil {
				lua.Error(state, lg.Appendf("could not open shader file: %s", log.LEVEL_ERROR, err))
			}

			shd.Device.AddProgram(string(b))

			return 0
		})

	/// @func kernel(device, name, []globals, []locals) -> int<shader.KERNEL>
	/// @arg device {int<collection.CRATE_SHADER>}
	/// @arg name {string}
	/// @arg globals {[]int} - List of the global work sizes to use.
	/// @arg locals {[]int} - List of the local work sizes to use.
	/// @returns int<shader.KERNEL>
	lib.CreateFunction(tab, "kernel",
		[]lua.Arg{
			{Type: lua.INT, Name: "device"},
			{Type: lua.STRING, Name: "name"},
			lua.ArgArray("global", lua.ArrayType{Type: lua.INT}, false),
			lua.ArgArray("local", lua.ArrayType{Type: lua.INT}, false),
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			shd, err := r.CR_SHD.Item(args["device"].(int))
			if err != nil {
				lua.Error(state, lg.Append("could not retrieve shader device", log.LEVEL_ERROR))
			}

			name := args["name"].(string)

			glb := args["global"].([]any)
			global := make([]int, len(glb))
			for i, v := range glb {
				global[i] = v.(int)
			}

			lcl := args["local"].([]any)
			local := make([]int, len(lcl))
			for i, v := range lcl {
				local[i] = v.(int)
			}

			kernel := shd.Device.Kernel(name).Global(global...).Local(local...)
			id := len(shd.Kernels)
			shd.Kernels = append(shd.Kernels, &kernel)

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func run(device, kernel, buffers...)
	/// @arg device {int<collection.CRATE_SHADER>}
	/// @arg kernel {int<shader.KERNEL>}
	/// @arg buffers {struct<shader.Buffer>...}
	lib.CreateFunction(tab, "run",
		[]lua.Arg{
			{Type: lua.INT, Name: "device"},
			{Type: lua.INT, Name: "kernel"},
			lua.ArgVariadic("buffers", lua.ArrayType{Type: lua.RAW_TABLE}, true),
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			shd, err := r.CR_SHD.Item(args["device"].(int))
			if err != nil {
				lua.Error(state, lg.Append("could not retrieve shader device", log.LEVEL_ERROR))
			}

			kernelid := args["kernel"].(int)
			if kernelid < 0 || kernelid >= len(shd.Kernels) {
				lua.Error(state, lg.Appendf("kernel id %d out of range 0-%d", log.LEVEL_ERROR, kernelid, len(shd.Kernels)))
			}

			kernel := shd.Kernels[kernelid]
			buffers := getBuffers(shd, args["buffers"].([]any))

			err = <-kernel.Run(buffers...)
			if err != nil {
				lua.Error(state, lg.Appendf("failed to run kernel: %s", log.LEVEL_ERROR, err))
			}

			return 0
		})

	/// @func buffer_image(device, type, width, height) -> struct<shader.BufferImage>
	/// @arg device {int<collection.CRATE_SHADER>}
	/// @arg type {int<shader.ImageType>}
	/// @arg width {int}
	/// @arg height {int}
	/// @returns {struct<shader.BufferImage>}
	lib.CreateFunction(tab, "buffer_image",
		[]lua.Arg{
			{Type: lua.INT, Name: "device"},
			{Type: lua.INT, Name: "type"},
			{Type: lua.INT, Name: "width"},
			{Type: lua.INT, Name: "height"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			device := args["device"].(int)
			shd, err := r.CR_SHD.Item(device)
			if err != nil {
				lua.Error(state, lg.Append("could not retrieve shader device", log.LEVEL_ERROR))
			}

			width := args["width"].(int)
			height := args["height"].(int)

			buff, err := shd.Device.NewImage(blackcl.ImageType(args["type"].(int)), image.Rect(0, 0, width, height))
			if err != nil {
				lua.Error(state, lg.Appendf("failed to create new image buffer: %s", log.LEVEL_ERROR, err))
			}

			id := len(shd.BuffersImage)
			shd.BuffersImage = append(shd.BuffersImage, buff)

			t := bufferImage(r, lib, state, lg, device, id)

			state.Push(t)
			return 1
		})

	/// @func buffer_image_ext(device, type, p1, p2) -> struct<shader.BufferImage>
	/// @arg device {int<collection.CRATE_SHADER>}
	/// @arg type {int<shader.ImageType>}
	/// @arg p1 {struct<image.Point>}
	/// @arg p2 {struct<image.Point>}
	/// @returns {struct<shader.BufferImage>}
	lib.CreateFunction(tab, "buffer_image_ext",
		[]lua.Arg{
			{Type: lua.INT, Name: "device"},
			{Type: lua.INT, Name: "type"},
			{Type: lua.RAW_TABLE, Name: "p1"},
			{Type: lua.RAW_TABLE, Name: "p2"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			device := args["device"].(int)
			shd, err := r.CR_SHD.Item(device)
			if err != nil {
				lua.Error(state, lg.Append("could not retrieve shader device", log.LEVEL_ERROR))
			}

			p1 := imageutil.TableToPoint(args["p1"].(*golua.LTable))
			p2 := imageutil.TableToPoint(args["p2"].(*golua.LTable))

			buff, err := shd.Device.NewImage(blackcl.ImageType(args["type"].(int)), image.Rect(p1.X, p1.Y, p2.X, p2.Y))
			if err != nil {
				lua.Error(state, lg.Appendf("failed to create new image buffer: %s", log.LEVEL_ERROR, err))
			}

			id := len(shd.BuffersImage)
			shd.BuffersImage = append(shd.BuffersImage, buff)

			t := bufferImage(r, lib, state, lg, device, id)

			state.Push(t)
			return 1
		})

	/// @func buffer_image_ext_xy(device, type, x1, y1, x2, y2) -> struct<shader.BufferImage>
	/// @arg device {int<collection.CRATE_SHADER>}
	/// @arg type {int<shader.ImageType>}
	/// @arg x1 {int}
	/// @arg y1 {int}
	/// @arg x2 {int}
	/// @arg y2 {int}
	/// @returns {struct<shader.BufferImage>}
	lib.CreateFunction(tab, "buffer_image_ext_xy",
		[]lua.Arg{
			{Type: lua.INT, Name: "device"},
			{Type: lua.INT, Name: "type"},
			{Type: lua.INT, Name: "x1"},
			{Type: lua.INT, Name: "y1"},
			{Type: lua.INT, Name: "x2"},
			{Type: lua.INT, Name: "y2"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			device := args["device"].(int)
			shd, err := r.CR_SHD.Item(device)
			if err != nil {
				lua.Error(state, lg.Append("could not retrieve shader device", log.LEVEL_ERROR))
			}

			x1 := args["x1"].(int)
			y1 := args["y1"].(int)
			x2 := args["x2"].(int)
			y2 := args["y2"].(int)

			buff, err := shd.Device.NewImage(blackcl.ImageType(args["type"].(int)), image.Rect(x1, y1, x2, y2))
			if err != nil {
				lua.Error(state, lg.Appendf("failed to create new image buffer: %s", log.LEVEL_ERROR, err))
			}

			id := len(shd.BuffersImage)
			shd.BuffersImage = append(shd.BuffersImage, buff)

			t := bufferImage(r, lib, state, lg, device, id)

			state.Push(t)
			return 1
		})

	/// @func buffer_image_from(device, id) -> struct<shader.BufferImage>
	/// @arg device {int<collection.CRATE_SHADER>}
	/// @arg id {int<collection.IMAGE>}
	/// @returns {struct<shader.BufferImage>}
	/// @blocking
	lib.CreateFunction(tab, "buffer_image_from",
		[]lua.Arg{
			{Type: lua.INT, Name: "device"},
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var img image.Image

			device := args["device"].(int)
			shd, err := r.CR_SHD.Item(device)
			if err != nil {
				lua.Error(state, lg.Append("could not retrieve shader device", log.LEVEL_ERROR))
			}

			<-r.IC.Schedule(state, args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					img = i.Self.Image
				},
			})

			buff, err := shd.Device.NewImageFromImage(img)
			if err != nil {
				lua.Error(state, lg.Appendf("failed to create new image buffer: %s", log.LEVEL_ERROR, err))
			}

			id := len(shd.BuffersImage)
			shd.BuffersImage = append(shd.BuffersImage, buff)

			t := bufferImage(r, lib, state, lg, device, id)

			state.Push(t)
			return 1
		})

	/// @func buffer_vector(device, length) -> struct<shader.BufferVector>
	/// @arg device {int<collection.CRATE_SHADER>}
	/// @arg length {int}
	/// @returns {struct<shader.BufferVector>}
	lib.CreateFunction(tab, "buffer_vector",
		[]lua.Arg{
			{Type: lua.INT, Name: "device"},
			{Type: lua.INT, Name: "length"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			device := args["device"].(int)
			shd, err := r.CR_SHD.Item(device)
			if err != nil {
				lua.Error(state, lg.Append("could not retrieve shader device", log.LEVEL_ERROR))
			}

			buff, err := shd.Device.NewVector(args["length"].(int))
			if err != nil {
				lua.Error(state, lg.Appendf("failed to create new vector buffer: %s", log.LEVEL_ERROR, err))
			}

			id := len(shd.BuffersVector)
			shd.BuffersVector = append(shd.BuffersVector, buff)

			t := bufferVector(r, lib, state, lg, device, id)

			state.Push(t)
			return 1
		})

	/// @func buffer_bytes(device, size) -> struct<shader.BufferBytes>
	/// @arg device {int<collection.CRATE_SHADER>}
	/// @arg size {int}
	/// @returns {struct<shader.BufferBytes>}
	lib.CreateFunction(tab, "buffer_bytes",
		[]lua.Arg{
			{Type: lua.INT, Name: "device"},
			{Type: lua.INT, Name: "bytes"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			device := args["device"].(int)
			shd, err := r.CR_SHD.Item(device)
			if err != nil {
				lua.Error(state, lg.Append("could not retrieve shader device", log.LEVEL_ERROR))
			}

			buff, err := shd.Device.NewBytes(args["size"].(int))
			if err != nil {
				lua.Error(state, lg.Appendf("failed to create new bytes buffer: %s", log.LEVEL_ERROR, err))
			}

			id := len(shd.BuffersBytes)
			shd.BuffersBytes = append(shd.BuffersBytes, buff)

			t := bufferBytes(r, lib, state, lg, device, id)

			state.Push(t)
			return 1
		})

	/// @constants BufferType {int}
	/// @const BUFFER_IMAGE
	/// @const BUFFER_VECTOR
	/// @const BUFFER_BYTES
	tab.RawSetString("BUFFER_IMAGE", golua.LNumber(buffertype_image))
	tab.RawSetString("BUFFER_VECTOR", golua.LNumber(buffertype_vector))
	tab.RawSetString("BUFFER_BYTES", golua.LNumber(buffertype_bytes))

	/// @constants ImageType {int}
	/// @const IMAGE_RGBA
	/// @const IMAGE_GRAY
	tab.RawSetString("IMAGE_RGBA", golua.LNumber(blackcl.ImageTypeRGBA))
	tab.RawSetString("IMAGE_GRAY", golua.LNumber(blackcl.ImageTypeGray))

	/// @constants ID {int}
	/// @const KERNEL
	tab.RawSetString("KERNEL", golua.LNumber(0))
}

const (
	buffertype_image int = iota
	buffertype_vector
	buffertype_bytes
)

func getBuffers(shd *collection.ShaderItem, buffers []any) []any {
	/// @interface Buffer
	/// @prop type {int<shader.BufferType>}
	/// @prop device {int<collection.CRATE_SHADER>}
	/// @prop id {int}

	result := make([]any, len(buffers))

	for i, v := range buffers {
		buff := v.(*golua.LTable)

		typ := int(buff.RawGetString("type").(golua.LNumber))
		id := int(buff.RawGetString("id").(golua.LNumber))

		switch typ {
		case buffertype_image:
			result[i] = shd.BuffersImage[id]
		case buffertype_vector:
			result[i] = shd.BuffersVector[id]
		case buffertype_bytes:
			result[i] = shd.BuffersBytes[id]
		}
	}

	return result
}

func bufferImage(r *lua.Runner, lib *lua.Lib, state *golua.LState, lg *log.Logger, device, id int) *golua.LTable {
	/// @struct BufferImage
	/// @prop type {int<shader.BufferType>}
	/// @prop device {int<collection.CRATE_SHADER>}
	/// @prop id {int}
	/// @method bounds() -> int, int
	/// @method copy(self, int<collection.IMAGE>) -> self
	/// @method data(name string, int<image.Encoding>) -> int<collection.IMAGE>
	/// @method data_into(self, int<collection.IMAGE>) -> self

	t := state.NewTable()

	t.RawSetString("type", golua.LNumber(buffertype_image))
	t.RawSetString("device", golua.LNumber(device))
	t.RawSetString("id", golua.LNumber(id))

	lib.TableFunction(state, t, "bounds",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			shd, err := r.CR_SHD.Item(int(t.RawGetString("device").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append("could not retrieve shader device", log.LEVEL_ERROR))
			}

			buff := shd.BuffersImage[int(t.RawGetString("id").(golua.LNumber))]
			rect := buff.Bounds()

			state.Push(golua.LNumber(rect.Dx()))
			state.Push(golua.LNumber(rect.Dy()))
			return 2
		})

	lib.BuilderFunction(state, t, "copy",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			var img image.Image

			<-r.IC.Schedule(state, args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  LIB_SHADER,
				Name: "BufferImage.copy",
				Fn: func(i *collection.Item[collection.ItemImage]) {
					img = i.Self.Image
				},
			})

			shd, err := r.CR_SHD.Item(int(t.RawGetString("device").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append("could not retrieve shader device", log.LEVEL_ERROR))
			}

			buff := shd.BuffersImage[int(t.RawGetString("id").(golua.LNumber))]

			err = <-buff.Copy(img)
			if err != nil {
				lua.Error(state, lg.Appendf("failed to copy image data to buffer: %s", log.LEVEL_ERROR, err))
			}
		})

	lib.TableFunction(state, t, "data",
		[]lua.Arg{
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "encoding"},
		},
		func(state *golua.LState, args map[string]any) int {
			shd, err := r.CR_SHD.Item(int(t.RawGetString("device").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append("could not retrieve shader device", log.LEVEL_ERROR))
			}

			buff := shd.BuffersImage[int(t.RawGetString("id").(golua.LNumber))]
			img, err := buff.Data()
			if err != nil {
				lua.Error(state, lg.Appendf("failed to get image data from buffer: %s", log.LEVEL_ERROR, err))
			}

			name := args["name"].(string)
			id := r.IC.ScheduleAdd(state, name, lg, LIB_SHADER, "BufferImage.data", func(i *collection.Item[collection.ItemImage]) {
				var model imageutil.ColorModel
				if _, ok := img.(*image.RGBA); ok {
					model = imageutil.MODEL_RGBA
				} else if _, ok := img.(*image.Gray); ok {
					model = imageutil.MODEL_GRAY
				} else {
					img = imageutil.CopyImage(img, imageutil.MODEL_RGBA)
					model = imageutil.MODEL_RGBA
				}

				i.Self = &collection.ItemImage{
					Image:    img,
					Name:     name,
					Encoding: lua.ParseEnum(args["encoding"].(int), imageutil.EncodingList, lib),
					Model:    model,
				}
			})

			state.Push(golua.LNumber(id))
			return 1
		})

	lib.BuilderFunction(state, t, "data_into",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			shd, err := r.CR_SHD.Item(int(t.RawGetString("device").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append("could not retrieve shader device", log.LEVEL_ERROR))
			}

			buff := shd.BuffersImage[int(t.RawGetString("id").(golua.LNumber))]
			img, err := buff.Data()
			if err != nil {
				lua.Error(state, lg.Appendf("failed to get image data from buffer: %s", log.LEVEL_ERROR, err))
			}

			r.IC.Schedule(state, args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  LIB_SHADER,
				Name: "BufferImage.data_into",
				Fn: func(i *collection.Item[collection.ItemImage]) {
					var model imageutil.ColorModel
					if _, ok := img.(*image.RGBA); ok {
						model = imageutil.MODEL_RGBA
					} else if _, ok := img.(*image.Gray); ok {
						model = imageutil.MODEL_GRAY
					} else {
						img = imageutil.CopyImage(img, imageutil.MODEL_RGBA)
						model = imageutil.MODEL_RGBA
					}

					i.Self.Image = img
					i.Self.Model = model
				},
			})
		})

	return t
}

func bufferVector(r *lua.Runner, lib *lua.Lib, state *golua.LState, lg *log.Logger, device, id int) *golua.LTable {
	/// @struct BufferVector
	/// @prop type {int<shader.BufferType>}
	/// @prop device {int<collection.CRATE_SHADER>}
	/// @prop id {int}
	/// @method copy(self, []float) -> self
	/// @method data() -> []float
	/// @method length() -> int

	t := state.NewTable()

	t.RawSetString("type", golua.LNumber(buffertype_vector))
	t.RawSetString("device", golua.LNumber(device))
	t.RawSetString("id", golua.LNumber(id))

	lib.BuilderFunction(state, t, "copy",
		[]lua.Arg{
			lua.ArgArray("data", lua.ArrayType{Type: lua.FLOAT}, false),
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			shd, err := r.CR_SHD.Item(int(t.RawGetString("device").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append("could not retrieve shader device", log.LEVEL_ERROR))
			}

			buff := shd.BuffersVector[int(t.RawGetString("id").(golua.LNumber))]

			dataraw := args["data"].([]any)
			data := make([]float32, len(dataraw))
			for i, v := range dataraw {
				data[i] = float32(v.(float64))
			}

			err = <-buff.Copy(data)
			if err != nil {
				lua.Error(state, lg.Appendf("failed to copy vector data to buffer: %s", log.LEVEL_ERROR, err))
			}
		})

	lib.TableFunction(state, t, "data",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			shd, err := r.CR_SHD.Item(int(t.RawGetString("device").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append("could not retrieve shader device", log.LEVEL_ERROR))
			}

			buff := shd.BuffersVector[int(t.RawGetString("id").(golua.LNumber))]

			data, err := buff.Data()
			if err != nil {
				lua.Error(state, lg.Appendf("failed to get vector data from buffer: %s", log.LEVEL_ERROR, err))
			}

			datalist := state.NewTable()
			for i, v := range data {
				datalist.RawSetInt(i+1, golua.LNumber(v))
			}

			state.Push(datalist)
			return 1
		})

	lib.TableFunction(state, t, "length",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			shd, err := r.CR_SHD.Item(int(t.RawGetString("device").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append("could not retrieve shader device", log.LEVEL_ERROR))
			}

			buff := shd.BuffersVector[int(t.RawGetString("id").(golua.LNumber))]

			length := buff.Length()
			state.Push(golua.LNumber(length))
			return 1
		})

	return t
}

func bufferBytes(r *lua.Runner, lib *lua.Lib, state *golua.LState, lg *log.Logger, device, id int) *golua.LTable {
	/// @struct BufferBytes
	/// @prop type {int<shader.BufferType>}
	/// @prop device {int<collection.CRATE_SHADER>}
	/// @prop id {int}
	/// @method copy(self, []int) -> self
	/// @method data() -> []int
	/// @method size() -> int

	t := state.NewTable()

	t.RawSetString("type", golua.LNumber(buffertype_bytes))
	t.RawSetString("device", golua.LNumber(device))
	t.RawSetString("id", golua.LNumber(id))

	lib.BuilderFunction(state, t, "copy",
		[]lua.Arg{
			lua.ArgArray("data", lua.ArrayType{Type: lua.INT}, false),
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			shd, err := r.CR_SHD.Item(int(t.RawGetString("device").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append("could not retrieve shader device", log.LEVEL_ERROR))
			}

			buff := shd.BuffersBytes[int(t.RawGetString("id").(golua.LNumber))]

			dataraw := args["data"].([]any)
			data := make([]byte, len(dataraw))
			for i, v := range dataraw {
				data[i] = byte(v.(int))
			}

			err = <-buff.Copy(data)
			if err != nil {
				lua.Error(state, lg.Appendf("failed to copy bytes data to buffer: %s", log.LEVEL_ERROR, err))
			}
		})

	lib.TableFunction(state, t, "data",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			shd, err := r.CR_SHD.Item(int(t.RawGetString("device").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append("could not retrieve shader device", log.LEVEL_ERROR))
			}

			buff := shd.BuffersBytes[int(t.RawGetString("id").(golua.LNumber))]

			data, err := buff.Data()
			if err != nil {
				lua.Error(state, lg.Appendf("failed to get bytes data from buffer: %s", log.LEVEL_ERROR, err))
			}

			datalist := state.NewTable()
			for i, v := range data {
				datalist.RawSetInt(i+1, golua.LNumber(v))
			}

			state.Push(datalist)
			return 1
		})

	lib.TableFunction(state, t, "size",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			shd, err := r.CR_SHD.Item(int(t.RawGetString("device").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append("could not retrieve shader device", log.LEVEL_ERROR))
			}

			buff := shd.BuffersBytes[int(t.RawGetString("id").(golua.LNumber))]

			size := buff.Size()
			state.Push(golua.LNumber(size))
			return 1
		})

	return t
}
