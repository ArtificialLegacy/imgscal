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

/// @lib References
/// @import ref
/// @desc
/// Library for creating and referencing pointers.
/// @section
/// This is used when both go and lua need to reference the same mutable data.

func RegisterRef(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_REF, r, r.State, lg)

	/// @func new(value, type?) -> int<collection.CRATE_REF>
	/// @arg value {any}
	/// @arg? type {int<ref.REFType>} - Must be compatible with the above value.
	/// @returns {int<collection.CRATE_REF>}
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
				c := imageutil.ColorTableToRGBAColor(v)

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

	/// @func new_slice(ln, type?) -> int<collection.CRATE_REF>
	/// @arg ln {int} - A fixed length for the slice.
	/// @arg? type {int<ref.REFType>} - Must be compatible with the values set into the slice.
	/// @returns {int<collection.CRATE_REF>}
	/// @desc
	/// References are used when go and lua need to share a reference to the same value.
	/// The primitive type versions must be used when that value must be a Go value.
	/// Also note that when refs are used for indexes in Go that they must start at 0 not 1.
	/// This also includes the slice created here, it's index starts at 0.
	lib.CreateFunction(tab, "new_slice",
		[]lua.Arg{
			{Type: lua.INT, Name: "ln"},
			{Type: lua.INT, Name: "type", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var id int

			ln := args["len"].(int)

			switch args["type"].(int) {
			case REFTYPE_LUA:
				val := make([]golua.LValue, ln)
				id = r.CR_REF.Add(&collection.RefItem[any]{Value: &val})
			case REFTYPE_BOOL:
				val := make([]bool, ln)
				id = r.CR_REF.Add(&collection.RefItem[any]{Value: &val})
			case REFTYPE_INT:
				val := make([]int, ln)
				id = r.CR_REF.Add(&collection.RefItem[any]{Value: &val})
			case REFTYPE_INT32:
				val := make([]int32, ln)
				r.CR_REF.Add(&collection.RefItem[any]{Value: &val})
			case REFTYPE_FLOAT:
				val := make([]float64, ln)
				id = r.CR_REF.Add(&collection.RefItem[any]{Value: &val})
			case REFTYPE_FLOAT32:
				val := make([]float32, ln)
				id = r.CR_REF.Add(&collection.RefItem[any]{Value: &val})
			case REFTYPE_STRING:
				val := make([]string, ln)
				id = r.CR_REF.Add(&collection.RefItem[any]{Value: &val})
			case REFTYPE_RGBA:
				val := make([]*color.RGBA, ln)
				id = r.CR_REF.Add(&collection.RefItem[any]{Value: &val})
			case REFTYPE_TIME:
				val := make([]*time.Time, ln)
				id = r.CR_REF.Add(&collection.RefItem[any]{Value: &val})

			default:
				state.Error(golua.LString(lg.Append(fmt.Sprintf("unknown reftype for new_slice: %d", args["type"]), log.LEVEL_ERROR)), 0)
			}

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func get(id) -> any
	/// @arg id {int<collection.CRATE_REF>}
	/// @returns {any} - Will be a set type deteremined by the REFType.
	/// @desc
	/// Note: this is a copy of the value being referenced, to mutate the ref use 'ref.set()'.
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
				state.Push(imageutil.RGBAColorToColorTable(state, item))
			case *time.Time:
				state.Push(golua.LNumber(item.UnixMilli()))
			case *g.FontInfo:
				state.Push(golua.LString(item.String()))

			default:
				state.Error(golua.LString(lg.Append(fmt.Sprintf("unknown reftype for get: %T", item), log.LEVEL_ERROR)), 0)
			}

			return 1
		})

	/// @func get_slice(id, index) -> any
	/// @arg id {int<collection.CRATE_REF>}
	/// @arg index {int} - Must be within the range of 0 to 1 less than the set length.
	/// @returns {any} - Will be a set type determined by the REFType.
	/// @desc
	/// Note: this is a copy of the value being referenced, to mutate the ref use 'ref.set_slice()'.
	lib.CreateFunction(tab, "get_slice",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "index"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			i, err := r.CR_REF.Item(args["id"].(int))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("unable to find ref: %s", err), log.LEVEL_ERROR)), 0)
			}

			index := args["index"].(int)

			switch item := i.Value.(type) {
			case *[]golua.LValue:
				state.Push((*item)[index])
			case *[]bool:
				state.Push(golua.LBool((*item)[index]))
			case *[]int:
				state.Push(golua.LNumber((*item)[index]))
			case *[]int32:
				state.Push(golua.LNumber((*item)[index]))
			case *[]float64:
				state.Push(golua.LNumber((*item)[index]))
			case *[]float32:
				state.Push(golua.LNumber((*item)[index]))
			case *[]string:
				state.Push(golua.LString((*item)[index]))
			case *[]*color.RGBA:
				state.Push(imageutil.RGBAColorToColorTable(state, (*item)[index]))
			case *[]*time.Time:
				state.Push(golua.LNumber((*item)[index].UnixMilli()))
			case *[]*g.FontInfo:
				state.Push(golua.LString((*item)[index].String()))

			default:
				state.Error(golua.LString(lg.Append(fmt.Sprintf("unknown reftype for get_slice: %T", item), log.LEVEL_ERROR)), 0)
			}

			return 1
		})

	/// @func len_slice(id) -> int
	/// @arg id {int<collection.CRATE_REF>}
	/// @returns {int} - The length of the slice in the ref.
	lib.CreateFunction(tab, "len_slice",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			i, err := r.CR_REF.Item(args["id"].(int))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("unable to find ref: %s", err), log.LEVEL_ERROR)), 0)
			}

			switch item := i.Value.(type) {
			case *[]golua.LValue:
				state.Push(golua.LNumber(len(*item)))
			case *[]bool:
				state.Push(golua.LNumber(len(*item)))
			case *[]int:
				state.Push(golua.LNumber(len(*item)))
			case *[]int32:
				state.Push(golua.LNumber(len(*item)))
			case *[]float64:
				state.Push(golua.LNumber(len(*item)))
			case *[]float32:
				state.Push(golua.LNumber(len(*item)))
			case *[]string:
				state.Push(golua.LNumber(len(*item)))
			case *[]*color.RGBA:
				state.Push(golua.LNumber(len(*item)))
			case *[]*time.Time:
				state.Push(golua.LNumber(len(*item)))
			case *[]*g.FontInfo:
				state.Push(golua.LNumber(len(*item)))

			default:
				state.Error(golua.LString(lg.Append(fmt.Sprintf("unknown reftype for len_slice: %T", item), log.LEVEL_ERROR)), 0)
			}

			return 1
		})

	/// @func set(id, value)
	/// @arg id {int<collection.CRATE_REF>}
	/// @arg value {any} - Must be compatible with the set REFType.
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
				i.Value = imageutil.ColorTableToRGBAColor(v.(*golua.LTable))
			case *time.Time:
				v := time.UnixMilli(int64(v.(golua.LNumber)))
				i.Value = &v

			default:
				state.Error(golua.LString(lg.Append(fmt.Sprintf("unknown reftype for set: %T", i.Value), log.LEVEL_ERROR)), 0)
			}

			return 0
		})

	/// @func set_slice(id, index, value)
	/// @arg id {int<collection.CRATE_REF>}
	/// @arg index {int} - Must be in range 0 to 1 less than the set length.
	/// @arg value {any} - Must be compatible with the set REFType.
	lib.CreateFunction(tab, "set_slice",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "index"},
			{Type: lua.ANY, Name: "value"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			i, err := r.CR_REF.Item(args["id"].(int))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("unable to find ref: %s", err), log.LEVEL_ERROR)), 0)
			}

			index := args["index"].(int)

			val := args["value"].(golua.LValue)

			switch v := i.Value.(type) {
			case *[]golua.LValue:
				(*v)[index] = val
			case *[]bool:
				(*v)[index] = bool(val.(golua.LBool))
			case *[]int:
				(*v)[index] = int(val.(golua.LNumber))
			case *[]int32:
				(*v)[index] = int32(val.(golua.LNumber))
			case *[]float64:
				(*v)[index] = float64(val.(golua.LNumber))
			case *[]float32:
				(*v)[index] = float32(val.(golua.LNumber))
			case *[]string:
				(*v)[index] = string(val.(golua.LString))
			case *[]*color.RGBA:
				(*v)[index] = imageutil.ColorTableToRGBAColor(val.(*golua.LTable))
			case *[]*time.Time:
				t := time.UnixMilli(int64(val.(golua.LNumber)))
				(*v)[index] = &t

			default:
				state.Error(golua.LString(lg.Append(fmt.Sprintf("unknown reftype for set_slice: %T", i.Value), log.LEVEL_ERROR)), 0)
			}

			return 0
		})

	/// @func del(id)
	/// @arg id {int<collection.CRATE_REF>}
	lib.CreateFunction(tab, "del",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CR_REF.Clean(args["id"].(int))
			return 0
		})

	/// @func del_many(ids)
	/// @arg ids {[]int<collection.CRATE_REF>}
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
	/// @const RGBA - Use in lua as struct<image.Color>.
	/// @const TIME - Use in lua as a number representing the time in ms.
	/// @const FONT - Internal 'giu.FontInfo', cannot be created or set, getting will return the result of 'font.String()'.
	tab.RawSetString("LUA", golua.LNumber(REFTYPE_LUA))
	tab.RawSetString("BOOL", golua.LNumber(REFTYPE_BOOL))
	tab.RawSetString("INT", golua.LNumber(REFTYPE_INT))
	tab.RawSetString("INT32", golua.LNumber(REFTYPE_INT32))
	tab.RawSetString("FLOAT", golua.LNumber(REFTYPE_FLOAT))
	tab.RawSetString("FLOAT32", golua.LNumber(REFTYPE_FLOAT32))
	tab.RawSetString("STRING", golua.LNumber(REFTYPE_STRING))
	tab.RawSetString("RGBA", golua.LNumber(REFTYPE_RGBA))
	tab.RawSetString("TIME", golua.LNumber(REFTYPE_TIME))
	tab.RawSetString("FONT", golua.LNumber(REFTYPE_FONT))
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
