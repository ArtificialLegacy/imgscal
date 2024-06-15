package collection

import (
	"fmt"
	"sync"

	"github.com/ArtificialLegacy/imgscal/pkg/log"
)

const TASK_QUEUE_SIZE = 64

type Item[T any] struct {
	Self *T
	Name string

	lg *log.Logger

	cleaned bool
	collect bool

	TaskQueue chan *Task[T]
}

func NewItem[T any](name string, lg *log.Logger) *Item[T] {
	i := &Item[T]{
		Self: nil,
		Name: name,
		lg:   lg,

		cleaned:   false,
		collect:   false,
		TaskQueue: make(chan *Task[T], TASK_QUEUE_SIZE),
	}

	go i.process()

	return i
}

func (i *Item[T]) process() {
	for {
		task := <-i.TaskQueue
		i.lg.Append(fmt.Sprintf("%s.%s task called", task.Lib, task.Name), log.LEVEL_INFO)
		task.Fn(i)
		i.lg.Append(fmt.Sprintf("%s.%s task finished", task.Lib, task.Name), log.LEVEL_INFO)

		if i.cleaned {
			i.lg.Append(fmt.Sprintf("item %s cleaned", i.Name), log.LEVEL_INFO)
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

func (c *Collection[T]) AddItem(name string) int {
	item := NewItem[T](name, c.lg)
	id := len(c.items)

	c.items = append(c.items, item)

	return id
}

func (c *Collection[T]) Schedule(id int, tk *Task[T]) <-chan struct{} {
	c.lg.Append(fmt.Sprintf("task scheduled for %d", id), log.LEVEL_INFO)

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

func (c *Collection[T]) CollectAll() {
	wg := sync.WaitGroup{}

	for id, i := range c.items {
		if i.collect {
			continue
		}

		wg.Add(1)
		idHere := id
		i.collect = true

		c.lg.Append(fmt.Sprintf("item %d collection queued [%T]", idHere, i.Self), log.LEVEL_INFO)
		c.Schedule(id, &Task[T]{
			Lib:  "internal",
			Name: "collect_all",
			Fn: func(i *Item[T]) {
				c.lg.Append(fmt.Sprintf("item %d collected  [%T]", idHere, i.Self), log.LEVEL_INFO)

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
}

func (c *Collection[T]) Collect(id int) {
	i := c.items[id]
	i.collect = true

	c.lg.Append(fmt.Sprintf("item %d collection queued [%T]", id, i.Self), log.LEVEL_INFO)
	c.Schedule(id, &Task[T]{
		Lib:  "internal",
		Name: "collect",
		Fn: func(i *Item[T]) {
			c.lg.Append(fmt.Sprintf("item %d collected  [%T]", id, i.Self), log.LEVEL_INFO)

			if c.onCollect != nil {
				c.onCollect(i)
			}
			i.Self = nil
			i.cleaned = true
		},
	})
}
