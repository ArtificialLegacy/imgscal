package lib

import (
	"os"

	"github.com/ArtificialLegacy/imgscal/pkg/image"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	golua "github.com/Shopify/go-lua"
	"github.com/qeesung/image2ascii/convert"
)

const LIB_ASCII = "ascii"

func RegisterASCII(r *lua.Runner, lg *log.Logger) {
	r.State.NewTable()

	/// @func to_file()
	/// @arg image_id
	/// @arg filepath - directories to file must exist.
	/// @arg color - boolean, for terminal dislay
	/// @arg reverse - boolean
	r.State.PushGoFunction(func(state *golua.State) int {
		lg.Append("ascii.to_file called", log.LEVEL_INFO)

		id, ok := r.State.ToInteger(-4)
		if !ok {
			r.State.PushString(lg.Append("invalid image id provided to ascii.to_file", log.LEVEL_ERROR))
			r.State.Error()
		}

		pth, ok := r.State.ToString(-3)
		if !ok {
			r.State.PushString(lg.Append("invalid file path provided to ascii.to_file", log.LEVEL_ERROR))
			r.State.Error()
		}

		color := r.State.ToBoolean(-2)
		reverse := r.State.ToBoolean(-1)

		r.IC.Schedule(id, &image.ImageTask{
			Fn: func(i *image.Image) {
				lg.Append("ascii.to_file task called", log.LEVEL_INFO)

				converter := convert.NewImageConverter()
				str := converter.Image2ASCIIString(i.Img, &convert.Options{
					Colored:  color,
					Reversed: reverse,
				})

				f, err := os.OpenFile(pth, os.O_CREATE|os.O_TRUNC, 0o666)
				if err != nil {
					r.State.PushString(lg.Append("failed to open file for saving ascii string", log.LEVEL_ERROR))
					r.State.Error()
				}
				defer f.Close()

				f.WriteString(str)

				lg.Append("ascii.to_file task finished", log.LEVEL_INFO)
			},
		})

		return 0
	})
	r.State.SetField(-2, "to_file")

	/// @func to_file_size()
	/// @arg image_id
	/// @arg filepath - directories to file must exist.
	/// @arg width
	/// @arg height
	/// @arg color - boolean, for terminal dislay
	/// @arg reverse - boolean
	r.State.PushGoFunction(func(state *golua.State) int {
		lg.Append("ascii.to_file_size called", log.LEVEL_INFO)

		id, ok := r.State.ToInteger(-6)
		if !ok {
			r.State.PushString(lg.Append("invalid image id provided to ascii.to_file_size", log.LEVEL_ERROR))
			r.State.Error()
		}

		pth, ok := r.State.ToString(-5)
		if !ok {
			r.State.PushString(lg.Append("invalid file path provided to ascii.to_file_size", log.LEVEL_ERROR))
			r.State.Error()
		}

		width, ok := r.State.ToInteger(-4)
		if !ok {
			r.State.PushString(lg.Append("invalid width provided to ascii.to_file_size", log.LEVEL_ERROR))
			r.State.Error()
		}

		height, ok := r.State.ToInteger(-3)
		if !ok {
			r.State.PushString(lg.Append("invalid height provided to ascii.to_file_size", log.LEVEL_ERROR))
			r.State.Error()
		}

		color := r.State.ToBoolean(-2)
		reverse := r.State.ToBoolean(-1)

		r.IC.Schedule(id, &image.ImageTask{
			Fn: func(i *image.Image) {
				lg.Append("ascii.to_file_size task called", log.LEVEL_INFO)

				converter := convert.NewImageConverter()
				str := converter.Image2ASCIIString(i.Img, &convert.Options{
					FixedWidth:  width,
					FixedHeight: height,
					Colored:     color,
					Reversed:    reverse,
				})

				f, err := os.OpenFile(pth, os.O_CREATE|os.O_TRUNC, 0o666)
				if err != nil {
					r.State.PushString(lg.Append("failed to open file for saving ascii string", log.LEVEL_ERROR))
					r.State.Error()
				}
				defer f.Close()

				f.WriteString(str)

				lg.Append("ascii.to_file_size task finished", log.LEVEL_INFO)
			},
		})

		return 0
	})
	r.State.SetField(-2, "to_file_size")

	/// @func to_string()
	/// @arg image_id
	/// @arg color - boolean, for terminal dislay
	/// @arg reverse - boolean
	/// @returns the ascii art as a string
	/// @blocking
	r.State.PushGoFunction(func(state *golua.State) int {
		lg.Append("ascii.to_string called", log.LEVEL_INFO)

		id, ok := r.State.ToInteger(-3)
		if !ok {
			r.State.PushString(lg.Append("invalid image id provided to ascii.to_string", log.LEVEL_ERROR))
			r.State.Error()
		}

		color := r.State.ToBoolean(-2)
		reverse := r.State.ToBoolean(-1)

		str := ""
		wait := make(chan bool)

		r.IC.Schedule(id, &image.ImageTask{
			Fn: func(i *image.Image) {
				lg.Append("ascii.to_string task called", log.LEVEL_INFO)

				converter := convert.NewImageConverter()
				str = converter.Image2ASCIIString(i.Img, &convert.Options{
					Colored:  color,
					Reversed: reverse,
				})

				wait <- true

				lg.Append("ascii.to_string task finished", log.LEVEL_INFO)
			},
		})

		<-wait
		r.State.PushString(str)
		return 0
	})
	r.State.SetField(-2, "to_string")

	/// @func to_string_size()
	/// @arg image_id
	/// @arg width
	/// @arg height
	/// @arg color - boolean, for terminal dislay
	/// @arg reverse - boolean
	/// @returns the ascii art as a string
	/// @blocking
	r.State.PushGoFunction(func(state *golua.State) int {
		lg.Append("ascii.to_string_size called", log.LEVEL_INFO)

		id, ok := r.State.ToInteger(-5)
		if !ok {
			r.State.PushString(lg.Append("invalid image id provided to ascii.to_string_size", log.LEVEL_ERROR))
			r.State.Error()
		}

		width, ok := r.State.ToInteger(-4)
		if !ok {
			r.State.PushString(lg.Append("invalid width provided to ascii.to_string_size", log.LEVEL_ERROR))
			r.State.Error()
		}

		height, ok := r.State.ToInteger(-3)
		if !ok {
			r.State.PushString(lg.Append("invalid height provided to ascii.to_string_size", log.LEVEL_ERROR))
			r.State.Error()
		}

		color := r.State.ToBoolean(-2)
		reverse := r.State.ToBoolean(-1)

		str := ""
		wait := make(chan bool)

		r.IC.Schedule(id, &image.ImageTask{
			Fn: func(i *image.Image) {
				lg.Append("ascii.to_string_size task called", log.LEVEL_INFO)

				converter := convert.NewImageConverter()
				str = converter.Image2ASCIIString(i.Img, &convert.Options{
					FixedWidth:  width,
					FixedHeight: height,
					Colored:     color,
					Reversed:    reverse,
				})

				wait <- true

				lg.Append("ascii.to_string_size task finished", log.LEVEL_INFO)
			},
		})

		<-wait
		r.State.PushString(str)
		return 0
	})
	r.State.SetField(-2, "to_string_size")

	r.State.SetGlobal(LIB_ASCII)
}
