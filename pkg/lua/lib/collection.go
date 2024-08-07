package lib

import (
	"fmt"
	"sync"

	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	golua "github.com/yuin/gopher-lua"
)

const LIB_COLLECTION = "collection"

func RegisterCollection(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_COLLECTION, r, r.State, lg)

	/// @func task()
	/// @arg name
	/// @returns id
	lib.CreateFunction(tab, "task",
		[]lua.Arg{
			{Type: lua.STRING, Name: "name"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {

			name := args["name"].(string)

			chLog := log.NewLogger(fmt.Sprintf("task_%s", name))
			chLog.Parent = lg
			lg.Append(fmt.Sprintf("child log created: task_%s", name), log.LEVEL_INFO)

			id := r.TC.AddItem(&chLog)

			r.TC.Schedule(id, &collection.Task[collection.ItemTask]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemTask]) {
					i.Self = &collection.ItemTask{
						Name: name,
					}
				},
			})

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func schedule()
	/// @arg type - collection type
	/// @arg id
	/// @arg func
	/// @desc
	/// Schedules a lua func to be called from the queue.
	lib.CreateFunction(tab, "schedule",
		[]lua.Arg{
			{Type: lua.INT, Name: "type"},
			{Type: lua.INT, Name: "id"},
			{Type: lua.FUNC, Name: "func"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			switch lua.ParseEnum(args["type"].(int), collection.CollectionList, lib) {
			case collection.TYPE_TASK:
				<-r.TC.Schedule(args["id"].(int), &collection.Task[collection.ItemTask]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemTask]) {
						callScheduledFunction(state, args["func"].(*golua.LFunction))
					},
				})
			case collection.TYPE_IMAGE:
				<-r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemImage]) {
						callScheduledFunction(state, args["func"].(*golua.LFunction))
					},
				})
			case collection.TYPE_FILE:
				<-r.FC.Schedule(args["id"].(int), &collection.Task[collection.ItemFile]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemFile]) {
						callScheduledFunction(state, args["func"].(*golua.LFunction))
					},
				})
			case collection.TYPE_CONTEXT:
				<-r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemContext]) {
						callScheduledFunction(state, args["func"].(*golua.LFunction))
					},
				})
			case collection.TYPE_QR:
				<-r.QR.Schedule(args["id"].(int), &collection.Task[collection.ItemQR]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemQR]) {
						callScheduledFunction(state, args["func"].(*golua.LFunction))
					},
				})
			}

			return 0
		})

	/// @func wait()
	/// @arg type - collection type
	/// @arg id - id of item in the collection
	/// @blocking
	lib.CreateFunction(tab, "wait",
		[]lua.Arg{
			{Type: lua.INT, Name: "type"},
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			switch lua.ParseEnum(args["type"].(int), collection.CollectionList, lib) {
			case collection.TYPE_TASK:
				<-r.TC.Schedule(args["id"].(int), &collection.Task[collection.ItemTask]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn:   func(i *collection.Item[collection.ItemTask]) {},
				})
			case collection.TYPE_IMAGE:
				<-r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn:   func(i *collection.Item[collection.ItemImage]) {},
				})
			case collection.TYPE_FILE:
				<-r.FC.Schedule(args["id"].(int), &collection.Task[collection.ItemFile]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn:   func(i *collection.Item[collection.ItemFile]) {},
				})
			case collection.TYPE_CONTEXT:
				<-r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn:   func(i *collection.Item[collection.ItemContext]) {},
				})
			case collection.TYPE_QR:
				<-r.QR.Schedule(args["id"].(int), &collection.Task[collection.ItemQR]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn:   func(i *collection.Item[collection.ItemQR]) {},
				})
			}

			return 0
		})

	/// @func wait_all()
	/// @arg type - collection type
	/// @blocking
	lib.CreateFunction(tab, "wait_all",
		[]lua.Arg{
			{Type: lua.INT, Name: "type"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			switch lua.ParseEnum(args["type"].(int), collection.CollectionList, lib) {
			case collection.TYPE_TASK:
				<-r.TC.ScheduleAll(&collection.Task[collection.ItemTask]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn:   func(i *collection.Item[collection.ItemTask]) {},
				})
			case collection.TYPE_IMAGE:
				<-r.IC.ScheduleAll(&collection.Task[collection.ItemImage]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn:   func(i *collection.Item[collection.ItemImage]) {},
				})
			case collection.TYPE_FILE:
				<-r.FC.ScheduleAll(&collection.Task[collection.ItemFile]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn:   func(i *collection.Item[collection.ItemFile]) {},
				})
			case collection.TYPE_CONTEXT:
				<-r.CC.ScheduleAll(&collection.Task[collection.ItemContext]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn:   func(i *collection.Item[collection.ItemContext]) {},
				})
			case collection.TYPE_QR:
				<-r.QR.ScheduleAll(&collection.Task[collection.ItemQR]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn:   func(i *collection.Item[collection.ItemQR]) {},
				})
			}

			return 0
		})

	/// @func wait_extensive()
	/// @blocking
	/// @desc
	/// This waits for all items across all collections to sync.
	lib.CreateFunction(tab, "wait_extensive",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			chans := []<-chan struct{}{}

			chans = append(chans, r.TC.ScheduleAll(&collection.Task[collection.ItemTask]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn:   func(i *collection.Item[collection.ItemTask]) {},
			}))
			chans = append(chans, r.IC.ScheduleAll(&collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn:   func(i *collection.Item[collection.ItemImage]) {},
			}))
			chans = append(chans, r.FC.ScheduleAll(&collection.Task[collection.ItemFile]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn:   func(i *collection.Item[collection.ItemFile]) {},
			}))
			chans = append(chans, r.CC.ScheduleAll(&collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn:   func(i *collection.Item[collection.ItemContext]) {},
			}))
			chans = append(chans, r.QR.ScheduleAll(&collection.Task[collection.ItemQR]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn:   func(i *collection.Item[collection.ItemQR]) {},
			}))

			wg := sync.WaitGroup{}
			wg.Add(len(chans))

			for _, co := range chans {
				c := co
				go func() {
					<-c
					wg.Done()
				}()
			}

			wg.Wait()
			return 0
		})

	/// @func collect()
	/// @arg type - collection type
	/// @arg id - id of item to collect early
	/// @desc
	/// items are collected automatically at the end of execution,
	/// but if this can be used to collect early on workflows that create many items.
	/// This is important for collections that open files, as they are only closed when collected.
	lib.CreateFunction(tab, "collect",
		[]lua.Arg{
			{Type: lua.INT, Name: "type"},
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			switch lua.ParseEnum(args["type"].(int), collection.CollectionList, lib) {
			case collection.TYPE_TASK:
				r.TC.Collect(args["id"].(int))
			case collection.TYPE_IMAGE:
				r.IC.Collect(args["id"].(int))
			case collection.TYPE_FILE:
				r.FC.Collect(args["id"].(int))
			case collection.TYPE_CONTEXT:
				r.CC.Collect(args["id"].(int))
			case collection.TYPE_QR:
				r.QR.Collect(args["id"].(int))
			}
			return 0
		})

	/// @func log()
	/// @arg type - collection type
	/// @arg id - id of the item to log to
	/// @arg msg
	lib.CreateFunction(tab, "log",
		[]lua.Arg{
			{Type: lua.INT, Name: "type"},
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "msg"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			msg := args["msg"].(string)
			switch lua.ParseEnum(args["type"].(int), collection.CollectionList, lib) {
			case collection.TYPE_TASK:
				<-r.TC.Schedule(args["id"].(int), &collection.Task[collection.ItemTask]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemTask]) {
						i.Lg.Append(fmt.Sprintf("lua log: %s", msg), log.LEVEL_INFO)
					},
				})
			case collection.TYPE_IMAGE:
				<-r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemImage]) {
						i.Lg.Append(fmt.Sprintf("lua log: %s", msg), log.LEVEL_INFO)
					},
				})
			case collection.TYPE_FILE:
				<-r.FC.Schedule(args["id"].(int), &collection.Task[collection.ItemFile]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemFile]) {
						i.Lg.Append(fmt.Sprintf("lua log: %s", msg), log.LEVEL_INFO)
					},
				})
			case collection.TYPE_CONTEXT:
				<-r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemContext]) {
						i.Lg.Append(fmt.Sprintf("lua log: %s", msg), log.LEVEL_INFO)
					},
				})
			case collection.TYPE_QR:
				<-r.QR.Schedule(args["id"].(int), &collection.Task[collection.ItemQR]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemQR]) {
						i.Lg.Append(fmt.Sprintf("lua log: %s", msg), log.LEVEL_INFO)
					},
				})
			}
			return 0
		})

	/// @func warn()
	/// @arg type - collection type
	/// @arg id - id of the item to log to
	/// @arg msg
	lib.CreateFunction(tab, "warn",
		[]lua.Arg{
			{Type: lua.INT, Name: "type"},
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "msg"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			msg := args["msg"].(string)
			switch lua.ParseEnum(args["type"].(int), collection.CollectionList, lib) {
			case collection.TYPE_TASK:
				<-r.TC.Schedule(args["id"].(int), &collection.Task[collection.ItemTask]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemTask]) {
						i.Lg.Append(fmt.Sprintf("lua warn: %s", msg), log.LEVEL_WARN)
					},
				})
			case collection.TYPE_IMAGE:
				<-r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemImage]) {
						i.Lg.Append(fmt.Sprintf("lua warn: %s", msg), log.LEVEL_WARN)
					},
				})
			case collection.TYPE_FILE:
				<-r.FC.Schedule(args["id"].(int), &collection.Task[collection.ItemFile]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemFile]) {
						i.Lg.Append(fmt.Sprintf("lua warn: %s", msg), log.LEVEL_WARN)
					},
				})
			case collection.TYPE_CONTEXT:
				<-r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemContext]) {
						i.Lg.Append(fmt.Sprintf("lua warn: %s", msg), log.LEVEL_WARN)
					},
				})
			case collection.TYPE_QR:
				<-r.QR.Schedule(args["id"].(int), &collection.Task[collection.ItemQR]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemQR]) {
						i.Lg.Append(fmt.Sprintf("lua warn: %s", msg), log.LEVEL_WARN)
					},
				})
			}
			return 0
		})

	/// @func panic()
	/// @arg type - collection type
	/// @arg id - id of the item to panic
	/// @arg msg
	lib.CreateFunction(tab, "panic",
		[]lua.Arg{
			{Type: lua.INT, Name: "type"},
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "msg"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			msg := args["msg"].(string)
			switch lua.ParseEnum(args["type"].(int), collection.CollectionList, lib) {
			case collection.TYPE_TASK:
				<-r.TC.Schedule(args["id"].(int), &collection.Task[collection.ItemTask]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemTask]) {
						state.Error(golua.LString(i.Lg.Append(fmt.Sprintf("lua panic: %s", msg), log.LEVEL_ERROR)), 0)
					},
				})
			case collection.TYPE_IMAGE:
				<-r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemImage]) {
						state.Error(golua.LString(i.Lg.Append(fmt.Sprintf("lua panic: %s", msg), log.LEVEL_ERROR)), 0)
					},
				})
			case collection.TYPE_FILE:
				<-r.FC.Schedule(args["id"].(int), &collection.Task[collection.ItemFile]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemFile]) {
						state.Error(golua.LString(i.Lg.Append(fmt.Sprintf("lua panic: %s", msg), log.LEVEL_ERROR)), 0)
					},
				})
			case collection.TYPE_CONTEXT:
				<-r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemContext]) {
						state.Error(golua.LString(i.Lg.Append(fmt.Sprintf("lua panic: %s", msg), log.LEVEL_ERROR)), 0)
					},
				})
			case collection.TYPE_QR:
				<-r.QR.Schedule(args["id"].(int), &collection.Task[collection.ItemQR]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemQR]) {
						state.Error(golua.LString(i.Lg.Append(fmt.Sprintf("lua panic: %s", msg), log.LEVEL_ERROR)), 0)
					},
				})
			}
			return 0
		})

	/// @constants Collection Types
	/// @const TYPE_TASK
	/// @const TYPE_IMAGE
	/// @const TYPE_FILE
	/// @const TYPE_CONTEXT
	/// @const TYPE_QR
	r.State.SetField(tab, "TYPE_TASK", golua.LNumber(collection.TYPE_TASK))
	r.State.SetField(tab, "TYPE_IMAGE", golua.LNumber(collection.TYPE_IMAGE))
	r.State.SetField(tab, "TYPE_FILE", golua.LNumber(collection.TYPE_FILE))
	r.State.SetField(tab, "TYPE_CONTEXT", golua.LNumber(collection.TYPE_CONTEXT))
	r.State.SetField(tab, "TYPE_QR", golua.LNumber(collection.TYPE_QR))
}

func callScheduledFunction(state *golua.LState, f *golua.LFunction) {
	state.Push(f)
	state.Call(0, 0)
}
