package lib

import (
	"fmt"
	"image/color"
	"time"

	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	imageutil "github.com/ArtificialLegacy/imgscal/pkg/image_util"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	golua "github.com/yuin/gopher-lua"
)

const LIB_STD = "std"

func RegisterStd(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_STD, r, r.State, lg)

	/// @func log()
	/// @arg msg - the message to display in the log
	lib.CreateFunction(tab, "log",
		[]lua.Arg{
			{Type: lua.ANY, Name: "msg"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			lg.Append(fmt.Sprintf("lua log: %s", args["msg"]), log.LEVEL_INFO)
			return 0
		})

	/// @func warn()
	/// @arg msg - the message to display as a warning in the log
	lib.CreateFunction(tab, "warn",
		[]lua.Arg{
			{Type: lua.STRING, Name: "msg"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			lg.Append(fmt.Sprintf("lua warn: %s", args["msg"]), log.LEVEL_WARN)
			return 0
		})

	/// @func panic()
	/// @arg msg - the message to display in the error
	lib.CreateFunction(tab, "panic",
		[]lua.Arg{
			{Type: lua.STRING, Name: "msg"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			state.Error(golua.LString(lg.Append(fmt.Sprintf("lua panic: %s", args["msg"]), log.LEVEL_ERROR)), 0)
			return 0
		})

	/// @func ref()
	/// @arg value
	/// @arg? type
	/// @returns id
	/// @desc
	/// References are used when go and lua need to share a reference to the same value.
	/// The primitive type versions must be used when that value must be a Go value.
	/// Also note that when refs are used for indexes in Go that they must start at 0 not 1.
	lib.CreateFunction(tab, "ref",
		[]lua.Arg{
			{Type: lua.ANY, Name: "value"},
			{Type: lua.INT, Name: "type", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			v := args["value"].(golua.LValue)
			var id int

			switch args["type"].(int) {
			case REFTYPE_LUA:
				id = r.CR_REF.Add(&collection.RefItem[any]{Value: &v})
			case REFTYPE_BOOL:
				v := bool(v.(golua.LBool))
				id = r.CR_REF.Add(&collection.RefItem[any]{Value: &v})
			case REFTYPE_INT:
				v := int(v.(golua.LNumber))
				id = r.CR_REF.Add(&collection.RefItem[any]{Value: &v})
			case REFTYPE_INT32:
				v := int32(v.(golua.LNumber))
				id = r.CR_REF.Add(&collection.RefItem[any]{Value: &v})
			case REFTYPE_FLOAT:
				v := float64(v.(golua.LNumber))
				id = r.CR_REF.Add(&collection.RefItem[any]{Value: &v})
			case REFTYPE_FLOAT32:
				v := float32(v.(golua.LNumber))
				id = r.CR_REF.Add(&collection.RefItem[any]{Value: &v})
			case REFTYPE_STRING:
				v := string(v.(golua.LString))
				id = r.CR_REF.Add(&collection.RefItem[any]{Value: &v})
			case REFTYPE_RGBA:
				v := v.(*golua.LTable)
				c := imageutil.TableToRGBA(state, v)

				id = r.CR_REF.Add(&collection.RefItem[any]{Value: c})
			case REFTYPE_TIME:
				v := v.(golua.LNumber)
				t := time.UnixMilli(int64(v))
				id = r.CR_REF.Add(&collection.RefItem[any]{Value: &t})
			}

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func ref_get()
	/// @arg id
	/// @returns value
	/// @desc
	/// Note: this is a copy of the value being referenced, to mutate the ref use ref_set().
	lib.CreateFunction(tab, "ref_get",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			i, err := r.CR_REF.Item(args["id"].(int))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("unable to find ref: %s", err), log.LEVEL_ERROR)), 0)
			}

			switch item := i.Value.(type) {
			case *golua.LValue:
				state.Push(*item)
			case *bool:
				state.Push(golua.LBool(*item))
			case *int:
				state.Push(golua.LNumber(*item))
			case *int32:
				state.Push(golua.LNumber(*item))
			case *float64:
				state.Push(golua.LNumber(*item))
			case *float32:
				state.Push(golua.LNumber(*item))
			case *string:
				state.Push(golua.LString(*item))
			case *color.RGBA:
				state.Push(imageutil.RGBAToTable(state, item))
			case *time.Time:
				state.Push(golua.LNumber(item.UnixMilli()))

			default:
				state.Error(golua.LString(lg.Append(fmt.Sprintf("unknown ref type: %T", item), log.LEVEL_ERROR)), 0)
			}

			return 1
		})

	/// @func ref_set()
	/// @arg id
	/// @arg value
	lib.CreateFunction(tab, "ref_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.ANY, Name: "value"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			i, err := r.CR_REF.Item(args["id"].(int))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("unable to find ref: %s", err), log.LEVEL_ERROR)), 0)
			}

			v := args["value"].(golua.LValue)

			switch i.Value.(type) {
			case *golua.LValue:
				i.Value = &v
			case *bool:
				v := bool(v.(golua.LBool))
				i.Value = &v
			case *int:
				v := int(v.(golua.LNumber))
				i.Value = &v
			case *int32:
				v := int32(v.(golua.LNumber))
				i.Value = &v
			case *float64:
				v := float64(v.(golua.LNumber))
				i.Value = &v
			case *float32:
				v := float32(v.(golua.LNumber))
				i.Value = &v
			case *string:
				v := string(v.(golua.LString))
				i.Value = &v
			case *color.RGBA:
				i.Value = imageutil.TableToRGBA(state, v.(*golua.LTable))
			case *time.Time:
				v := time.UnixMilli(int64(v.(golua.LNumber)))
				i.Value = &v

			default:
				state.Error(golua.LString(lg.Append(fmt.Sprintf("unknown ref type: %T", i.Value), log.LEVEL_ERROR)), 0)
			}

			return 0
		})

	/// @constants Ref Types
	/// @const REFTYPE_LUA
	/// @const REFTYPE_BOOL
	/// @const REFTYPE_INT
	/// @const REFTYPE_INT32
	/// @const REFTYPE_FLOAT
	/// @const REFTYPE_FLOAT32
	/// @const REFTYPE_STRING
	/// @const REFTYPE_RGBA
	/// @const REFTYPE_TIME - timestamp in ms
	r.State.SetField(tab, "REFTYPE_LUA", golua.LNumber(REFTYPE_LUA))
	r.State.SetField(tab, "REFTYPE_BOOL", golua.LNumber(REFTYPE_BOOL))
	r.State.SetField(tab, "REFTYPE_INT", golua.LNumber(REFTYPE_INT))
	r.State.SetField(tab, "REFTYPE_INT32", golua.LNumber(REFTYPE_INT32))
	r.State.SetField(tab, "REFTYPE_FLOAT", golua.LNumber(REFTYPE_FLOAT))
	r.State.SetField(tab, "REFTYPE_FLOAT32", golua.LNumber(REFTYPE_FLOAT32))
	r.State.SetField(tab, "REFTYPE_STRING", golua.LNumber(REFTYPE_STRING))
	r.State.SetField(tab, "REFTYPE_RGBA", golua.LNumber(REFTYPE_RGBA))
	r.State.SetField(tab, "REFTYPE_TIME", golua.LNumber(REFTYPE_TIME))
}

const (
	REFTYPE_LUA int = iota
	REFTYPE_BOOL
	REFTYPE_INT
	REFTYPE_INT32
	REFTYPE_FLOAT
	REFTYPE_FLOAT32
	REFTYPE_STRING
	REFTYPE_RGBA
	REFTYPE_TIME
)
