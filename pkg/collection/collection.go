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

	failed bool

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
			for len(i.TaskQueue) > 0 {
				task := <-i.TaskQueue
				if task.Lib == "internal" {
					task.Fn(i)
				}
			}
		}
	}()

	for {
		task := <-i.TaskQueue
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
	Fn func(i *Item[T])

	Lib  string
	Name string
}

type Collection[T any] struct {
	items []*Item[T]
	lg    *log.Logger

	onCollect func(i *Item[T])
}

func NewCollection[T any](lg *log.Logger) *Collection[T] {
	return &Collection[T]{
		items: []*Item[T]{},
		lg:    lg,
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
	}

	item := c.items[id]
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

func (c *Collection[T]) CollectAll() error {
	wg := sync.WaitGroup{}

	var err error

	for id, i := range c.items {
		wg.Add(1)
		idHere := id
		iHere := i
		iHere.collect = true

		iHere.Lg.Append(fmt.Sprintf("item %d collection queued [%T]", idHere, i.Self), log.LEVEL_INFO)
		c.Schedule(idHere, &Task[T]{
			Lib:  "internal",
			Name: "collect_all",
			Fn: func(i *Item[T]) {
				if i.cleaned && i.failed {
					wg.Done()
					return
				}

				i.Lg.Append(fmt.Sprintf("item %d collected  [%T]", idHere, i.Self), log.LEVEL_INFO)

				if i.failed {
					err = fmt.Errorf(i.Lg.Append(fmt.Sprintf("item %d was marked as failed  [%T]", idHere, i.Self), log.LEVEL_ERROR))
				}
				i.Lg.Close()

				if c.onCollect != nil {
					c.onCollect(i)
				}
				i.Self = nil
				i.cleaned = true

				wg.Done()
			},
		})
	}

	wg.Wait()
	c.lg.Append("all items collected", log.LEVEL_INFO)

	return err
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
