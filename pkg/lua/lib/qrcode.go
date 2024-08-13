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
	golua "github.com/yuin/gopher-lua"
)

const LIB_QRCODE = "qrcode"

func RegisterQRCode(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_QRCODE, r, r.State, lg)

	/// @func new()
	/// @arg content
	/// @arg recovery
	/// @returns id
	lib.CreateFunction(tab, "new",
		[]lua.Arg{
			{Type: lua.STRING, Name: "content"},
			{Type: lua.INT, Name: "recovery"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			chLog := log.NewLogger(fmt.Sprintf("qrcode_%d", r.QR.Next()))
			chLog.Parent(lg)
			lg.Append(fmt.Sprintf("child log created: qrcode_%d", r.QR.Next()), log.LEVEL_INFO)

			id := r.QR.AddItem(&chLog)

			r.QR.Schedule(id, &collection.Task[collection.ItemQR]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemQR]) {
					qr, err := qrcode.New(args["content"].(string), lua.ParseEnum(args["recovery"].(int), recoveryLevelList, lib))
					if err != nil {
						state.Error(golua.LString(i.Lg.Append(fmt.Sprintf("unable to create qrcode: %s", err), log.LEVEL_ERROR)), 0)
					}

					i.Self = &collection.ItemQR{
						QR: qr,
					}
				},
			})

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func to_image()
	/// @arg id
	/// @arg name
	/// @arg size - positive sets a fixed size, negative sets a scaled size
	/// @arg encoding
	/// @returns image id
	lib.CreateFunction(tab, "to_image",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "size"},
			{Type: lua.INT, Name: "encoding"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
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
			chLog.Parent(lg)
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

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func to_string()
	/// @arg id
	/// @arg inverse
	/// @returns string
	/// @blocking
	lib.CreateFunction(tab, "to_string",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.BOOL, Name: "inverse"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var str string

			<-r.QR.Schedule(args["id"].(int), &collection.Task[collection.ItemQR]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemQR]) {
					str = i.Self.QR.ToString(args["inverse"].(bool))
				},
			})

			state.Push(golua.LString(str))
			return 1
		})

	/// @func to_small_string()
	/// @arg id
	/// @arg inverse
	/// @returns string
	/// @blocking
	lib.CreateFunction(tab, "to_small_string",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.BOOL, Name: "inverse"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var str string

			<-r.QR.Schedule(args["id"].(int), &collection.Task[collection.ItemQR]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemQR]) {
					str = i.Self.QR.ToSmallString(args["inverse"].(bool))
				},
			})

			state.Push(golua.LString(str))
			return 1
		})

	/// @func color_set_foreground()
	/// @arg id
	/// @arg color {red, green, blue, alpha}
	lib.CreateFunction(tab, "color_set_foreground",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.TABLE, Name: "color", Table: &[]lua.Arg{
				{Type: lua.INT, Name: "red"},
				{Type: lua.INT, Name: "green"},
				{Type: lua.INT, Name: "blue"},
				{Type: lua.INT, Name: "alpha"},
			}},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
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
	lib.CreateFunction(tab, "color_set_background",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.TABLE, Name: "color", Table: &[]lua.Arg{
				{Type: lua.INT, Name: "red"},
				{Type: lua.INT, Name: "green"},
				{Type: lua.INT, Name: "blue"},
				{Type: lua.INT, Name: "alpha"},
			}},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
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
	lib.CreateFunction(tab, "color_foreground",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
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

			t := state.NewTable()
			state.SetField(t, "red", golua.LNumber(re))
			state.SetField(t, "green", golua.LNumber(gr))
			state.SetField(t, "blue", golua.LNumber(bl))
			state.SetField(t, "alpha", golua.LNumber(al))
			state.Push(t)
			return 1
		})

	/// @func color_background()
	/// @arg id
	/// @returns {red, green, blue, alpha}
	/// @blocking
	lib.CreateFunction(tab, "color_background",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
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

			t := state.NewTable()
			state.SetField(t, "red", golua.LNumber(re))
			state.SetField(t, "green", golua.LNumber(gr))
			state.SetField(t, "blue", golua.LNumber(bl))
			state.SetField(t, "alpha", golua.LNumber(al))
			state.Push(t)
			return 1
		})

	/// @func border()
	/// @arg id
	/// @returns bool
	/// @blocking
	lib.CreateFunction(tab, "border",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var active bool

			<-r.QR.Schedule(args["id"].(int), &collection.Task[collection.ItemQR]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemQR]) {
					active = !i.Self.QR.DisableBorder
				},
			})

			state.Push(golua.LBool(active))
			return 1
		})

	/// @func recovery_level()
	/// @arg id
	/// @returns int
	/// @blocking
	lib.CreateFunction(tab, "recovery_level",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var level qrcode.RecoveryLevel

			<-r.QR.Schedule(args["id"].(int), &collection.Task[collection.ItemQR]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemQR]) {
					level = i.Self.QR.Level
				},
			})

			state.Push(golua.LNumber(level))
			return 1
		})

	/// @func version()
	/// @arg id
	/// @returns int
	/// @blocking
	lib.CreateFunction(tab, "version",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var version int

			<-r.QR.Schedule(args["id"].(int), &collection.Task[collection.ItemQR]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemQR]) {
					version = i.Self.QR.VersionNumber
				},
			})

			state.Push(golua.LNumber(version))
			return 1
		})

	/// @func border_set()
	/// @arg id
	/// @arg? border
	lib.CreateFunction(tab, "border_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.BOOL, Name: "border", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
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
	lib.CreateFunction(tab, "content_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "content"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
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
	lib.CreateFunction(tab, "content",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var content string

			<-r.QR.Schedule(args["id"].(int), &collection.Task[collection.ItemQR]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemQR]) {
					content = i.Self.QR.Content
				},
			})

			state.Push(golua.LString(content))
			return 1
		})

	/// @constants Recovery Levels
	/// @const RECOVERY_LOW
	/// @const RECOVERY_MEDIUM
	/// @const RECOVERY_HIGH
	/// @const RECOVERY_HIGHEST
	tab.RawSetString("RECOVERY_LOW", golua.LNumber(qrcode.Low))
	tab.RawSetString("RECOVERY_MEDIUM", golua.LNumber(qrcode.Medium))
	tab.RawSetString("RECOVERY_HIGH", golua.LNumber(qrcode.High))
	tab.RawSetString("RECOVERY_HIGHEST", golua.LNumber(qrcode.Highest))
}

var recoveryLevelList = []qrcode.RecoveryLevel{
	qrcode.Low,
	qrcode.Medium,
	qrcode.High,
	qrcode.Highest,
}
