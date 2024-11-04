package collection

import (
	"context"
	"fmt"
	"image"
	"runtime/debug"
	"sync"

	imageutil "github.com/ArtificialLegacy/imgscal/pkg/image_util"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/fogleman/gg"
	"github.com/skip2/go-qrcode"
	golua "github.com/yuin/gopher-lua"
)

const TASK_QUEUE_SIZE = 64

type CollectionType int

const (
	TYPE_TASK CollectionType = iota
	TYPE_IMAGE
	TYPE_CONTEXT
	TYPE_QR
)

var CollectionList = []CollectionType{
	TYPE_TASK,
	TYPE_IMAGE,
	TYPE_CONTEXT,
	TYPE_QR,
}

type ItemSelf interface {
	Identifier() CollectionType
}

type ItemImage struct {
	Image    image.Image
	Name     string
	Encoding imageutil.ImageEncoding
	Model    imageutil.ColorModel
}

func (img ItemImage) Identifier() CollectionType { return TYPE_IMAGE }

type ItemContext struct {
	Context *gg.Context
}

func (context ItemContext) Identifier() CollectionType { return TYPE_CONTEXT }

type ItemQR struct {
	QR *qrcode.QRCode
}

func (qr ItemQR) Identifier() CollectionType { return TYPE_QR }

type ItemTask struct {
	Name string
}

func (qr ItemTask) Identifier() CollectionType { return TYPE_TASK }

type Item[T ItemSelf] struct {
	Self *T
	Lg   *log.Logger
	wg   *sync.WaitGroup

	cleaned bool
	collect bool

	currTask *Task[T]

	failed bool
	Err    error

	TaskQueue chan *Task[T]
}

func NewItem[T ItemSelf](lg *log.Logger, wg *sync.WaitGroup, fn func(i *Item[T])) *Item[T] {
	i := &Item[T]{
		Self:      nil,
		Lg:        lg,
		cleaned:   false,
		collect:   false,
		failed:    false,
		TaskQueue: make(chan *Task[T], TASK_QUEUE_SIZE),
		wg:        wg,
	}

	go i.process(fn)

	return i
}

func (i *Item[T]) process(fn func(i *Item[T])) {
	defer func() {
		if p := recover(); p != nil {
			i.Lg.Append(fmt.Sprintf("recovered from panic within collection item: %+v\n%s", p, debug.Stack()), log.LEVEL_ERROR)
			if fn != nil {
				fn(i)
			}
			i.failed = true
			i.cleaned = true
			i.Err = fmt.Errorf("%+v", p)
			if i.currTask != nil {
				i.wg.Done()
				if i.currTask.Fail != nil {
					i.currTask.Fail(i)
				}
			}
			for c := 0; len(i.TaskQueue) > 0; c++ {
				task := <-i.TaskQueue
				i.wg.Done()
				if task.Fail != nil {
					task.Fail(i)
				}
				i.Lg.Append(fmt.Sprintf("drained task %d from item [%T]", c, i.Self), log.LEVEL_WARN)
			}
		}

		i.Lg.Close()
	}()

	for {
		if len(i.TaskQueue) == 0 {
			continue
		}

		i.currTask = <-i.TaskQueue
		i.Lg.Append(fmt.Sprintf("%s.%s task called", i.currTask.Lib, i.currTask.Name), log.LEVEL_VERBOSE)
		i.currTask.Fn(i)
		i.Lg.Append(fmt.Sprintf("%s.%s task finished", i.currTask.Lib, i.currTask.Name), log.LEVEL_VERBOSE)
		i.currTask = nil
		i.wg.Done()

		if i.cleaned {
			i.Lg.Append(fmt.Sprintf("item [%T] cleaned", i.Self), log.LEVEL_INFO)
			break
		}
	}
}

type Task[T ItemSelf] struct {
	Fn   func(i *Item[T])
	Fail func(i *Item[T])

	Lib  string
	Name string
}

type Collection[T ItemSelf] struct {
	items []*Item[T]
	lg    *log.Logger

	wg *sync.WaitGroup

	Errs []error

	onCollect func(i *Item[T])

	Identifier CollectionType
}

func NewCollection[T ItemSelf](lg *log.Logger, wg *sync.WaitGroup, identifier CollectionType) *Collection[T] {
	return &Collection[T]{
		items:      []*Item[T]{},
		lg:         lg,
		Errs:       []error{},
		wg:         wg,
		Identifier: identifier,
	}
}

func (c *Collection[T]) OnCollect(fn func(i *Item[T])) *Collection[T] {
	c.onCollect = fn
	return c
}

