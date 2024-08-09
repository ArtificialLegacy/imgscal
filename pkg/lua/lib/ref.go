package lib

import (
	"fmt"
	"image/color"
	"time"

	g "github.com/AllenDang/giu"
	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	imageutil "github.com/ArtificialLegacy/imgscal/pkg/image_util"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	golua "github.com/yuin/gopher-lua"
)

const LIB_REF = "ref"

func RegisterRef(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_REF, r, r.State, lg)

	/// @func new()
	/// @arg value
	/// @arg? type
	/// @returns id
	/// @desc
	/// References are used when go and lua need to share a reference to the same value.
	/// The primitive type versions must be used when that value must be a Go value.
	/// Also note that when refs are used for indexes in Go that they must start at 0 not 1.
	lib.CreateFunction(tab, "new",
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

			default:
				state.Error(golua.LString(lg.Append(fmt.Sprintf("unknown reftype for new: %d", args["type"]), log.LEVEL_ERROR)), 0)
			}

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func get()
	/// @arg id
	/// @returns value
	/// @desc
	/// Note: this is a copy of the value being referenced, to mutate the ref use ref.set().
	lib.CreateFunction(tab, "get",
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
			case *g.FontInfo:
				state.Push(golua.LString(item.String()))

			default:
				state.Error(golua.LString(lg.Append(fmt.Sprintf("unknown reftype for get: %T", item), log.LEVEL_ERROR)), 0)
			}

			return 1
		})

	/// @func set()
	/// @arg id
	/// @arg value
	lib.CreateFunction(tab, "set",
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
				state.Error(golua.LString(lg.Append(fmt.Sprintf("unknown reftype for set: %T", i.Value), log.LEVEL_ERROR)), 0)
			}

			return 0
		})

	/// @func del()
	/// @arg id
	lib.CreateFunction(tab, "del",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CR_REF.Clean(args["id"].(int))
			return 0
		})

	/// @func del_many()
	/// @arg []ids
	lib.CreateFunction(tab, "del_many",
		[]lua.Arg{
			lua.ArgArray("ids", lua.ArrayType{Type: lua.INT}, false),
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			ids := args["ids"].(map[string]any)
			for _, id := range ids {
				r.CR_REF.Clean(id.(int))
			}

			return 0
		})

	/// @constants Ref Types
	/// @const LUA
	/// @const BOOL
	/// @const INT
	/// @const INT32
	/// @const FLOAT
	/// @const FLOAT32
	/// @const STRING
	/// @const RGBA
	/// @const TIME - timestamp in ms
	/// @const FONT - giu.FontInfo - cannot be created or set, getting will return .String()
	r.State.SetField(tab, "LUA", golua.LNumber(REFTYPE_LUA))
	r.State.SetField(tab, "BOOL", golua.LNumber(REFTYPE_BOOL))
	r.State.SetField(tab, "INT", golua.LNumber(REFTYPE_INT))
	r.State.SetField(tab, "INT32", golua.LNumber(REFTYPE_INT32))
	r.State.SetField(tab, "FLOAT", golua.LNumber(REFTYPE_FLOAT))
	r.State.SetField(tab, "FLOAT32", golua.LNumber(REFTYPE_FLOAT32))
	r.State.SetField(tab, "STRING", golua.LNumber(REFTYPE_STRING))
	r.State.SetField(tab, "RGBA", golua.LNumber(REFTYPE_RGBA))
	r.State.SetField(tab, "TIME", golua.LNumber(REFTYPE_TIME))
	r.State.SetField(tab, "FONT", golua.LNumber(REFTYPE_FONT))
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
	REFTYPE_FONT
)
