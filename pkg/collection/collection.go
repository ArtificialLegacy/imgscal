package collection

import (
	"fmt"
	"sync"

	"github.com/ArtificialLegacy/imgscal/pkg/log"
)

const TASK_QUEUE_SIZE = 64

type CollectionType int

const (
	TYPE_TASK CollectionType = iota
	TYPE_IMAGE
	TYPE_FILE
)

type Item[T any] struct {
	Self *T
	Name string

	Lg *log.Logger

	cleaned bool
	collect bool
	waiting bool

	failed bool
	Err    error

	TaskQueue chan *Task[T]
}

func NewItem[T any](name string, lg *log.Logger, fn func(i *Item[T])) *Item[T] {
	i := &Item[T]{
		Self: nil,
		Name: name,
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
			i.Lg.Append("recovered from panic within collection item.", log.LEVEL_ERROR)
			if fn != nil {
				fn(i)
			}
			i.failed = true
			i.waiting = true
			i.cleaned = true
			i.Err = fmt.Errorf("%+v", p)
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
		i.waiting = false
		i.Lg.Append(fmt.Sprintf("%s.%s task called", task.Lib, task.Name), log.LEVEL_INFO)
		task.Fn(i)
		i.Lg.Append(fmt.Sprintf("%s.%s task finished", task.Lib, task.Name), log.LEVEL_INFO)

		if i.cleaned {
			i.Lg.Append(fmt.Sprintf("item %s cleaned", i.Name), log.LEVEL_INFO)
			break
		}
	}
}

type Task[T any] struct {
	Fn   func(i *Item[T])
	Fail func(i *Item[T])

	Lib  string
	Name string
}

type Collection[T any] struct {
	items []*Item[T]
	lg    *log.Logger

	Errs []error

	onCollect func(i *Item[T])
}

func NewCollection[T any](lg *log.Logger) *Collection[T] {
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

func (c *Collection[T]) AddItem(name string, lg *log.Logger) int {
	item := NewItem(name, lg, c.onCollect)
	id := len(c.items)

	c.items = append(c.items, item)

	return id
}

func (c *Collection[T]) Schedule(id int, tk *Task[T]) <-chan struct{} {
	wait := make(chan struct{}, 1)

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
		wait <- struct{}{}
		return wait
	}

	item.Lg.Append(fmt.Sprintf("task scheduled for %d", id), log.LEVEL_INFO)
	item.TaskQueue <- task

	return wait
}

func (c *Collection[T]) ScheduleAll(tk *Task[T]) <-chan struct{} {
	c.lg.Append(fmt.Sprintf("tasks scheduled for all items: [%T]", c), log.LEVEL_INFO)

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
