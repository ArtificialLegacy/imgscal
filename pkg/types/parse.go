package types

import (
	"fmt"
	"strings"
)

const (
	escape_open  = '{'
	escape_close = '}'
)

func ParseEscaped(str string) (string, []string) {
	escaped := []string{}
	result := ""

	for i := 0; i < len(str); i++ {
		s := str[i]

		if s != escape_open {
			result += string(s)
			continue
		} else {
			result += "%s"

			current := ""
			escapeCount := 1
			i++

			for ; i < len(str) && escapeCount > 0; i++ {
				s = str[i]

				if s == escape_close {
					escapeCount--
				}

				if escapeCount > 0 {
					current += string(s)
				}
			}

			escaped = append(escaped, strings.TrimSpace(current))
			i--
		}
	}

	return result, escaped
}

var typeMap = map[string]string{
	"int":    "integer",
	"float":  "number",
	"bool":   "boolean",
	"string": "string",
	"any":    "any",
	"%s":     "%s",
}

func IsVariadic(str string) bool {
	return strings.HasSuffix(str, "...")
}

func ParseType(str string, self string) string {
	opt := ""
	if strings.HasSuffix(str, "?") {
		opt = "?"
		str = strings.TrimSuffix(str, "?")
	}

	if str == "self" {
		return self + opt
	}

	if strings.HasPrefix(str, "[]") {
		return fmt.Sprintf("%s[]%s", ParseType(strings.TrimPrefix(str, "[]"), self), opt)
	}

	if str == "table<any>" {
		return "table<any, any>" + opt
	}

	if strings.HasPrefix(str, "function") {
		return ParseMethodType(str, self) + opt
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

func ParseMethodType(str string, self string) string {
	var escaped []string
	str, escaped = ParseEscaped(str)
	for i, v := range escaped {
		escaped[i] = ParseType(v, self)
	}

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
			if IsVariadic(split[len(split)-1]) {
				argname = "..."
			}

			argsFormatted[i] = fmt.Sprintf("%s: %s", argname, ParseType(strings.Join(split, " "), self))
			continue
		}

		if len(split) == 1 {
			argname := fmt.Sprintf("arg%d", i)
			if IsVariadic(split[0]) {
				argname = "..."
			}
			argsFormatted[i] = fmt.Sprintf("%s: %s", argname, ParseType(split[0], self))
		} else if len(split) == 2 {
			argname := split[0]
			if IsVariadic(split[0]) {
				argname = "..."
			}
			argsFormatted[i] = fmt.Sprintf("%s: %s", argname, ParseType(split[1], self))
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
		retFormatted[i] = ParseType(strings.TrimSpace(r), self)
	}
	retString := ""
	if len(retFormatted) > 0 {
		retString = fmt.Sprintf(": %s", strings.Join(retFormatted, ", "))
	}

	return fmt.Sprintf(fmt.Sprintf("fun(%s)%s", strings.Join(argsFormatted, ", "), retString), anyArray(escaped)...)
}

func anyArray[T any](a []T) []any {
	result := make([]any, len(a))

	for k, v := range a {
		result[k] = v
	}

	return result
}