func (c *Collection[T]) AddItem(lg *log.Logger) int {
	item := NewItem(lg, c.wg, c.onCollect)
	id := len(c.items)

	c.items = append(c.items, item)

	return id
}

func (c *Collection[T]) Item(id int) *Item[T] {
	if id < 0 && id >= len(c.items) {
		c.lg.Append(fmt.Sprintf("invald item index: %d range of 0-%d", id, len(c.items)), log.LEVEL_WARN)
		return nil
	}

	item := c.items[id]
	return item
}

func (c *Collection[T]) ItemExists(id int) bool {
	if id < 0 && id >= len(c.items) {
		return false
	}

	item := c.items[id]
	if item == nil {
		return false
	}

	return !item.collect && !item.cleaned && !item.failed
}

func (c *Collection[T]) ScheduleAdd(state *golua.LState, name string, lg *log.Logger, tl, tn string, fn func(i *Item[T])) int {
	chLog := log.NewLogger(fmt.Sprintf("image_%s", name), lg)
	lg.Append(fmt.Sprintf("child log created: image_%s", name), log.LEVEL_INFO)

	id := c.AddItem(&chLog)

	c.Schedule(state, id, &Task[T]{
		Lib:  tl,
		Name: tn,
		Fn: func(i *Item[T]) {
			fn(i)
		},
	})

	return id
}

func (c *Collection[T]) Schedule(state *golua.LState, id int, tk *Task[T]) <-chan struct{} {
	wait := make(chan struct{}, 2)

	if !c.IDValid(id) {
		c.lg.Append(fmt.Sprintf("invald item index: %d range of 0-%d", id, len(c.items)), log.LEVEL_WARN)
		if tk.Fail != nil {
			tk.Fail(nil)
		}
		wait <- struct{}{}
		return wait
	}

	task := &Task[T]{
		Lib:  tk.Lib,
		Name: tk.Name,
		Fn: func(i *Item[T]) {
			tk.Fn(i)
			wait <- struct{}{}
		},
		Fail: func(i *Item[T]) {
			if tk.Fail != nil {
				tk.Fail(i)
			}
			wait <- struct{}{}
		},
	}

	item := c.items[id]

	if item.failed {
		item.Lg.Append(fmt.Sprintf("cannot schedule task for failed item: %d (%s.%s)", id, tk.Lib, tk.Name), log.LEVEL_WARN)
		task.Fail(item)
		return wait
	}

	ctx := state.Context()
	nested := SearchContext(ctx, id, c.Identifier)

	if !nested {
		item.Lg.Append(fmt.Sprintf("task scheduled for %d (%s.%s)", id, tk.Lib, tk.Name), log.LEVEL_VERBOSE)
		c.wg.Add(1)
		item.TaskQueue <- task
	} else {
		item.Lg.Append(fmt.Sprintf("task skipped scheduling, already within nested schedule: %d (%s.%s)", id, tk.Lib, tk.Name), log.LEVEL_VERBOSE)
		defer func() {
			if p := recover(); p != nil {
				task.Fail(item)
			}
		}()
		task.Fn(item)
	}

	return wait
}

type ContextKey string

const CONTEXT_STACK ContextKey = "__scheduleStack"

type StackContext struct {
	Identifier CollectionType
	ID         int
}

func CreateContext(state *golua.LState) {
	ctx := state.Context()
	if ctx == nil {
		ctx = context.TODO()
	}

	state.SetContext(context.WithValue(ctx, CONTEXT_STACK, []StackContext{}))
}

func NewThread(state *golua.LState, id int, identifier CollectionType) *golua.LState {
	thread, _ := state.NewThread()
	tctx := thread.Context()
	ctx := AddContext(tctx, id, identifier)
	thread.SetContext(ctx)

	return thread
}

func AddContext(ctx context.Context, id int, identifier CollectionType) context.Context {
	value := ctx.Value(CONTEXT_STACK)
	stack, ok := value.([]StackContext)
	copy(stack, stack)

	if !ok {
		stack = []StackContext{{Identifier: identifier, ID: id}}
	} else {
		stack = append(stack, StackContext{Identifier: identifier, ID: id})
	}

	newCtx := context.WithValue(ctx, CONTEXT_STACK, stack)
	return newCtx
}

func SearchContext(ctx context.Context, id int, identifier CollectionType) bool {
	value := ctx.Value(CONTEXT_STACK)
	stack, ok := value.([]StackContext)

	if !ok {
		return false
	}

	for _, v := range stack {
		if v.Identifier == identifier && v.ID == id {
			return true
		}
	}

	return false
}

func (c *Collection[T]) IDValid(id int) bool {
	return id >= 0 && id < len(c.items)
}

