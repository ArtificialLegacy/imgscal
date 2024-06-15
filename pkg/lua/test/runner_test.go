package test

import (
	"testing"

	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	golua "github.com/Shopify/go-lua"
)

func setupLib() *lua.Lib {
	state := golua.NewState()
	lg := log.NewLoggerEmpty()
	lib := lua.NewLib("testing", state, &lg)

	return lib
}

func TestParseArgs_Int(t *testing.T) {
	lib := setupLib()

	lib.State.PushInteger(1)
	lib.State.PushInteger(2)

	argMap := lib.ParseArgs("test_int", []lua.Arg{{Type: lua.INT, Name: "int1"}, {Type: lua.INT, Name: "int2"}})

	if v, ok := argMap["int1"]; ok {
		if v.(int) != 1 {
			t.Errorf("got wrong int: wanted=%d, got=%d", 1, v)
		}
	} else {
		t.Error("failed to parse int1 argument")
	}

	if v, ok := argMap["int2"]; ok {
		if v.(int) != 2 {
			t.Errorf("got wrong int: wanted=%d, got=%d", 2, v)
		}
	} else {
		t.Error("failed to parse int2 argument")
	}
}

func TestParseArgs_Float(t *testing.T) {
	lib := setupLib()

	lib.State.PushNumber(1.5)
	lib.State.PushNumber(2.5)

	argMap := lib.ParseArgs("test_float", []lua.Arg{{Type: lua.FLOAT, Name: "float1"}, {Type: lua.FLOAT, Name: "float2"}})

	if v, ok := argMap["float1"]; ok {
		if v.(float64) != 1.5 {
			t.Errorf("got wrong float: wanted=%f, got=%f", 1.5, v)
		}
	} else {
		t.Error("failed to parse float1 argument")
	}

	if v, ok := argMap["float2"]; ok {
		if v.(float64) != 2.5 {
			t.Errorf("got wrong float: wanted=%f, got=%f", 2.5, v)
		}
	} else {
		t.Error("failed to parse float2 argument")
	}
}

func TestParseArgs_Bool(t *testing.T) {
	lib := setupLib()

	lib.State.PushBoolean(true)
	lib.State.PushBoolean(false)

	argMap := lib.ParseArgs("test_bool", []lua.Arg{{Type: lua.BOOL, Name: "bool1"}, {Type: lua.BOOL, Name: "bool2"}})

	if v, ok := argMap["bool1"]; ok {
		if v.(bool) != true {
			t.Errorf("got wrong bool: wanted=%t, got=%t", true, v)
		}
	} else {
		t.Error("failed to parse float1 argument")
	}

	if v, ok := argMap["bool2"]; ok {
		if v.(bool) != false {
			t.Errorf("got wrong bool: wanted=%t, got=%t", true, v)
		}
	} else {
		t.Error("failed to parse bool2 argument")
	}
}

func TestParseArgs_String(t *testing.T) {
	lib := setupLib()

	lib.State.PushString("test1")
	lib.State.PushString("test2")

	argMap := lib.ParseArgs("test_string", []lua.Arg{{Type: lua.STRING, Name: "string1"}, {Type: lua.STRING, Name: "string2"}})

	if v, ok := argMap["string1"]; ok {
		if v.(string) != "test1" {
			t.Errorf("got wrong string: wanted=%s, got=%s", "test1", v)
		}
	} else {
		t.Error("failed to parse string1 argument")
	}

	if v, ok := argMap["string2"]; ok {
		if v.(string) != "test2" {
			t.Errorf("got wrong string: wanted=%s, got=%s", "test2", v)
		}
	} else {
		t.Error("failed to parse string2 argument")
	}
}

