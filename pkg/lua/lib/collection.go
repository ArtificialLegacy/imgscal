package lib

import (
	"sync"

	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
)

const LIB_COLLECTION = "collection"

func RegisterCollection(r *lua.Runner, lg *log.Logger) {
	lib := lua.NewLib(LIB_COLLECTION, r.State, lg)

	/// @func wait()
	/// @arg type - collection type
	/// @arg id - id of item in the collection
	/// @blocking
	lib.CreateFunction("wait",
		[]lua.Arg{
			{Type: lua.INT, Name: "type"},
			{Type: lua.INT, Name: "id"},
		},
		func(d lua.TaskData, args map[string]any) int {
			switch args["type"].(int) {
			case int(collection.TYPE_TASK):
				<-r.TC.Schedule(args["id"].(int), &collection.Task[collection.ItemTask]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn:   func(i *collection.Item[collection.ItemTask]) {},
				})
			case int(collection.TYPE_IMAGE):
				<-r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn:   func(i *collection.Item[collection.ItemImage]) {},
				})
			case int(collection.TYPE_FILE):
				<-r.FC.Schedule(args["id"].(int), &collection.Task[collection.ItemFile]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn:   func(i *collection.Item[collection.ItemFile]) {},
				})
			case int(collection.TYPE_CONTEXT):
				<-r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn:   func(i *collection.Item[collection.ItemContext]) {},
				})
			case int(collection.TYPE_QR):
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
	lib.CreateFunction("wait_all",
		[]lua.Arg{
			{Type: lua.INT, Name: "type"},
		},
		func(d lua.TaskData, args map[string]any) int {
			switch args["type"].(int) {
			case int(collection.TYPE_TASK):
				<-r.TC.ScheduleAll(&collection.Task[collection.ItemTask]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn:   func(i *collection.Item[collection.ItemTask]) {},
				})
			case int(collection.TYPE_IMAGE):
				<-r.IC.ScheduleAll(&collection.Task[collection.ItemImage]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn:   func(i *collection.Item[collection.ItemImage]) {},
				})
			case int(collection.TYPE_FILE):
				<-r.FC.ScheduleAll(&collection.Task[collection.ItemFile]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn:   func(i *collection.Item[collection.ItemFile]) {},
				})
			case int(collection.TYPE_CONTEXT):
				<-r.CC.ScheduleAll(&collection.Task[collection.ItemContext]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn:   func(i *collection.Item[collection.ItemContext]) {},
				})
			case int(collection.TYPE_QR):
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
	lib.CreateFunction("wait_extensive",
		[]lua.Arg{},
		func(d lua.TaskData, args map[string]any) int {
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
	lib.CreateFunction("collect",
		[]lua.Arg{
			{Type: lua.INT, Name: "type"},
			{Type: lua.INT, Name: "id"},
		},
		func(d lua.TaskData, args map[string]any) int {
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

	/// @constants Collection Types
	/// @const TYPE_TASK
	/// @const TYPE_IMAGE
	/// @const TYPE_FILE
	/// @const TYPE_CONTEXT
	/// @const TYPE_QR
	lib.State.PushInteger(int(collection.TYPE_TASK))
	lib.State.SetField(-2, "TYPE_TASK")
	lib.State.PushInteger(int(collection.TYPE_IMAGE))
	lib.State.SetField(-2, "TYPE_IMAGE")
	lib.State.PushInteger(int(collection.TYPE_FILE))
	lib.State.SetField(-2, "TYPE_FILE")
	lib.State.PushInteger(int(collection.TYPE_CONTEXT))
	lib.State.SetField(-2, "TYPE_CONTEXT")
	lib.State.PushInteger(int(collection.TYPE_QR))
	lib.State.SetField(-2, "TYPE_QR")
}
