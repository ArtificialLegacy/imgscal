package test

import (
	"testing"

	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	golua "github.com/yuin/gopher-lua"
)

func setupLib() *lua.Lib {
	state := golua.NewState()
	lg := log.NewLoggerEmpty()
	r := lua.NewRunner(state, &lg, false)
	lib, _ := lua.NewLib("testing", &r, state, &lg)

	return lib
}

func TestParseArgs_Int(t *testing.T) {
	lib := setupLib()

	lib.State.Push(golua.LNumber(1))
	lib.State.Push(golua.LNumber(2))

	argMap, _ := lib.ParseArgs(lib.State, "test_int", []lua.Arg{{Type: lua.INT, Name: "v1"}, {Type: lua.INT, Name: "v2"}}, 2, 0)

	if v, ok := argMap["v1"]; ok {
		if v.(int) != 1 {
			t.Errorf("got wrong int: wanted=%d, got=%d", 1, v)
		}
	} else {
		t.Error("failed to parse v1 argument")
	}

	if v, ok := argMap["v2"]; ok {
		if v.(int) != 2 {
			t.Errorf("got wrong int: wanted=%d, got=%d", 2, v)
		}
	} else {
		t.Error("failed to parse v2 argument")
	}
}

func TestParseArgs_Float(t *testing.T) {
	lib := setupLib()

	lib.State.Push(golua.LNumber(1.5))
	lib.State.Push(golua.LNumber(2.5))

	argMap, _ := lib.ParseArgs(lib.State, "test_float", []lua.Arg{{Type: lua.FLOAT, Name: "v1"}, {Type: lua.FLOAT, Name: "v2"}}, 2, 0)

	if v, ok := argMap["v1"]; ok {
		if v.(float64) != 1.5 {
			t.Errorf("got wrong float: wanted=%f, got=%f", 1.5, v)
		}
	} else {
		t.Error("failed to parse v1 argument")
	}

	if v, ok := argMap["v2"]; ok {
		if v.(float64) != 2.5 {
			t.Errorf("got wrong float: wanted=%f, got=%f", 2.5, v)
		}
	} else {
		t.Error("failed to parse v2 argument")
	}
}

func TestParseArgs_Bool(t *testing.T) {
	lib := setupLib()

	lib.State.Push(golua.LBool(true))
	lib.State.Push(golua.LBool(false))

	argMap, _ := lib.ParseArgs(lib.State, "test_bool", []lua.Arg{{Type: lua.BOOL, Name: "v1"}, {Type: lua.BOOL, Name: "v2"}}, 2, 0)

	if v, ok := argMap["v1"]; ok {
		if v.(bool) != true {
			t.Errorf("got wrong bool: wanted=%t, got=%t", true, v)
		}
	} else {
		t.Error("failed to parse v1 argument")
	}

	if v, ok := argMap["v2"]; ok {
		if v.(bool) != false {
			t.Errorf("got wrong bool: wanted=%t, got=%t", false, v)
		}
	} else {
		t.Error("failed to parse v2 argument")
	}
}

func TestParseArgs_String(t *testing.T) {
	lib := setupLib()

	lib.State.Push(golua.LString("A"))
	lib.State.Push(golua.LString("B"))

	argMap, _ := lib.ParseArgs(lib.State, "test_string", []lua.Arg{{Type: lua.STRING, Name: "v1"}, {Type: lua.STRING, Name: "v2"}}, 2, 0)

	if v, ok := argMap["v1"]; ok {
		if v.(string) != "A" {
			t.Errorf("got wrong string: wanted=%s, got=%s", "A", v)
		}
	} else {
		t.Error("failed to parse v1 argument")
	}

	if v, ok := argMap["v2"]; ok {
		if v.(string) != "B" {
			t.Errorf("got wrong string: wanted=%s, got=%s", "B", v)
		}
	} else {
		t.Error("failed to parse v2 argument")
	}
}

