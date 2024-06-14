package lua

import (
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/ArtificialLegacy/imgscal/pkg/image"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/Shopify/go-lua"
)

type Runner struct {
	State *lua.State
	IC    *image.ImageCollection
	lg    *log.Logger
}

func NewRunner(state *lua.State, lg *log.Logger) Runner {
	return Runner{
		State: state,
		IC:    image.NewImageCollection(lg),
		lg:    lg,
	}
}

func (r *Runner) Run(file string) error {
	pwd, _ := os.Getwd()

	err := lua.DoFile(r.State, path.Join(pwd, file))
	if err != nil {
		return err
	}

	return nil
}

type Lib struct {
	Lib   string
	State *lua.State
	Lg    *log.Logger
}

func NewLib(lib string, state *lua.State, lg *log.Logger) *Lib {
	state.NewTable()
	state.SetGlobal(lib)
	state.Global(lib)

	return &Lib{
		Lib:   lib,
		State: state,
		Lg:    lg,
	}
}

type ArgType int

const (
	INT ArgType = iota
	FLOAT
	BOOL
	STRING
	TABLE
	ARRAY
	ANY
)

type Arg struct {
	Type     ArgType
	Name     string
	Optional bool
	Table    *[]Arg
}

type ArrayType struct {
	Type  ArgType
	Table *[]Arg
}

func ArgArray(name string, arrType ArrayType, optional bool) Arg {
	table := []Arg{
		{Type: arrType.Type, Name: name, Table: arrType.Table},
	}

	return Arg{
		Type:     ARRAY,
		Name:     name,
		Optional: optional,
		Table:    &table,
	}
}

func (l *Lib) ParseArgs(name string, args []Arg) map[string]any {
	argMap := map[string]any{}

	for i, a := range args {
		switch a.Type {
		case INT:
			v, ok := l.State.ToInteger(i - len(args))
			if !ok && !a.Optional {
				l.State.PushString(l.Lg.Append(fmt.Sprintf("invalid int provided to %s in arg pos %d", name, i), log.LEVEL_ERROR))
				l.State.Error()
			}
			rm := l.State.AbsIndex(i - len(args))
			l.State.Remove(rm)
			argMap[a.Name] = v

		case FLOAT:
			v, ok := l.State.ToNumber(i - len(args))
			if !ok && !a.Optional {
				l.State.PushString(l.Lg.Append(fmt.Sprintf("invalid float provided to %s in arg pos %d", name, i), log.LEVEL_ERROR))
				l.State.Error()
			}
			rm := l.State.AbsIndex(i - len(args))
			l.State.Remove(rm)
			argMap[a.Name] = v

		case BOOL:
			v := l.State.ToBoolean(i - len(args))
			argMap[a.Name] = v

		case STRING:
			v, ok := l.State.ToString(i - len(args))
			if !ok && !a.Optional {
				l.State.PushString(l.Lg.Append(fmt.Sprintf("invalid string provided to %s in arg pos %d", name, i), log.LEVEL_ERROR))
				l.State.Error()
			}
			rm := l.State.AbsIndex(i - len(args))
			l.State.Remove(rm)
			argMap[a.Name] = v

		case TABLE:
			exists := l.State.IsTable(i - len(args))
			if !exists && !a.Optional {
				l.State.PushString(l.Lg.Append(fmt.Sprintf("invalid table provided to %s in arg pos %d", name, i), log.LEVEL_ERROR))
				l.State.Error()
			} else if !exists && a.Optional {
				argMap[a.Name] = map[string]any{}
			} else {
				l.flattenTable(*a.Table)
				argMap[a.Name] = l.ParseArgs(name, *a.Table)
				rm := l.State.AbsIndex(i - len(args))
				l.State.Remove(rm)
			}

		case ARRAY:
			exists := l.State.IsTable(i - len(args))
			if !exists && !a.Optional {
				l.State.PushString(l.Lg.Append(fmt.Sprintf("invalid array provided to %s in arg pos %d", name, i), log.LEVEL_ERROR))
				l.State.Error()
			} else if !exists && a.Optional {
				argMap[a.Name] = []any{}
			} else {
				ln := l.State.RawLength(i - len(args))
				argTable := []Arg{}

				for i := 1; i <= ln; i++ {
					argTable = append(argTable, Arg{
						Type:  (*a.Table)[0].Type,
						Name:  fmt.Sprint(i),
						Table: (*a.Table)[0].Table,
					})
				}

				for i, arg := range argTable {
					ind, _ := strconv.ParseInt(arg.Name, 10, 64)
					l.State.PushInteger(int(ind))
					l.State.Table(-i - 2)
				}

				argMap[a.Name] = l.ParseArgs(name, argTable)
				rm := l.State.AbsIndex(i - len(args))
				l.State.Remove(rm)
			}

		case ANY:
			v := l.State.ToValue(i - len(args))
			rm := l.State.AbsIndex(i - len(args))
			l.State.Remove(rm)
			argMap[a.Name] = v
		default:
			panic(fmt.Sprintf("attempting to parse an arg with an unknown type: %d", a.Type))
		}
	}

	return argMap
}

func (l *Lib) flattenTable(args []Arg) {
	for i, arg := range args {
		l.State.Field(-i-1, arg.Name)
	}
}

func (l *Lib) CreateFunction(name string, args []Arg, fn func(state *lua.State, args map[string]any) int) {
	l.State.PushGoFunction(func(state *lua.State) int {
		l.Lg.Append(fmt.Sprintf("%s.%s called.", l.Lib, name), log.LEVEL_INFO)

		argMap := l.ParseArgs(name, args)

		return fn(state, argMap)
	})
	l.State.SetField(-2, name)
}