func TestParseArgs_Any(t *testing.T) {
	lib := setupLib()

	lib.State.PushString("test1")
	lib.State.PushString("test2")

	argMap := lib.ParseArgs("test_any", []lua.Arg{{Type: lua.ANY, Name: "string1"}, {Type: lua.ANY, Name: "string2"}})

	if v, ok := argMap["string1"]; ok {
		if v.(string) != "test1" {
			t.Errorf("got wrong string: wanted=%s, got=%s", "test1", v)
		}
	} else {
		t.Error("failed to parse string1 argument")
	}

	if v, ok := argMap["string2"]; ok {
		if v.(string) != "test2" {
			t.Errorf("got wrong string: wanted=%s, got=%s", "test2", v)
		}
	} else {
		t.Error("failed to parse string2 argument")
	}
}

func TestParseArgs_Table(t *testing.T) {
	lib := setupLib()

	lib.State.NewTable()
	lib.State.PushString("test1")
	lib.State.SetField(-2, "test")

	argMap := lib.ParseArgs("test_table", []lua.Arg{
		{Type: lua.TABLE, Name: "table1", Table: &[]lua.Arg{
			{Type: lua.STRING, Name: "test"},
		}},
	})

	if v, ok := argMap["table1"]; ok {
		if v, ok := v.(map[string]any)["test"]; ok {
			if v.(string) != "test1" {
				t.Errorf("got wrong string: wanted=%s, got=%s", "test1", v)
			}
		} else {
			t.Error("failed to parse table field test")
		}
	} else {
		t.Error("failed to parse table1 argument")
	}
}

func TestParseArgs_TableNested(t *testing.T) {
	lib := setupLib()

	lib.State.NewTable()
	lib.State.NewTable()
	lib.State.PushString("test1")
	lib.State.SetField(-2, "test")
	lib.State.SetField(-2, "table2")

	argMap := lib.ParseArgs("test_table", []lua.Arg{
		{Type: lua.TABLE, Name: "table1", Table: &[]lua.Arg{
			{Type: lua.TABLE, Name: "table2", Table: &[]lua.Arg{
				{Type: lua.STRING, Name: "test"},
			}},
		}},
	})

	if v, ok := argMap["table1"]; ok {
		if v, ok := v.(map[string]any)["table2"]; ok {
			if v, ok := v.(map[string]any)["test"]; ok {
				if v.(string) != "test1" {
					t.Errorf("got wrong string: wanted=%s, got=%s", "test1", v)
				}
			} else {
				t.Error("failed to parse table2 field test")
			}
		} else {
			t.Error("failed to parse table1 field table2")
		}
	} else {
		t.Error("failed to parse table1 argument")
	}
}

func TestParseArgs_Array(t *testing.T) {
	lib := setupLib()

	lib.State.NewTable()
	lib.State.PushInteger(1)
	lib.State.PushString("test1")
	lib.State.SetTable(-3)

	argMap := lib.ParseArgs("test_array", []lua.Arg{
		lua.ArgArray("array1", lua.ArrayType{Type: lua.STRING}, false),
	})

	if v, ok := argMap["array1"]; ok {
		if v, ok := v.(map[string]any)["1"]; ok {
			if v.(string) != "test1" {
				t.Errorf("got wrong string: wanted=%s, got=%s", "test1", v)
			}
		} else {
			t.Error("failed to parse array field ")
		}
	} else {
		t.Error("failed to parse array1 argument")
	}
}

func TestParseArgs_Optional(t *testing.T) {
	lib := setupLib()

	argMap := lib.ParseArgs("test_optional", []lua.Arg{
		{Type: lua.STRING, Name: "test", Optional: true},
	})

	if argMap["test"] != "" {
		t.Errorf("got incorrect value, expected=0, got=%v", argMap["test"])
	}
}

func TestParseArgs_OptionalTableField(t *testing.T) {
	lib := setupLib()

	lib.State.NewTable()

	argMap := lib.ParseArgs("test_optional", []lua.Arg{
		{Type: lua.TABLE, Name: "test_table", Table: &[]lua.Arg{
			{Type: lua.STRING, Name: "test", Optional: true},
		}},
	})

	if argMap["test_table"].(map[string]any)["test"] != "" {
		t.Errorf("got incorrect value, expected='', got=%v", argMap["test"])
	}
}
