package lib

import (
	"fmt"
	"image"
	"image/color"

	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	imageutil "github.com/ArtificialLegacy/imgscal/pkg/image_util"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	"github.com/skip2/go-qrcode"
)

const LIB_QRCODE = "qrcode"

func RegisterQRCode(r *lua.Runner, lg *log.Logger) {
	lib := lua.NewLib(LIB_QRCODE, r.State, lg)

	/// @func new()
	/// @arg content
	/// @arg recovery
	/// @returns id
	lib.CreateFunction("new",
		[]lua.Arg{
			{Type: lua.STRING, Name: "content"},
			{Type: lua.INT, Name: "recovery"},
		},
		func(d lua.TaskData, args map[string]any) int {
			chLog := log.NewLogger(fmt.Sprintf("qrcode_%d", r.QR.Next()))
			chLog.Parent = lg
			lg.Append(fmt.Sprintf("child log created: qrcode_%d", r.QR.Next()), log.LEVEL_INFO)

			id := r.QR.AddItem(&chLog)

			r.QR.Schedule(id, &collection.Task[collection.ItemQR]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemQR]) {
					qr, err := qrcode.New(args["content"].(string), lua.ParseEnum(args["recovery"].(int), recoveryLevelList, lib))
					if err != nil {
						r.State.PushString(lg.Append(fmt.Sprintf("unable to create qrcode: %s", err), log.LEVEL_ERROR))
						r.State.Error()
					}

					i.Self = &collection.ItemQR{
						QR: qr,
					}
				},
			})

			r.State.PushInteger(id)
			return 1
		})

	/// @func to_image()
	/// @arg id
	/// @arg name
	/// @arg size - positive sets a fixed size, negative sets a scaled size
	/// @arg encoding
	/// @returns image id
	lib.CreateFunction("to_image",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "size"},
			{Type: lua.INT, Name: "encoding"},
		},
		func(d lua.TaskData, args map[string]any) int {
			var img image.Image
			imgReady := make(chan struct{}, 1)

			r.QR.Schedule(args["id"].(int), &collection.Task[collection.ItemQR]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemQR]) {
					img = i.Self.QR.Image(args["size"].(int))
					imgReady <- struct{}{}
				},
				Fail: func(i *collection.Item[collection.ItemQR]) {
					imgReady <- struct{}{}
				},
			})

			name := args["name"].(string)

			chLog := log.NewLogger(fmt.Sprintf("image_%s", name))
			chLog.Parent = lg
			lg.Append(fmt.Sprintf("child log created: image_%s", name), log.LEVEL_INFO)

			id := r.IC.AddItem(&chLog)

			r.IC.Schedule(id, &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					<-imgReady
					i.Self = &collection.ItemImage{
						Image:    img,
						Name:     name,
						Encoding: lua.ParseEnum(args["encoding"].(int), imageutil.EncodingList, lib),
					}
				},
			})

			r.State.PushInteger(id)
			return 1
		})

	/// @func to_string()
	/// @arg id
	/// @arg inverse
	/// @returns string
	/// @blocking
	lib.CreateFunction("to_string",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.BOOL, Name: "inverse"},
		},
		func(d lua.TaskData, args map[string]any) int {
			var str string

			<-r.QR.Schedule(args["id"].(int), &collection.Task[collection.ItemQR]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemQR]) {
					str = i.Self.QR.ToString(args["inverse"].(bool))
				},
			})

			r.State.PushString(str)
			return 1
		})

	/// @func to_small_string()
	/// @arg id
	/// @arg inverse
	/// @returns string
	/// @blocking
	lib.CreateFunction("to_small_string",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.BOOL, Name: "inverse"},
		},
		func(d lua.TaskData, args map[string]any) int {
			var str string

			<-r.QR.Schedule(args["id"].(int), &collection.Task[collection.ItemQR]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemQR]) {
					str = i.Self.QR.ToSmallString(args["inverse"].(bool))
				},
			})

			r.State.PushString(str)
			return 1
		})

	/// @func color_set_foreground()
	/// @arg id
	/// @arg color {red, green, blue, alpha}
	lib.CreateFunction("color_set_foreground",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.TABLE, Name: "color", Table: &[]lua.Arg{
				{Type: lua.INT, Name: "red"},
				{Type: lua.INT, Name: "green"},
				{Type: lua.INT, Name: "blue"},
				{Type: lua.INT, Name: "alpha"},
			}},
		},
		func(d lua.TaskData, args map[string]any) int {
			r.QR.Schedule(args["id"].(int), &collection.Task[collection.ItemQR]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemQR]) {
					colors := args["color"].(map[string]any)
					col := color.RGBA{
						R: uint8(colors["red"].(int)),
						G: uint8(colors["green"].(int)),
						B: uint8(colors["blue"].(int)),
						A: uint8(colors["alpha"].(int)),
					}
					i.Self.QR.ForegroundColor = col
				},
			})
			return 0
		})

	/// @func color_set_background()
	/// @arg id
	/// @arg color {red, green, blue, alpha}
	lib.CreateFunction("color_set_background",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.TABLE, Name: "color", Table: &[]lua.Arg{
				{Type: lua.INT, Name: "red"},
				{Type: lua.INT, Name: "green"},
				{Type: lua.INT, Name: "blue"},
				{Type: lua.INT, Name: "alpha"},
			}},
		},
		func(d lua.TaskData, args map[string]any) int {
			r.QR.Schedule(args["id"].(int), &collection.Task[collection.ItemQR]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemQR]) {
					colors := args["color"].(map[string]any)
					col := color.RGBA{
						R: uint8(colors["red"].(int)),
						G: uint8(colors["green"].(int)),
						B: uint8(colors["blue"].(int)),
						A: uint8(colors["alpha"].(int)),
					}
					i.Self.QR.BackgroundColor = col
				},
			})
			return 0
		})

	/// @func color_foreground()
	/// @arg id
	/// @returns {red, green, blue, alpha}
	/// @blocking
	lib.CreateFunction("color_foreground",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(d lua.TaskData, args map[string]any) int {
			re := uint32(0)
			gr := uint32(0)
			bl := uint32(0)
			al := uint32(0)

			<-r.QR.Schedule(args["id"].(int), &collection.Task[collection.ItemQR]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemQR]) {
					re, gr, bl, al = i.Self.QR.ForegroundColor.RGBA()
				},
			})

			r.State.NewTable()
			r.State.PushInteger(int(re))
			r.State.SetField(-2, "red")
			r.State.PushInteger(int(gr))
			r.State.SetField(-2, "green")
			r.State.PushInteger(int(bl))
			r.State.SetField(-2, "blue")
			r.State.PushInteger(int(al))
			r.State.SetField(-2, "alpha")
			return 1
		})

	/// @func color_background()
	/// @arg id
	/// @returns {red, green, blue, alpha}
	/// @blocking
	lib.CreateFunction("color_background",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(d lua.TaskData, args map[string]any) int {
			re := uint32(0)
			gr := uint32(0)
			bl := uint32(0)
			al := uint32(0)

			<-r.QR.Schedule(args["id"].(int), &collection.Task[collection.ItemQR]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemQR]) {
					re, gr, bl, al = i.Self.QR.BackgroundColor.RGBA()
				},
			})

			r.State.NewTable()
			r.State.PushInteger(int(re))
			r.State.SetField(-2, "red")
			r.State.PushInteger(int(gr))
			r.State.SetField(-2, "green")
			r.State.PushInteger(int(bl))
			r.State.SetField(-2, "blue")
			r.State.PushInteger(int(al))
			r.State.SetField(-2, "alpha")
			return 1
		})

	/// @func border()
	/// @arg id
	/// @returns bool
	/// @blocking
	lib.CreateFunction("border",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(d lua.TaskData, args map[string]any) int {
			var active bool

			<-r.QR.Schedule(args["id"].(int), &collection.Task[collection.ItemQR]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemQR]) {
					active = !i.Self.QR.DisableBorder
				},
			})

			r.State.PushBoolean(active)
			return 1
		})

	/// @func recovery_level()
	/// @arg id
	/// @returns int
	/// @blocking
	lib.CreateFunction("recovery_level",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(d lua.TaskData, args map[string]any) int {
			var level qrcode.RecoveryLevel

			<-r.QR.Schedule(args["id"].(int), &collection.Task[collection.ItemQR]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemQR]) {
					level = i.Self.QR.Level
				},
			})

			r.State.PushInteger(int(level))
			return 1
		})

	/// @func version()
	/// @arg id
	/// @returns int
	/// @blocking
	lib.CreateFunction("version",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(d lua.TaskData, args map[string]any) int {
			var version int

			<-r.QR.Schedule(args["id"].(int), &collection.Task[collection.ItemQR]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemQR]) {
					version = i.Self.QR.VersionNumber
				},
			})

			r.State.PushInteger(int(version))
			return 1
		})

	/// @func border_set()
	/// @arg id
	/// @arg? border
	lib.CreateFunction("border_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.BOOL, Name: "border", Optional: true},
		},
		func(d lua.TaskData, args map[string]any) int {
			r.QR.Schedule(args["id"].(int), &collection.Task[collection.ItemQR]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemQR]) {
					i.Self.QR.DisableBorder = !args["border"].(bool)
				},
			})
			return 0
		})

	/// @func content_set()
	/// @arg id
	/// @arg content
	lib.CreateFunction("content_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "content"},
		},
		func(d lua.TaskData, args map[string]any) int {
			r.QR.Schedule(args["id"].(int), &collection.Task[collection.ItemQR]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemQR]) {
					i.Self.QR.Content = args["content"].(string)
				},
			})
			return 0
		})

	/// @func content()
	/// @arg id
	/// @returns string
	/// @blocking
	lib.CreateFunction("content",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(d lua.TaskData, args map[string]any) int {
			var content string

			<-r.QR.Schedule(args["id"].(int), &collection.Task[collection.ItemQR]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemQR]) {
					content = i.Self.QR.Content
				},
			})

			r.State.PushString(content)
			return 1
		})

	/// @constants Recovery Levels
	/// @const RECOVERY_LOW
	/// @const RECOVERY_MEDIUM
	/// @const RECOVERY_HIGH
	/// @const RECOVERY_HIGHEST
	r.State.PushInteger(int(qrcode.Low))
	r.State.SetField(-2, "RECOVERY_LOW")
	r.State.PushInteger(int(qrcode.Medium))
	r.State.SetField(-2, "RECOVERY_MEDIUM")
	r.State.PushInteger(int(qrcode.High))
	r.State.SetField(-2, "RECOVERY_HIGH")
	r.State.PushInteger(int(qrcode.Highest))
	r.State.SetField(-2, "RECOVERY_HIGHEST")
}

var recoveryLevelList = []qrcode.RecoveryLevel{
	qrcode.Low,
	qrcode.Medium,
	qrcode.High,
	qrcode.Highest,
}
