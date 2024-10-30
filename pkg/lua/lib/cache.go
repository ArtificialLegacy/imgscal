package lib

import (
	"image"

	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	imageutil "github.com/ArtificialLegacy/imgscal/pkg/image_util"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	golua "github.com/yuin/gopher-lua"
)

const LIB_CACHE = "cache"

/// @lib Cache
/// @import cache
/// @desc
/// Allows for keeping images in memory, without the overhead of the task scheduler.

func RegisterCache(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_CACHE, r, r.State, lg)

	/// @func store(id, nocopy?) -> int<collection.CRATE_CACHEDIMAGE>
	/// @arg id {int<collection.IMAGE>} - ID of the image to cache.
	/// @arg? nocopy {bool}
	/// @return {int<collection.CRATE_CACHEDIMAGE>}
	/// @blocking
	/// @desc
	/// This stores an image in non-accessable storage. This allows the original image item to be reused without losing the image data.
	/// Cached images do not have a log file, and do not have a goroutine for scheduling.
	lib.CreateFunction(tab, "store",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.BOOL, Name: "nocopy", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var id int

			<-r.IC.Schedule(state, args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					model := i.Self.Model

					var img image.Image
					if args["nocopy"].(bool) {
						img = i.Self.Image
					} else {
						img = imageutil.CopyImage(i.Self.Image, model)
					}

					id = r.CR_CIM.Add(&collection.CachedImageItem{
						Image: img,
						Model: model,
					})
				},
			})

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func store_into(cid, id, nocopy?)
	/// @arg id {int<collection.IMAGE>} - ID of the image to cache.
	/// @arg cid {int<collection.CRATE_CACHEDIMAGE>} - Preexisting cache item to store the image into.
	/// @arg? nocopy {bool}
	/// @blocking
	/// @desc
	/// This stores an image in non-accessable storage. This allows the original image item to be reused without losing the image data.
	/// Cached images do not have a log file, and do not have a goroutine for scheduling.
	lib.CreateFunction(tab, "store_into",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "cid"},
			{Type: lua.BOOL, Name: "nocopy", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			<-r.IC.Schedule(state, args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					model := i.Self.Model

					var img image.Image
					if args["nocopy"].(bool) {
						img = i.Self.Image
					} else {
						img = imageutil.CopyImage(i.Self.Image, model)
					}

					citem, err := r.CR_CIM.Item(args["cid"].(int))
					if err != nil {
						lua.Error(state, lg.Appendf("failed to get cached image: %s", log.LEVEL_ERROR, err))
					}

					citem.Model = model
					citem.Image = img
				},
			})

			return 0
		})

	/// @func retrieve(cid, id, nocopy?)
	/// @arg cid {int<collection.CRATE_CACHEDIMAGE>} - ID of the image to retrieve from the cache.
	/// @arg id {int<collection.IMAGE>} - ID of the image item to move the cached image into.
	/// @arg? nocopy {bool}
	/// @blocking
	/// @desc
	/// This keeps the encoding and name of the image item.
	lib.CreateFunction(tab, "retrieve",
		[]lua.Arg{
			{Type: lua.INT, Name: "cid"},
			{Type: lua.INT, Name: "id"},
			{Type: lua.BOOL, Name: "nocopy", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			<-r.IC.Schedule(state, args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					citem, err := r.CR_CIM.Item(args["cid"].(int))
					if err != nil {
						lua.Error(state, lg.Appendf("failed to get cached image: %s", log.LEVEL_ERROR, err))
					}

					i.Self.Model = citem.Model
					if args["nocopy"].(bool) {
						i.Self.Image = citem.Image
					} else {
						i.Self.Image = imageutil.CopyImage(citem.Image, citem.Model)
					}
				},
			})
			return 0
		})

	/// @func retrieve_ext(cid, id, name, encoding, nocopy?)
	/// @arg cid {int<collection.CRATE_CACHEDIMAGE>} - ID of the image to retrieve from the cache.
	/// @arg id {int<collection.IMAGE>} - ID of the image item to move the cached image into.
	/// @arg name {string}
	/// @arg encoding {int<image.Encoding>}
	/// @arg? nocopy {bool}
	/// @blocking
	lib.CreateFunction(tab, "retrieve_ext",
		[]lua.Arg{
			{Type: lua.INT, Name: "cid"},
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "encoding"},
			{Type: lua.BOOL, Name: "nocopy", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			<-r.IC.Schedule(state, args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					citem, err := r.CR_CIM.Item(args["cid"].(int))
					if err != nil {
						lua.Error(state, lg.Appendf("failed to get cached image: %s", log.LEVEL_ERROR, err))
					}

					i.Self.Model = citem.Model
					if args["nocopy"].(bool) {
						i.Self.Image = citem.Image
					} else {
						i.Self.Image = imageutil.CopyImage(citem.Image, citem.Model)
					}

					i.Self.Name = args["name"].(string)
					i.Self.Encoding = lua.ParseEnum(args["encoding"].(int), imageutil.EncodingList, lib)
				},
			})
			return 0
		})

	/// @func remove(id)
	/// @arg id {int<collection.CRATE_CACHEDIMAGE>} - ID of the cached image to clean.
	/// @desc
	/// After this, the cached image cannot be used.
	lib.CreateFunction(tab, "remove",
		[]lua.Arg{
			{Type: lua.INT, Name: "cid"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CR_CIM.Clean(args["cid"].(int))
			return 0
		})

	/// @func empty(id)
	/// @arg id {int<collection.CRATE_CACHEDIMAGE>} - ID of the cached image to set to an empty image.
	/// @desc
	/// Sets the cached image to a 1px by 1px gray image.
	lib.CreateFunction(tab, "empty",
		[]lua.Arg{
			{Type: lua.INT, Name: "cid"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			citem, err := r.CR_CIM.Item(args["cid"].(int))
			if err != nil {
				lua.Error(state, lg.Appendf("failed to get cached image: %s", log.LEVEL_ERROR, err))
			}

			citem.Model = imageutil.MODEL_GRAY
			citem.Image = image.NewGray(image.Rect(0, 0, 1, 1))

			return 0
		})
}
