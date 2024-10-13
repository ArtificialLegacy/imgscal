package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/ArtificialLegacy/imgscal/pkg/doc"
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

		if name == "imgscal" {
			fmt.Fprintf(fs, "\n")
			for _, mod := range docs {
				if mod.Name == "~" {
					continue
				}
				fmt.Fprintf(fs, "---@module '%s'\n", mod.Name)
			}
		}

		fmt.Fprintf(fs, "\n---@meta %s\n\n", name)
		fmt.Fprintf(fs, "---@class %s\n", name)

		for _, cn := range lib.Cns {
			alias := strings.TrimSpace(fmt.Sprintf("%s.%s", name, cn.Group))
			fmt.Fprintf(fs, "---@alias %s %s\n", alias, parseType(cn.Type, name))

			entries := make([]string, len(cn.Consts))

			for i, c := range cn.Consts {
				split := strings.Fields(c)
				prop := split[0]
				desc := ""
				if len(split) > 1 {
					desc = " " + strings.Join(split[1:], " ")
				}
				fmt.Fprintf(fs, "---@alias %s.%s %s\n", name, prop, alias)
				fmt.Fprintf(fs, "---@field %s.%s %s%s\n", name, prop, alias, desc)

				entries[i] = fmt.Sprintf("%s.%s", name, prop)
			}

			fmt.Fprintf(fs, "---@alias %s.* (%s)\n", alias, strings.Join(entries, " | "))
		}

		for _, desc := range lib.Desc {
			fmt.Fprintf(fs, "---%s\n", desc)
		}
		fmt.Fprintf(fs, "%s = {}\n\n", name)

		for _, fn := range lib.Fns {
			formatFunction(fs, name, fn)
		}

		for _, st := range lib.Sts {
			formatStruct(fs, name, st)
		}

		for _, it := range lib.Its {
			formatInterface(fs, name, it)
		}

		fmt.Fprintf(fs, "return %s\n", name)
	}
}