func (c *Collection[T]) SchedulePipe(state *golua.LState, id1, id2 int, tk1, tk2 *Task[T]) <-chan struct{} {
	wait := make(chan struct{}, 2)

	if !c.IDValid(id1) || !c.IDValid(id2) {
		c.lg.Append(fmt.Sprintf("invald item index: (1: %d) (2: %d) range of 0-%d", id1, id2, len(c.items)), log.LEVEL_WARN)
		if tk1.Fail != nil {
			tk1.Fail(nil)
		}
		if tk2.Fail != nil {
			tk2.Fail(nil)
		}
		wait <- struct{}{}
		return wait
	}

	if id1 != id2 {
		ready := make(chan struct{}, 2)
		finished := make(chan struct{}, 2)

		item1 := c.items[id1]
		item2 := c.items[id2]

		if item1.failed || item2.failed {
			if item1.failed {
				item1.Lg.Append(fmt.Sprintf("cannot schedule task for failed item: %d", id1), log.LEVEL_WARN)
			} else {
				item2.Lg.Append(fmt.Sprintf("cannot schedule task for failed item: %d", id2), log.LEVEL_WARN)
			}
			tk1.Fail(item1)
			tk2.Fail(item2)
			wait <- struct{}{}
			return wait
		}

		task1 := &Task[T]{
			Lib:  tk1.Lib,
			Name: tk1.Name,
			Fn: func(i *Item[T]) {
				tk1.Fn(i)
				ready <- struct{}{}
				<-finished
			},
			Fail: func(i *Item[T]) {
				if tk1.Fail != nil {
					tk1.Fail(i)
				}
				ready <- struct{}{}
			},
		}
		task2 := &Task[T]{
			Lib:  tk2.Lib,
			Name: tk2.Name,
			Fn: func(i *Item[T]) {
				<-ready
				tk2.Fn(i)
				finished <- struct{}{}
				wait <- struct{}{}
			},
			Fail: func(i *Item[T]) {
				if tk2.Fail != nil {
					tk2.Fail(i)
				}
				finished <- struct{}{}
				wait <- struct{}{}
			},
		}

		c.Schedule(state, id1, task1)
		c.Schedule(state, id2, task2)
	} else {
		item := c.items[id1]

		if item.failed {
			item.Lg.Append(fmt.Sprintf("cannot schedule task for failed item: %d", id1), log.LEVEL_WARN)
			tk1.Fail(item)
			wait <- struct{}{}
			return wait
		}

		task := &Task[T]{
			Lib:  tk1.Lib,
			Name: tk1.Name,
			Fn: func(i *Item[T]) {
				tk1.Fn(i)
				tk2.Fn(i)
			},
			Fail: func(i *Item[T]) {
				if tk1.Fail != nil {
					tk1.Fail(i)
				}
				if tk2.Fail != nil {
					tk2.Fail(i)
				}
				wait <- struct{}{}
			},
		}

		c.Schedule(state, id1, task)
	}

	return wait
}

func (c *Collection[T]) ScheduleAll(state *golua.LState, tk *Task[T]) <-chan struct{} {
	c.lg.Append(fmt.Sprintf("tasks scheduled for all items: [%T]", c), log.LEVEL_VERBOSE)

	wait := make(chan struct{}, 1)
	wg := sync.WaitGroup{}

	for id, i := range c.items {
		if i.collect {
			continue
		}

		wg.Add(1)
		task := &Task[T]{
			Lib:  tk.Lib,
			Name: tk.Name,
			Fn: func(i *Item[T]) {
				tk.Fn(i)
				wg.Done()
			},
		}

		c.Schedule(state, id, task)
	}

	go func() {
		wg.Wait()
		wait <- struct{}{}
	}()

	return wait
}

func (c *Collection[T]) CollectAll(state *golua.LState) {
	for id := range c.items {
		c.Collect(state, id)
	}
}

func (c *Collection[T]) Collect(state *golua.LState, id int) {
	i := c.items[id]
	if i.collect || i.failed || i.cleaned {
		return
	}

	i.collect = true

	c.lg.Append(fmt.Sprintf("item %d collection queued [%T]", id, i.Self), log.LEVEL_INFO)
	c.Schedule(state, id, &Task[T]{
		Lib:  "internal",
		Name: "collect",
		Fn: func(i *Item[T]) {
			i.Lg.Append(fmt.Sprintf("item %d collected  [%T]", id, i.Self), log.LEVEL_INFO)

			if c.onCollect != nil {
				c.onCollect(i)
			}
			i.Self = nil
			i.cleaned = true
		},
		Fail: func(i *Item[T]) {
			i.Self = nil
			i.cleaned = true
		},
	})
}

func (c *Collection[T]) Next() int {
	return len(c.items)
}
