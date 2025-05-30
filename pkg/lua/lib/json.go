package lib

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	"github.com/tailscale/hujson"
	golua "github.com/yuin/gopher-lua"
)

const LIB_JSON = "json"

/// @lib JSON
/// @import json
/// @desc
/// Library for parsing and saving arbitrary json data.

func RegisterJSON(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_JSON, r, r.State, lg)

	/// @func parse(path) -> table<any>
	/// @arg path {string}
	/// @returns {table<any>} - Table representing the json file parsed.
	lib.CreateFunction(tab, "parse",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			file, err := os.Stat(args["path"].(string))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("invalid json path provided to io.load_image: %s", args["path"]), log.LEVEL_ERROR)), 0)
			}
			if file.IsDir() {
				state.Error(golua.LString(lg.Append("cannot parse a directory as an json", log.LEVEL_ERROR)), 0)
			}
			if path.Ext(file.Name()) != ".json" {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("file is not recognized as json: %s has extension: '%s' not '.json'", file.Name(), path.Ext(file.Name())), log.LEVEL_ERROR)), 0)
			}

			fb, err := os.ReadFile(args["path"].(string))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("cannot open file %s: %s", args["path"], err.Error()), log.LEVEL_ERROR)), 0)
			}

			var data map[string]any
			err = json.Unmarshal(fb, &data)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to unmarshal json: %s", err.Error()), log.LEVEL_ERROR)), 0)
			}

			state.Push(lua.CreateValue(data, state))
			return 1
		})

	/// @func parse_schema(path, schema) -> table<any>
	/// @arg path {string}
	/// @arg schema {table<any>}
	/// @returns {table<any>} - Table representing the json file parsed.
	lib.CreateFunction(tab, "parse_schema",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
			{Type: lua.RAW_TABLE, Name: "schema"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			file, err := os.Stat(args["path"].(string))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("invalid json path provided to io.load_image: %s", args["path"]), log.LEVEL_ERROR)), 0)
			}
			if file.IsDir() {
				state.Error(golua.LString(lg.Append("cannot parse a directory as an json", log.LEVEL_ERROR)), 0)
			}
			if path.Ext(file.Name()) != ".json" {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("file is not recognized as json: %s has extension: '%s' not '.json'", file.Name(), path.Ext(file.Name())), log.LEVEL_ERROR)), 0)
			}

			fb, err := os.ReadFile(args["path"].(string))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("cannot open file %s: %s", args["path"], err.Error()), log.LEVEL_ERROR)), 0)
			}

			var data map[string]any
			err = json.Unmarshal(fb, &data)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to unmarshal json: %s", err.Error()), log.LEVEL_ERROR)), 0)
			}

			schema := lua.GetValue(args["schema"].(*golua.LTable)).(map[string]any)
			result := lua.MapSchema(schema, data)

			state.Push(lua.CreateValue(result, state))
			return 1
		})

	/// @func parse_string(str) -> table<any>
	/// @arg str {string}
	/// @returns {table<any>} - Table representing the json string parsed.
	lib.CreateFunction(tab, "parse_string",
		[]lua.Arg{
			{Type: lua.STRING, Name: "str"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			fb := []byte(args["str"].(string))

			var data map[string]any
			err := json.Unmarshal(fb, &data)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to unmarshal json: %s", err.Error()), log.LEVEL_ERROR)), 0)
			}

			state.Push(lua.CreateValue(data, state))
			return 1
		})

	/// @func parse_string_schema(str, schema) -> table<any>
	/// @arg str {string}
	/// @arg schema {table<any>}
	/// @returns {table<any>} - Table representing the json string parsed.
	lib.CreateFunction(tab, "parse_string_schema",
		[]lua.Arg{
			{Type: lua.STRING, Name: "str"},
			{Type: lua.RAW_TABLE, Name: "schema"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			fb := []byte(args["str"].(string))

			var data map[string]any
			err := json.Unmarshal(fb, &data)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to unmarshal json: %s", err.Error()), log.LEVEL_ERROR)), 0)
			}

			schema := lua.GetValue(args["schema"].(*golua.LTable)).(map[string]any)
			result := lua.MapSchema(schema, data)

			state.Push(lua.CreateValue(result, state))
			return 1
		})

	/// @func save(value, path, compact?)
	/// @arg value {table<any>} - Table to convert to json.
	/// @arg path {string}
	/// @arg? compact {bool} - Defaults to false, use to remove indent and new lines.
	lib.CreateFunction(tab, "save",
		[]lua.Arg{
			{Type: lua.ANY, Name: "value"},
			{Type: lua.STRING, Name: "path"},
			{Type: lua.BOOL, Name: "compact", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			file, err := os.OpenFile(args["path"].(string), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0o666)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("cannot open file: %s", args["path"].(string)), log.LEVEL_ERROR)), 0)
			}
			defer file.Close()

			data := lua.GetValue(args["value"].(golua.LValue))

			var b []byte

			if args["compact"].(bool) {
				b, err = json.Marshal(data)
			} else {
				b, err = json.MarshalIndent(data, "", "    ")
			}

			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to marshal json: %s", err.Error()), log.LEVEL_ERROR)), 0)
			}

			_, err = file.Write(b)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to write json to file: %s", err.Error()), log.LEVEL_ERROR)), 0)
			}

			return 0
		})

	/// @func string(value, compact?) -> string
	/// @arg value {table<any>} - Table to convert to json.
	/// @arg? compact {bool} - Defaults to false, use to remove indent and new lines.
	/// @returns {string}
	lib.CreateFunction(tab, "string",
		[]lua.Arg{
			{Type: lua.ANY, Name: "value"},
			{Type: lua.BOOL, Name: "compact", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			data := lua.GetValue(args["value"].(golua.LValue))

			var b []byte
			var err error

			if args["compact"].(bool) {
				b, err = json.Marshal(data)
			} else {
				b, err = json.MarshalIndent(data, "", "    ")
			}
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to marshal json: %s", err.Error()), log.LEVEL_ERROR)), 0)
			}

			state.Push(golua.LString(b))
			return 1
		})

	/// @func standardize(str) -> string
	/// @arg str {string} - JSON data as a string.
	/// @returns {string} - The same JSON but with non-standard features removed.
	lib.CreateFunction(tab, "standardize",
		[]lua.Arg{
			{Type: lua.STRING, Name: "str"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			value := args["str"].(string)
			vout, err := hujson.Standardize([]byte(value))
			if err != nil {
				lua.Error(state, lg.Appendf("failed to standardize JSON: %s", log.LEVEL_ERROR, err))
			}

			state.Push(golua.LString(vout))
			return 1
		})
}
