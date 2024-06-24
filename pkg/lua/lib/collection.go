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
			case int(collection.TYPE_IMAGE):
				<-r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn:   func(i *collection.Item[collection.ItemImage]) {},
				})
			case int(collection.TYPE_FILE):
				<-r.FC.Schedule(args["id"].(int), &collection.Task[collection.ItemFile]{
					Lib:  LIB_COLLECTION,
					Name: "wait",
					Fn:   func(i *collection.Item[collection.ItemFile]) {},
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
			case int(collection.TYPE_IMAGE):
				<-r.IC.ScheduleAll(&collection.Task[collection.ItemImage]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn:   func(i *collection.Item[collection.ItemImage]) {},
				})
			case int(collection.TYPE_FILE):
				<-r.FC.ScheduleAll(&collection.Task[collection.ItemFile]{
					Lib:  LIB_COLLECTION,
					Name: "wait_all",
					Fn:   func(i *collection.Item[collection.ItemFile]) {},
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
	/// normally this shouldn't be needed as items are collected automatically,
	/// but this can allow items to be collected earlier if collections are very large.
	lib.CreateFunction("collect",
		[]lua.Arg{
			{Type: lua.INT, Name: "type"},
			{Type: lua.INT, Name: "id"},
		},
		func(d lua.TaskData, args map[string]any) int {
			switch args["type"].(int) {
			case int(collection.TYPE_IMAGE):
				r.IC.Collect(args["id"].(int))
			case int(collection.TYPE_FILE):
				r.FC.Collect(args["id"].(int))
			}
			return 0
		})

	/// @constants Collection Types
	/// @const TYPE_TASK
	/// @const TYPE_IMAGE
	/// @const TYPE_FILE
	lib.State.PushInteger(int(collection.TYPE_TASK))
	lib.State.SetField(-2, "TYPE_TASK")
	lib.State.PushInteger(int(collection.TYPE_IMAGE))
	lib.State.SetField(-2, "TYPE_IMAGE")
	lib.State.PushInteger(int(collection.TYPE_FILE))
	lib.State.SetField(-2, "TYPE_FILE")
}