func TestParseArgs_Optional(t *testing.T) {
	lib := setupLib()

	lib.State.Push(golua.LString("A"))

	argMap, _ := lib.ParseArgs(lib.State, "test_optional", []lua.Arg{{Type: lua.STRING, Name: "v1"}, {Type: lua.STRING, Name: "v2", Optional: true}}, 1, 0)

	if v, ok := argMap["v1"]; ok {
		if v.(string) != "A" {
			t.Errorf("got wrong string: wanted=%s, got=%s", "A", v)
		}
	} else {
		t.Error("failed to parse v1 argument")
	}

	if v, ok := argMap["v2"]; ok {
		if v.(string) != "" {
			t.Errorf("got wrong string: wanted=%s, got=%s", "", v)
		}
	} else {
		t.Error("failed to parse v2 argument")
	}
}

func TestParseArgs_Table(t *testing.T) {
	lib := setupLib()

	tab := lib.State.NewTable()
	tab.RawSetString("v1", golua.LNumber(1))
	lib.State.Push(tab)

	argMap, _ := lib.ParseArgs(lib.State, "test_table", []lua.Arg{{Type: lua.TABLE, Name: "v1", Table: &[]lua.Arg{{Type: lua.INT, Name: "v1"}}}}, 1, 0)

	if v, ok := argMap["v1"]; ok {
		if v, ok := v.(map[string]any)["v1"]; ok {
			if v.(int) != 1 {
				t.Errorf("got wrong number: wanted=%d, got=%d", 1, v)
			}
		} else {
			t.Error("failed to parse v1 field")
		}
	} else {
		t.Error("failed to parse v1 argument")
	}
}

func TestParseArgs_Array(t *testing.T) {
	lib := setupLib()

	tab := lib.State.NewTable()
	tab.RawSetInt(1, golua.LNumber(1))
	tab.RawSetInt(2, golua.LNumber(2))
	lib.State.Push(tab)

	argMap, _ := lib.ParseArgs(lib.State, "test_array", []lua.Arg{lua.ArgArray("v1", lua.ArrayType{Type: lua.INT}, false)}, 1, 0)

	if v, ok := argMap["v1"]; ok {
		if v, ok := v.([]any); ok {
			if v[0].(int) != 1 {
				t.Errorf("got wrong number: wanted=%d, got=%d", 1, v[0])
			}
			if v[1].(int) != 2 {
				t.Errorf("got wrong number: wanted=%d, got=%d", 2, v[1])
			}
		} else {
			t.Error("failed to parse v1 field")
		}
	} else {
		t.Error("failed to parse v1 argument")
	}
}

func TestParseArgs_Variadic(t *testing.T) {
	lib := setupLib()

	lib.State.Push(golua.LNumber(1))
	lib.State.Push(golua.LNumber(2))

	argMap, _ := lib.ParseArgs(lib.State, "test_variadic", []lua.Arg{lua.ArgVariadic("v1", lua.ArrayType{Type: lua.INT}, false)}, 2, 0)

	if v, ok := argMap["v1"]; ok {
		if v, ok := v.([]any); ok {
			if v[0].(int) != 1 {
				t.Errorf("got wrong number: wanted=%d, got=%d", 1, v[0])
			}
			if v[1].(int) != 2 {
				t.Errorf("got wrong number: wanted=%d, got=%d", 2, v[1])
			}
		} else {
			t.Error("failed to parse v1 field")
		}
	} else {
		t.Error("failed to parse v1 argument")
	}
}

func TestMapSchema(t *testing.T) {
	schema := map[string]any{"v1": 1, "v2": 2, "v3": "A"}
	data := map[string]any{"v1": 2, "v2": 3, "v4": "B"}
	result := lua.MapSchema(schema, data)

	if v, ok := result["v1"]; ok {
		if v != 2 {
			t.Errorf("got wrong number: wanted=%d, got=%d", 2, v)
		}
	} else {
		t.Error("failed to map v1 field")
	}

	if v, ok := result["v2"]; ok {
		if v != 3 {
			t.Errorf("got wrong number: wanted=%d, got=%d", 3, v)
		}
	} else {
		t.Error("failed to map v2 field")
	}

	if v, ok := result["v3"]; ok {
		if v != "A" {
			t.Errorf("got wrong string: wanted=%s, got=%s", "A", v)
		}
	} else {
		t.Error("failed to map v3 field")
	}

	if v, ok := result["v4"]; ok {
		t.Errorf("v4 should not be in the result, but got=%s", v)
	}
}