func formatFunction(fs io.Writer, lib string, fn doc.Fn) {
	argList := make([]string, len(fn.Args))

	for i, arg := range fn.Args {
		argname := strings.TrimSpace(arg.Str)
		if isVariadic(arg.Type) {
			argname = "..."
		}

		opt := ""
		if arg.Opt {
			opt = "?"
		}

		fmt.Fprintf(fs, "---@param %s%s %s %s\n", argname, opt, parseType(arg.Type, fn.Name), strings.TrimSpace(arg.Desc))
		argList[i] = argname
	}

	for _, ret := range fn.Returns {
		retdesc := ""
		if ret.Str != "" {
			retdesc = fmt.Sprintf(" # %s", strings.TrimPrefix(strings.TrimSpace(ret.Str), "- "))
		}

		fmt.Fprintf(fs, "---@return %s%s\n", parseType(ret.Type, fn.Name), retdesc)
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
	alias := fmt.Sprintf("%s.%s", lib, st.Struct)
	fmt.Fprintf(fs, "---@class %s\n", alias)

	for _, prop := range st.Props {
		propdesc := ""
		if prop.Desc != "" {
			propdesc = prop.Desc
		}
		fmt.Fprintf(fs, "---@field %s %s%s\n", strings.TrimSpace(prop.Str), parseType(prop.Type, alias), propdesc)
	}

	for _, method := range st.Methods {
		methoddesc := ""
		if method.Desc != "" {
			methoddesc = " " + method.Desc
		}
		fmt.Fprintf(fs, "---@field %s %s%s\n", strings.TrimSpace(method.Name), parseMethodType(method.Type, alias), methoddesc)
	}

	for _, desc := range st.Desc {
		fmt.Fprintf(fs, "---%s\n", desc)
	}

	fmt.Fprint(fs, "\n")
}

func formatInterface(fs io.Writer, lib string, it doc.Interface) {
	alias := fmt.Sprintf("%s.%s", lib, it.Interface)

	pairs := ""

	for _, prop := range it.Props {
		propdesc := ""
		if prop.Desc != "" {
			propdesc = fmt.Sprintf(" --[[%s]]", prop.Desc)
		}
		pairs += fmt.Sprintf("%s: %s%s,", prop.Str, parseType(prop.Type, alias), propdesc)
	}

	for _, method := range it.Methods {
		methoddesc := ""
		if method.Desc != "" {
			methoddesc = fmt.Sprintf(" --[[%s]]", method.Desc)
		}
		pairs += fmt.Sprintf("%s: %s%s,", method.Name, parseMethodType(method.Type, alias), methoddesc)
	}

	interfaceDesc := ""
	for _, desc := range it.Desc {
		interfaceDesc += fmt.Sprintf("<br>%s", desc)
	}

	fmt.Fprintf(fs, "---@alias %s { %s }%s\n", alias, pairs, interfaceDesc)
}

var typeMap = map[string]string{
	"int":    "integer",
	"float":  "number",
	"bool":   "boolean",
	"string": "string",
	"any":    "any",
}

func isVariadic(str string) bool {
	return strings.HasSuffix(str, "...")
}

func parseType(str string, self string) string {
	opt := ""
	if strings.HasSuffix(str, "?") {
		opt = "?"
		str = strings.TrimSuffix(str, "?")
	}

	if str == "self" {
		return self + opt
	}

	if strings.HasPrefix(str, "[]") {
		return fmt.Sprintf("%s[]%s", parseType(strings.TrimPrefix(str, "[]"), self), opt)
	}

	if str == "table<any>" {
		return "table<any, any>" + opt
	}

	if strings.HasPrefix(str, "function") {
		return parseMethodType(str, self) + opt
	}

	if strings.HasPrefix(str, "struct") {
		alias := strings.FieldsFunc(str, func(r rune) bool {
			return r == '<' || r == '>'
		})

		return alias[1] + opt
	}

	for k, v := range typeMap {
		if strings.HasPrefix(str, k) {
			if strings.IndexByte(str, '<') != -1 {
				alias := strings.FieldsFunc(str, func(r rune) bool {
					return r == '<' || r == '>'
				})

				return alias[1] + opt
			}

			return v + opt
		}
	}

	fmt.Printf("unknown type: %s [%s]\n", str, self)
	return "any" + opt
}

func parseMethodType(str string, self string) string {
	argStart := strings.IndexByte(str, '(')
	argEnd := strings.LastIndexByte(str, ')')

	if argEnd == -1 {
		argEnd = len(str)
	}

	args := str[argStart+1 : argEnd]
	argList := strings.Split(args, ", ")
	if argStart+1 == argEnd {
		argList = []string{}
	}

	argsFormatted := make([]string, len(argList))

	for i, a := range argList {
		split := strings.Split(a, " ")

		if strings.HasPrefix(split[0], "function") || (len(split) > 1 && strings.HasPrefix(split[1], "function")) {
			argname := fmt.Sprintf("arg%d", i)
			if !strings.HasPrefix(split[0], "function") {
				argname = split[0]
				split = split[1:]
			}
			if isVariadic(split[len(split)-1]) {
				argname = "..."
			}

			argsFormatted[i] = fmt.Sprintf("%s: %s", argname, parseType(strings.Join(split, " "), self))
			continue
		}

		if len(split) == 1 {
			argname := fmt.Sprintf("arg%d", i)
			if isVariadic(split[0]) {
				argname = "..."
			}
			argsFormatted[i] = fmt.Sprintf("%s: %s", argname, parseType(split[0], self))
		} else if len(split) == 2 {
			argname := split[0]
			if isVariadic(split[0]) {
				argname = "..."
			}
			argsFormatted[i] = fmt.Sprintf("%s: %s", argname, parseType(split[1], self))
		}
	}

	retStart := strings.LastIndex(str, "->")
	returns := str[retStart+2:]
	if retStart == -1 {
		returns = ""
	}
	retList := strings.Split(returns, ", ")
	if retStart == -1 {
		retList = []string{}
	}

	retFormatted := make([]string, len(retList))

	for i, r := range retList {
		retFormatted[i] = parseType(strings.TrimSpace(r), self)
	}
	retString := ""
	if len(retFormatted) > 0 {
		retString = fmt.Sprintf(": %s", strings.Join(retFormatted, ", "))
	}

	return fmt.Sprintf("fun(%s)%s", strings.Join(argsFormatted, ", "), retString)
}
