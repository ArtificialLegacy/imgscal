package collection

import (
	"fmt"
	"image"
	"os"
	"runtime/debug"
	"sync"

	imageutil "github.com/ArtificialLegacy/imgscal/pkg/image_util"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/fogleman/gg"
	"github.com/skip2/go-qrcode"
)

const TASK_QUEUE_SIZE = 64

type CollectionType int

const (
	TYPE_TASK CollectionType = iota
	TYPE_IMAGE
	TYPE_FILE
	TYPE_CONTEXT
	TYPE_QR
)

var CollectionList = []CollectionType{
	TYPE_TASK,
	TYPE_IMAGE,
	TYPE_FILE,
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

type ItemFile struct {
	File *os.File
	Name string
}

func (file ItemFile) Identifier() CollectionType { return TYPE_FILE }

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

	cleaned bool
	collect bool
	waiting bool

	currTask *Task[T]

	failed bool
	Err    error

	TaskQueue chan *Task[T]
}

func NewItem[T ItemSelf](lg *log.Logger, fn func(i *Item[T])) *Item[T] {
	i := &Item[T]{
		Self: nil,
		Lg:   lg,

		cleaned:   false,
		collect:   false,
		failed:    false,
		waiting:   true,
		TaskQueue: make(chan *Task[T], TASK_QUEUE_SIZE),
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
			i.waiting = true
			i.cleaned = true
			i.Err = fmt.Errorf("%+v", p)
			if i.currTask != nil {
				if i.currTask.Fail != nil {
					i.currTask.Fail(i)
				}
			}
			for c := 0; len(i.TaskQueue) > 0; c++ {
				task := <-i.TaskQueue
				if task.Fail != nil {
					task.Fail(i)
				}
				i.Lg.Append(fmt.Sprintf("drained task %d from item [%T]", c, i.Self), log.LEVEL_WARN)
			}
		}
	}()

	for {
		i.waiting = true
		task := <-i.TaskQueue
		i.currTask = task
		i.waiting = false
		i.Lg.Append(fmt.Sprintf("%s.%s task called", task.Lib, task.Name), log.LEVEL_VERBOSE)
		task.Fn(i)
		i.Lg.Append(fmt.Sprintf("%s.%s task finished", task.Lib, task.Name), log.LEVEL_VERBOSE)

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

	Errs []error

	onCollect func(i *Item[T])
}

func NewCollection[T ItemSelf](lg *log.Logger) *Collection[T] {
	return &Collection[T]{
		items: []*Item[T]{},
		lg:    lg,
		Errs:  []error{},
	}
}

func (c *Collection[T]) OnCollect(fn func(i *Item[T])) *Collection[T] {
	c.onCollect = fn
	return c
}

func (c *Collection[T]) AddItem(lg *log.Logger) int {
	item := NewItem(lg, c.onCollect)
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

func (c *Collection[T]) Schedule(id int, tk *Task[T]) <-chan struct{} {
	wait := make(chan struct{}, 2)

	if id < 0 && id >= len(c.items) {
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
		item.Lg.Append(fmt.Sprintf("cannot schedule task for failed item: %d", id), log.LEVEL_WARN)
		task.Fail(item)
		wait <- struct{}{}
		return wait
	}

	item.Lg.Append(fmt.Sprintf("task scheduled for %d", id), log.LEVEL_VERBOSE)
	item.TaskQueue <- task

	return wait
}

func (c *Collection[T]) ScheduleAll(tk *Task[T]) <-chan struct{} {
	c.lg.Append(fmt.Sprintf("tasks scheduled for all items: [%T]", c), log.LEVEL_VERBOSE)

	wait := make(chan struct{}, 1)
	wg := sync.WaitGroup{}

	for _, i := range c.items {
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

		i.TaskQueue <- task
	}

	go func() {
		wg.Wait()
		wait <- struct{}{}
	}()

	return wait
}

func (c *Collection[T]) TaskCount() (int, bool) {
	count := 0

	errList := []error{}
	busy := false

	for _, i := range c.items {
		if i.Err != nil {
			errList = append(errList, i.Err)
		}
		if !i.waiting && !i.failed {
			busy = true
		}
		count += len(i.TaskQueue)
	}

	c.Errs = errList
	return count, busy
}

func (c *Collection[T]) CollectAll() {
	for id, i := range c.items {
		if i.collect {
			continue
		}

		i.Lg.Append(fmt.Sprintf("item %d collected  [%T]", id, i.Self), log.LEVEL_INFO)
		i.Lg.Close()

		if c.onCollect != nil {
			c.onCollect(i)
		}
		i.Self = nil
		i.cleaned = true
	}
}

func (c *Collection[T]) Collect(id int) {
	i := c.items[id]
	i.collect = true

	c.lg.Append(fmt.Sprintf("item %d collection queued [%T]", id, i.Self), log.LEVEL_INFO)
	c.Schedule(id, &Task[T]{
		Lib:  "internal",
		Name: "collect",
		Fn: func(i *Item[T]) {
			i.Lg.Append(fmt.Sprintf("item %d collected  [%T]", id, i.Self), log.LEVEL_INFO)
			i.Lg.Close()

			if c.onCollect != nil {
				c.onCollect(i)
			}
			i.Self = nil
			i.cleaned = true
		},
	})
}

func (c *Collection[T]) Next() int {
	return len(c.items)
}
