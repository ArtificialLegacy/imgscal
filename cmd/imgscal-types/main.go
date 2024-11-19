package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/ArtificialLegacy/imgscal/pkg/doc"
	"github.com/ArtificialLegacy/imgscal/pkg/types"
)

const (
	LIB_DIR   = "./pkg/lua/lib"
	TYPES_DIR = "./types"
)

func main() {
	_, err := os.Stat(TYPES_DIR)
	if err != nil {
		os.Mkdir(TYPES_DIR, 0o777)
	}

	fs, err := os.ReadDir(LIB_DIR)
	if err != nil {
		panic(fmt.Sprintf("cannot open lua lib dir with err: %s", err))
	}

	docs := []doc.Lib{}

	for _, f := range fs {
		bs, err := os.ReadFile(path.Join(LIB_DIR, f.Name()))
		if err != nil {
			panic(fmt.Sprintf("failed file read on file %s", f.Name()))
		}

		docs = append(docs, doc.Parse(f.Name(), bs))
	}

	for _, lib := range docs {
		name := lib.Name
		if name == "~" {
			name = "imgscal"
		}
		fs, err := os.OpenFile(path.Join(TYPES_DIR, name+".lua"), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0o666)
		if err != nil {
			panic(err)
		}
		defer fs.Close()

		fmt.Fprintf(fs, "\n---@meta %s\n\n", name)
		fmt.Fprintf(fs, "---@class %s\n", name)

		for _, desc := range lib.Desc {
			fmt.Fprintf(fs, "---%s\n", desc)
		}
		fmt.Fprintf(fs, "%s = {}\n\n", name)

		for _, cn := range lib.Cns {
			alias := strings.TrimSpace(fmt.Sprintf("%s_%s", name, cn.Group))
			fmt.Fprintf(fs, "---@alias %s %s\n", alias, types.ParseType(cn.Type, name))

			entries := make([]string, len(cn.Consts))

			for i, c := range cn.Consts {
				split := strings.Fields(c)
				prop := split[0]
				desc := ""
				if len(split) > 1 {
					desc = " " + strings.Join(split[1:], " ")
				}
				fmt.Fprintf(fs, "---@alias %s_%s %s\n", name, prop, alias)
				fmt.Fprintf(fs, "---@class %s\n", name)
				fmt.Fprintf(fs, "---@field %s %s%s\n\n", prop, alias, desc)
				fmt.Fprintf(fs, "%s.%s = nil\n\n", name, prop)

				entries[i] = fmt.Sprintf("%s_%s", name, prop)
			}

			fmt.Fprintf(fs, "---@alias %s_* (%s)\n", alias, strings.Join(entries, " | "))
		}

		fmt.Fprintf(fs, "\n")

		for _, fn := range lib.Fns {
			formatFunction(fs, name, fn)
		}

		for _, st := range lib.Sts {
			formatStruct(fs, name, st)
		}

		for _, it := range lib.Its {
			formatInterface(fs, name, it)
		}

		if name == "imgscal" {
			fmt.Fprintf(fs, "---@alias imgscal_Imports\n")

			for _, lib := range docs {
				fmt.Fprintf(fs, "---| '\"%s\"' # %s\n", lib.Name, strings.Join(lib.Desc, " "))
			}
		}

		fmt.Fprintf(fs, "\nreturn %s\n", name)
	}
}

func formatFunction(fs io.Writer, lib string, fn doc.Fn) {
	argList := make([]string, len(fn.Args))

	for i, arg := range fn.Args {
		argname := strings.TrimSpace(arg.Str)
		if types.IsVariadic(arg.Type) {
			argname = "..."
		}

		opt := ""
		if arg.Opt {
			opt = "?"
		}

		fmt.Fprintf(fs, "---@param %s%s %s %s\n", argname, opt, types.ParseType(arg.Type, fn.Name), strings.TrimSpace(arg.Desc))
		argList[i] = argname
	}

	for _, ret := range fn.Returns {
		retdesc := ""
		if ret.Str != "" {
			retdesc = fmt.Sprintf(" # %s", strings.TrimPrefix(strings.TrimSpace(ret.Str), "- "))
		}

		fmt.Fprintf(fs, "---@return %s%s\n", types.ParseType(ret.Type, fn.Name), retdesc)
	}

	if fn.Block {
		fmt.Fprint(fs, "---@nodiscard\n")
	}

	for _, desc := range fn.Desc {
		fmt.Fprintf(fs, "---%s\n", desc)
	}

	fmt.Fprintf(fs, "function %s.%s(%s) end\n\n", lib, fn.Name, strings.Join(argList, ", "))
}

func formatStruct(fs io.Writer, lib string, st doc.Struct) {
	alias := fmt.Sprintf("%s_%s", lib, st.Struct)
	fmt.Fprintf(fs, "---@class %s\n", alias)

	for _, prop := range st.Props {
		propdesc := ""
		if prop.Desc != "" {
			propdesc = prop.Desc
		}
		fmt.Fprintf(fs, "---@field %s %s%s\n", strings.TrimSpace(prop.Str), types.ParseType(prop.Type, alias), propdesc)
	}

	for _, method := range st.Methods {
		methoddesc := ""
		if method.Desc != "" {
			methoddesc = " " + method.Desc
		}
		fmt.Fprintf(fs, "---@field %s %s%s\n", strings.TrimSpace(method.Name), types.ParseMethodType(method.Type, alias), methoddesc)
	}

	for _, desc := range st.Desc {
		fmt.Fprintf(fs, "---%s\n", desc)
	}

	fmt.Fprint(fs, "\n")
}

func formatInterface(fs io.Writer, lib string, it doc.Interface) {
	alias := fmt.Sprintf("%s_%s", lib, it.Interface)

	pairs := ""

	for _, prop := range it.Props {
		propdesc := ""
		if prop.Desc != "" {
			propdesc = fmt.Sprintf(" --[[%s]]", prop.Desc)
		}
		pairs += fmt.Sprintf("%s: %s%s,", prop.Str, types.ParseType(prop.Type, alias), propdesc)
	}

	for _, method := range it.Methods {
		methoddesc := ""
		if method.Desc != "" {
			methoddesc = fmt.Sprintf(" --[[%s]]", method.Desc)
		}
		pairs += fmt.Sprintf("%s: %s%s,", method.Name, types.ParseMethodType(method.Type, alias), methoddesc)
	}

	interfaceDesc := ""
	for _, desc := range it.Desc {
		interfaceDesc += fmt.Sprintf("<br>%s", desc)
	}

	fmt.Fprintf(fs, "---@alias %s { %s }%s\n", alias, pairs, interfaceDesc)
}
