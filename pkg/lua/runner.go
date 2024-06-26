package lua

import (
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/Shopify/go-lua"
)

type Runner struct {
	State *lua.State
	lg    *log.Logger

	TC *collection.Collection[collection.ItemTask]
	IC *collection.Collection[collection.ItemImage]
	FC *collection.Collection[collection.ItemFile]
	CC *collection.Collection[collection.ItemContext]
	QR *collection.Collection[collection.ItemQR]
}

func NewRunner(state *lua.State, lg *log.Logger) Runner {
	return Runner{
		State: state,
		lg:    lg,

		// collections
		IC: collection.NewCollection[collection.ItemImage](lg),
		FC: collection.NewCollection[collection.ItemFile](lg).OnCollect(
			func(i *collection.Item[collection.ItemFile]) {
				if i.Self != nil {
					i.Self.File.Close()
				}
			}),
		CC: collection.NewCollection[collection.ItemContext](lg),
		QR: collection.NewCollection[collection.ItemQR](lg),
		TC: collection.NewCollection[collection.ItemTask](lg),
	}
}

func (r *Runner) Run(file string) error {
	pwd, _ := os.Getwd()

	defer func() {
		if p := recover(); p != nil {
			r.lg.Append("recovered from panic during lua runtime.", log.LEVEL_ERROR)
		}
	}()

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
	FUNC
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

func (l *Lib) ParseArgs(name string, args []Arg, ln int) map[string]any {
	argMap := map[string]any{}

	for i, a := range args {
		ind := -ln + i
		switch a.Type {
		case INT:
			v, ok := l.State.ToInteger(ind)
			if i >= ln {
				ok = false
			}
			if (!ok || ind == 0) && !a.Optional {
				l.State.PushString(l.Lg.Append(fmt.Sprintf("invalid int provided to %s in arg pos %d", name, i), log.LEVEL_ERROR))
				l.State.Error()
			} else if (!ok || ind == 0) && a.Optional {
				argMap[a.Name] = l.getDefault(a)
			} else {
				rm := l.State.AbsIndex(ind)
				l.State.Remove(rm)
				argMap[a.Name] = v
			}

		case FLOAT:
			v, ok := l.State.ToNumber(ind)
			if i >= ln {
				ok = false
			}
			if (!ok || ind == 0) && !a.Optional {
				l.State.PushString(l.Lg.Append(fmt.Sprintf("invalid float provided to %s in arg pos %d", name, i), log.LEVEL_ERROR))
				l.State.Error()
			} else if (!ok || ind == 0) && a.Optional {
				argMap[a.Name] = l.getDefault(a)
			} else {
				rm := l.State.AbsIndex(ind)
				l.State.Remove(rm)
				argMap[a.Name] = v
			}

		case BOOL:
			if i >= ln {
				argMap[a.Name] = false
			} else {
				v := l.State.ToBoolean(ind)
				argMap[a.Name] = v
				rm := l.State.AbsIndex(ind)
				l.State.Remove(rm)
			}

		case STRING:
			v, ok := l.State.ToString(ind)
			if i >= ln {
				ok = false
			}
			if (!ok || ind == 0) && !a.Optional {
				l.State.PushString(l.Lg.Append(fmt.Sprintf("invalid string provided to %s in arg pos %d", name, i), log.LEVEL_ERROR))
				l.State.Error()
			} else if (!ok || ind == 0) && a.Optional {
				argMap[a.Name] = l.getDefault(a)
			} else {
				rm := l.State.AbsIndex(ind)
				l.State.Remove(rm)
				argMap[a.Name] = v
			}

		case TABLE:
			exists := l.State.IsTable(ind)
			if i >= ln {
				exists = false
			}
			if (!exists || ind == 0) && !a.Optional {
				l.State.PushString(l.Lg.Append(fmt.Sprintf("invalid table provided to %s in arg pos %d", name, i), log.LEVEL_ERROR))
				l.State.Error()
			} else if (!exists || ind == 0) && a.Optional {
				argMap[a.Name] = l.getDefault(a)
			} else {
				l.flattenTable(*a.Table)
				argMap[a.Name] = l.ParseArgs(name, *a.Table, len(*a.Table))
				rm := l.State.AbsIndex(ind)
				l.State.Remove(rm)
			}

		case ARRAY:
			exists := l.State.IsTable(ind)
			if i >= ln {
				exists = false
			}
			if (!exists || ind == 0) && !a.Optional {
				l.State.PushString(l.Lg.Append(fmt.Sprintf("invalid array provided to %s in arg pos %d", name, i), log.LEVEL_ERROR))
				l.State.Error()
			} else if (!exists || ind == 0) && a.Optional {
				argMap[a.Name] = l.getDefault(a)
			} else {
				ln := l.State.RawLength(ind)
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

				argMap[a.Name] = l.ParseArgs(name, argTable, ln)
				rm := l.State.AbsIndex(ind)
				l.State.Remove(rm)
			}

		case ANY:
			if i >= ln {
				argMap[a.Name] = nil
			} else {
				v := l.State.ToValue(ind)
				rm := l.State.AbsIndex(ind)
				l.State.Remove(rm)
				argMap[a.Name] = v
			}

		case FUNC:
			if i >= ln || !l.State.IsFunction(ind) {
				argMap[a.Name] = nil
			} else {
				aind := l.State.AbsIndex(ind)
				argMap[a.Name] = aind
			}
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

func (l *Lib) getDefault(a Arg) any {
	switch a.Type {
	case INT:
		return 0

	case FLOAT:
		return 0.0

	case BOOL:
		return false

	case STRING:
		return ""

	case TABLE:
		tab := map[string]any{}
		for _, v := range *a.Table {
			tab[v.Name] = l.getDefault(v)
		}
		return tab

	case ARRAY:
		return []any{}

	case ANY:
		fallthrough
	default:
		return nil
	}
}

type TaskData struct {
	Lib  string
	Name string
}

func (l *Lib) CreateFunction(name string, args []Arg, fn func(d TaskData, args map[string]any) int) {
	l.State.PushGoFunction(func(state *lua.State) int {
		l.Lg.Append(fmt.Sprintf("%s.%s called.", l.Lib, name), log.LEVEL_INFO)

		argMap := l.ParseArgs(name, args, l.State.Top())

		ret := fn(TaskData{Lib: l.Lib, Name: name}, argMap)
		l.Lg.Append(fmt.Sprintf("%s.%s finished.", l.Lib, name), log.LEVEL_INFO)
		return ret
	})
	l.State.SetField(-2, name)
}
