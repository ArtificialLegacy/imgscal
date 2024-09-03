package lib

import (
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	golua "github.com/yuin/gopher-lua"
)

const LIB_TEMPLATE = "template"

/// @lib Template
/// @import template
/// @desc
/// Library for formatting text using Go templates.

func RegisterTemplate(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_TEMPLATE, r, r.State, lg)

	/// @func parse(name, template, data) -> string
	/// @arg name {string}
	/// @arg template {string} - A string containing a Go template.
	/// @arg data {any} - The data to pass into the template.
	/// @returns {string}
	lib.CreateFunction(tab, "parse",
		[]lua.Arg{
			{Type: lua.STRING, Name: "name"},
			{Type: lua.STRING, Name: "template"},
			{Type: lua.ANY, Name: "data"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t, err := template.New(args["name"].(string)).Parse(args["template"].(string))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to parse template: %s", err), log.LEVEL_ERROR)), 0)
			}

			data := lua.GetValue(args["data"].(golua.LValue))
			str := &strings.Builder{}

			err = t.Execute(str, data)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed executing template: %s", err), log.LEVEL_ERROR)), 0)
			}

			state.Push(golua.LString(str.String()))
			return 1
		})

	/// @func parse_file(path, data) -> string
	/// @arg path {string} - Path to the file to parse as a Go template.
	/// @arg data {any} - The data to pass into the template.
	/// @returns {string}
	lib.CreateFunction(tab, "parse_file",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
			{Type: lua.ANY, Name: "data"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			pth := args["path"].(string)

			b, err := os.ReadFile(pth)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to read file: %s, %s", pth, err), log.LEVEL_ERROR)), 0)
			}

			t, err := template.New(strings.TrimSuffix(path.Base(pth), path.Ext(pth))).Parse(string(b))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to parse template: %s, %s", pth, err), log.LEVEL_ERROR)), 0)
			}

			data := lua.GetValue(args["data"].(golua.LValue))
			str := &strings.Builder{}

			err = t.Execute(str, data)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed executing template: %s", err), log.LEVEL_ERROR)), 0)
			}

			state.Push(golua.LString(str.String()))
			return 1
		})

	/// @func parse_to_file(name, template, outpath, data)
	/// @arg name {string}
	/// @arg template {string} - A string containing a Go template.
	/// @arg outpath {string} - Path to the file to output the result to.
	/// @arg data {any} - The data to pass into the template.
	lib.CreateFunction(tab, "parse_to_file",
		[]lua.Arg{
			{Type: lua.STRING, Name: "name"},
			{Type: lua.STRING, Name: "template"},
			{Type: lua.STRING, Name: "outpath"},
			{Type: lua.ANY, Name: "data"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t, err := template.New(args["name"].(string)).Parse(args["template"].(string))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to parse template: %s", err), log.LEVEL_ERROR)), 0)
			}

			outpath := args["outpath"].(string)
			f, err := os.OpenFile(outpath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0o666)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to open file for writing: %s, %s", outpath, err), log.LEVEL_ERROR)), 0)
			}
			defer f.Close()

			data := lua.GetValue(args["data"].(golua.LValue))
			err = t.Execute(f, data)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed executing template: %s", err), log.LEVEL_ERROR)), 0)
			}
			return 0
		})

	/// @func parse_file_to_file(path, outpath, data)
	/// @arg path {string} - Path to the file to parse as a Go template.
	/// @arg outpath {string} - Path to the file to output the result to.
	/// @arg data {any} - The data to pass into the template.
	lib.CreateFunction(tab, "parse_file_to_file",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
			{Type: lua.STRING, Name: "outpath"},
			{Type: lua.ANY, Name: "data"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			pth := args["path"].(string)

			b, err := os.ReadFile(pth)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to read file: %s, %s", pth, err), log.LEVEL_ERROR)), 0)
			}

			t, err := template.New(strings.TrimSuffix(path.Base(pth), path.Ext(pth))).Parse(string(b))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to parse template: %s, %s", pth, err), log.LEVEL_ERROR)), 0)
			}

			outpath := args["outpath"].(string)
			f, err := os.OpenFile(outpath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0o666)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to open file for writing: %s, %s", outpath, err), log.LEVEL_ERROR)), 0)
			}
			defer f.Close()

			data := lua.GetValue(args["data"].(golua.LValue))
			err = t.Execute(f, data)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed executing template: %s", err), log.LEVEL_ERROR)), 0)
			}
			return 0
		})
}
