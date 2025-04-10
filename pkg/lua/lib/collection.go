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

/// @lib Collection
/// @import collection
/// @desc
/// Utility library for manually interacting with ImgScal's concurrency system.

func RegisterCollection(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_COLLECTION, r, r.State, lg)

	/// @func task(name) -> int<collection.TASK>
	/// @arg name {string}
	/// @returns int<collection.TASK>
	/// @desc
	/// This creates an item in the task collection.
	/// This is a generic collection for any concurrency needs.
	lib.CreateFunction(tab, "task",
		[]lua.Arg{
			{Type: lua.STRING, Name: "name"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			name := args["name"].(string)

			chLog := log.NewLogger(fmt.Sprintf("task_%s", name), lg)
			lg.Append(fmt.Sprintf("child log created: task_%s", name), log.LEVEL_INFO)

			id := r.TC.AddItem(&chLog)

			r.TC.Schedule(state, id, &collection.Task[collection.ItemTask]{
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

	/// @func reference(type, id) -> id
	/// @arg type {int<collection.Type}
	/// @arg id {int<collection.Type.*} - An ID from the same collection as the above type.
	/// @returns {int<collection.Type.*} - ID for a new collection item from the above id.
	/// @desc
	/// This creates a new task queue that references the same data, this does not ensure thread safety.
	/// The reference may also get out of sync if the original item is overridden.
	lib.CreateFunction(tab, "reference",
		[]lua.Arg{
			{Type: lua.INT, Name: "type"},
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			id := args["id"].(int)
			var newId int

			switch lua.ParseEnum(args["type"].(int), collection.CollectionList, lib) {
			case collection.TYPE_TASK:
				item := r.TC.Item(id)
				if item == nil {
					lua.Error(state, lg.Appendf("cannot create reference; invalid item id: %d", log.LEVEL_ERROR, id))
				}

				newId = r.TC.ScheduleAdd(state, item.Self.Name, lg, d.Lib, d.Name, func(i *collection.Item[collection.ItemTask]) {
					i.Self = item.Self
				})
			case collection.TYPE_IMAGE:
				item := r.IC.Item(id)
				if item == nil {
					lua.Error(state, lg.Appendf("cannot create reference; invalid item id: %d", log.LEVEL_ERROR, id))
				}

				newId = r.IC.ScheduleAdd(state, item.Self.Name, lg, d.Lib, d.Name, func(i *collection.Item[collection.ItemImage]) {
					i.Self = item.Self
				})
			case collection.TYPE_CONTEXT:
				item := r.CC.Item(id)
				if item == nil {
					lua.Error(state, lg.Appendf("cannot create reference; invalid item id: %d", log.LEVEL_ERROR, id))
				}

				next := r.CC.Next()
				newId = r.CC.ScheduleAdd(state, fmt.Sprintf("context_ref%d_%d", id, next), lg, d.Lib, d.Name, func(i *collection.Item[collection.ItemContext]) {
					i.Self = item.Self
				})
			case collection.TYPE_QR:
				item := r.QR.Item(id)
				if item == nil {
					lua.Error(state, lg.Appendf("cannot create reference; invalid item id: %d", log.LEVEL_ERROR, id))
				}

				next := r.QR.Next()
				newId = r.QR.ScheduleAdd(state, fmt.Sprintf("qr_ref%d_%d", id, next), lg, d.Lib, d.Name, func(i *collection.Item[collection.ItemQR]) {
					i.Self = item.Self
				})
			}

			state.Push(golua.LNumber(newId))
			return 1
		})

	/// @func schedule(type, id, func)
	/// @arg type {int<collection.Type>}
	/// @arg id {int<collection.Type.*>} - An ID from the same collection as the above type.
	/// @arg func {function()}
	/// @desc
	/// Schedules a lua func to be called from the queue.
	lib.CreateFunction(tab, "schedule",
		[]lua.Arg{
			{Type: lua.INT, Name: "type"},
			{Type: lua.INT, Name: "id"},
			{Type: lua.FUNC, Name: "func"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			id := args["id"].(int)

			var scheduledState *golua.LState

			switch lua.ParseEnum(args["type"].(int), collection.CollectionList, lib) {
			case collection.TYPE_TASK:
				r.TC.Schedule(state, id, &collection.Task[collection.ItemTask]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemTask]) {
						scheduledState = collection.NewThread(state, id, collection.TYPE_TASK)
						callScheduledFunction(scheduledState, args["func"].(*golua.LFunction))
					},
					Fail: func(i *collection.Item[collection.ItemTask]) {
						scheduledState.Close()
					},
				})
			case collection.TYPE_IMAGE:
				r.IC.Schedule(state, id, &collection.Task[collection.ItemImage]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemImage]) {
						scheduledState = collection.NewThread(state, id, collection.TYPE_IMAGE)
						callScheduledFunction(scheduledState, args["func"].(*golua.LFunction))
					},
					Fail: func(i *collection.Item[collection.ItemImage]) {
						scheduledState.Close()
					},
				})
			case collection.TYPE_CONTEXT:
				r.CC.Schedule(state, id, &collection.Task[collection.ItemContext]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemContext]) {
						scheduledState = collection.NewThread(state, id, collection.TYPE_CONTEXT)
						callScheduledFunction(scheduledState, args["func"].(*golua.LFunction))
					},
					Fail: func(i *collection.Item[collection.ItemContext]) {
						scheduledState.Close()
					},
				})
			case collection.TYPE_QR:
				r.QR.Schedule(state, id, &collection.Task[collection.ItemQR]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemQR]) {
						scheduledState = collection.NewThread(state, id, collection.TYPE_QR)
						callScheduledFunction(scheduledState, args["func"].(*golua.LFunction))
					},
					Fail: func(i *collection.Item[collection.ItemQR]) {
						scheduledState.Close()
					},
				})
			}

			return 0
		})

	/// @func wait(type, id)
	/// @arg type {int<collection.Type>}
	/// @arg id {int<collection.Type.*>} - An ID from the same collection as the above type.
	/// @blocking
	lib.CreateFunction(tab, "wait",
		[]lua.Arg{
			{Type: lua.INT, Name: "type"},
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			switch lua.ParseEnum(args["type"].(int), collection.CollectionList, lib) {
			case collection.TYPE_TASK:
				<-r.TC.Schedule(state, args["id"].(int), &collection.Task[collection.ItemTask]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn:   func(i *collection.Item[collection.ItemTask]) {},
				})
			case collection.TYPE_IMAGE:
				<-r.IC.Schedule(state, args["id"].(int), &collection.Task[collection.ItemImage]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn:   func(i *collection.Item[collection.ItemImage]) {},
				})
			case collection.TYPE_CONTEXT:
				<-r.CC.Schedule(state, args["id"].(int), &collection.Task[collection.ItemContext]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn:   func(i *collection.Item[collection.ItemContext]) {},
				})
			case collection.TYPE_QR:
				<-r.QR.Schedule(state, args["id"].(int), &collection.Task[collection.ItemQR]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn:   func(i *collection.Item[collection.ItemQR]) {},
				})
			}

			return 0
		})

	/// @func wait_all(type)
	/// @arg type {int<collection.Type>}
	/// @blocking
	lib.CreateFunction(tab, "wait_all",
		[]lua.Arg{
			{Type: lua.INT, Name: "type"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			switch lua.ParseEnum(args["type"].(int), collection.CollectionList, lib) {
			case collection.TYPE_TASK:
				<-r.TC.ScheduleAll(state, &collection.Task[collection.ItemTask]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn:   func(i *collection.Item[collection.ItemTask]) {},
				})
			case collection.TYPE_IMAGE:
				<-r.IC.ScheduleAll(state, &collection.Task[collection.ItemImage]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn:   func(i *collection.Item[collection.ItemImage]) {},
				})
			case collection.TYPE_CONTEXT:
				<-r.CC.ScheduleAll(state, &collection.Task[collection.ItemContext]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn:   func(i *collection.Item[collection.ItemContext]) {},
				})
			case collection.TYPE_QR:
				<-r.QR.ScheduleAll(state, &collection.Task[collection.ItemQR]{
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
	/// Realistically, shouldn't be used.
	lib.CreateFunction(tab, "wait_extensive",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			chans := []<-chan struct{}{}

			chans = append(chans, r.TC.ScheduleAll(state, &collection.Task[collection.ItemTask]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn:   func(i *collection.Item[collection.ItemTask]) {},
			}))
			chans = append(chans, r.IC.ScheduleAll(state, &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn:   func(i *collection.Item[collection.ItemImage]) {},
			}))
			chans = append(chans, r.CC.ScheduleAll(state, &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn:   func(i *collection.Item[collection.ItemContext]) {},
			}))
			chans = append(chans, r.QR.ScheduleAll(state, &collection.Task[collection.ItemQR]{
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

	/// @func collect(type, id)
	/// @arg type {int<collection.Type>}
	/// @arg id {int<collection.Type.*>} - An ID from the same collection as the above type.
	/// @desc
	/// Items are collected automatically at the end of execution,
	/// but this can be used to collect early in workflows that create a large amount of items.
	/// This is important for collections that open files, as they are only closed when collected.
	lib.CreateFunction(tab, "collect",
		[]lua.Arg{
			{Type: lua.INT, Name: "type"},
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			switch lua.ParseEnum(args["type"].(int), collection.CollectionList, lib) {
			case collection.TYPE_TASK:
				r.TC.Collect(state, args["id"].(int))
			case collection.TYPE_IMAGE:
				r.IC.Collect(state, args["id"].(int))
			case collection.TYPE_CONTEXT:
				r.CC.Collect(state, args["id"].(int))
			case collection.TYPE_QR:
				r.QR.Collect(state, args["id"].(int))
			}
			return 0
		})

	/// @func exists(type, id) -> bool
	/// @arg type {int<collection.Type>}
	/// @arg id {int<collection.Type.*>} - An ID from the same collection as the above type.
	/// @returns {bool} - If the item exists, and has not been collected.
	/// @desc
	/// Note that this is non-blocking, and can return true for items that get collected soon after.
	lib.CreateFunction(tab, "exists",
		[]lua.Arg{
			{Type: lua.INT, Name: "type"},
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var exists bool

			switch lua.ParseEnum(args["type"].(int), collection.CollectionList, lib) {
			case collection.TYPE_TASK:
				exists = r.TC.ItemExists(args["id"].(int))
			case collection.TYPE_IMAGE:
				exists = r.IC.ItemExists(args["id"].(int))
			case collection.TYPE_CONTEXT:
				exists = r.CC.ItemExists(args["id"].(int))
			case collection.TYPE_QR:
				exists = r.QR.ItemExists(args["id"].(int))
			}

			state.Push(golua.LBool(exists))
			return 1
		})

	/// @func log(type, id, msg)
	/// @arg type {int<collection.Type>}
	/// @arg id {int<collection.Type.*>} - An ID from the same collection as the above type.
	/// @arg msg {string}
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
				r.TC.Schedule(state, args["id"].(int), &collection.Task[collection.ItemTask]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemTask]) {
						i.Lg.Append(fmt.Sprintf("lua log: %s", msg), log.LEVEL_INFO)
					},
				})
			case collection.TYPE_IMAGE:
				r.IC.Schedule(state, args["id"].(int), &collection.Task[collection.ItemImage]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemImage]) {
						i.Lg.Append(fmt.Sprintf("lua log: %s", msg), log.LEVEL_INFO)
					},
				})
			case collection.TYPE_CONTEXT:
				r.CC.Schedule(state, args["id"].(int), &collection.Task[collection.ItemContext]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemContext]) {
						i.Lg.Append(fmt.Sprintf("lua log: %s", msg), log.LEVEL_INFO)
					},
				})
			case collection.TYPE_QR:
				r.QR.Schedule(state, args["id"].(int), &collection.Task[collection.ItemQR]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemQR]) {
						i.Lg.Append(fmt.Sprintf("lua log: %s", msg), log.LEVEL_INFO)
					},
				})
			}
			return 0
		})

	/// @func warn(type, id, msg)
	/// @arg type {int<collection.Type>}
	/// @arg id {int<collection.Type.*>} - An ID from the same collection as the above type.
	/// @arg msg {string}
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
				r.TC.Schedule(state, args["id"].(int), &collection.Task[collection.ItemTask]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemTask]) {
						i.Lg.Append(fmt.Sprintf("lua warn: %s", msg), log.LEVEL_WARN)
					},
				})
			case collection.TYPE_IMAGE:
				r.IC.Schedule(state, args["id"].(int), &collection.Task[collection.ItemImage]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemImage]) {
						i.Lg.Append(fmt.Sprintf("lua warn: %s", msg), log.LEVEL_WARN)
					},
				})
			case collection.TYPE_CONTEXT:
				r.CC.Schedule(state, args["id"].(int), &collection.Task[collection.ItemContext]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemContext]) {
						i.Lg.Append(fmt.Sprintf("lua warn: %s", msg), log.LEVEL_WARN)
					},
				})
			case collection.TYPE_QR:
				r.QR.Schedule(state, args["id"].(int), &collection.Task[collection.ItemQR]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemQR]) {
						i.Lg.Append(fmt.Sprintf("lua warn: %s", msg), log.LEVEL_WARN)
					},
				})
			}
			return 0
		})

	/// @func panic(type, id, msg)
	/// @arg type {int<collection.Type>}
	/// @arg id {int<collection.Type.*>} - An ID from the same collection as the above type.
	/// @arg msg {string}
	/// @blocking
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
				<-r.TC.Schedule(state, args["id"].(int), &collection.Task[collection.ItemTask]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemTask]) {
						state.Error(golua.LString(i.Lg.Append(fmt.Sprintf("lua panic: %s", msg), log.LEVEL_ERROR)), 0)
					},
				})
			case collection.TYPE_IMAGE:
				<-r.IC.Schedule(state, args["id"].(int), &collection.Task[collection.ItemImage]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemImage]) {
						state.Error(golua.LString(i.Lg.Append(fmt.Sprintf("lua panic: %s", msg), log.LEVEL_ERROR)), 0)
					},
				})
			case collection.TYPE_CONTEXT:
				<-r.CC.Schedule(state, args["id"].(int), &collection.Task[collection.ItemContext]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemContext]) {
						state.Error(golua.LString(i.Lg.Append(fmt.Sprintf("lua panic: %s", msg), log.LEVEL_ERROR)), 0)
					},
				})
			case collection.TYPE_QR:
				<-r.QR.Schedule(state, args["id"].(int), &collection.Task[collection.ItemQR]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemQR]) {
						state.Error(golua.LString(i.Lg.Append(fmt.Sprintf("lua panic: %s", msg), log.LEVEL_ERROR)), 0)
					},
				})
			}
			return 0
		})

	/// @constants Type {int}
	/// @const TASK
	/// @const IMAGE
	/// @const CONTEXT
	/// @const QR
	tab.RawSetString("TASK", golua.LNumber(collection.TYPE_TASK))
	tab.RawSetString("IMAGE", golua.LNumber(collection.TYPE_IMAGE))
	tab.RawSetString("CONTEXT", golua.LNumber(collection.TYPE_CONTEXT))
	tab.RawSetString("QR", golua.LNumber(collection.TYPE_QR))

	/// @constants Crate {int}
	/// @const CRATE_WINDOW
	/// @const CRATE_REF
	/// @const CRATE_GAMEMAKER
	/// @const CRATE_TEA
	/// @const CRATE_LIPGLOSS
	/// @const CRATE_CACHEDIMAGE
	/// @const CRATE_SHADER
	tab.RawSetString("CRATE_WINDOW", golua.LNumber(0))
	tab.RawSetString("CRATE_REF", golua.LNumber(1))
	tab.RawSetString("CRATE_GAMEMAKER", golua.LNumber(2))
	tab.RawSetString("CRATE_TEA", golua.LNumber(3))
	tab.RawSetString("CRATE_LIPGLOSS", golua.LNumber(4))
	tab.RawSetString("CRATE_CACHEDIMAGE", golua.LNumber(5))
	tab.RawSetString("CRATE_SHADER", golua.LNumber(6))
}

func callScheduledFunction(state *golua.LState, f *golua.LFunction) {
	state.Push(f)
	state.Call(0, 0)
	state.Close()
}
