package test

import (
	"sync"
	"testing"
	"time"

	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	golua "github.com/yuin/gopher-lua"
)

type ItemString struct {
	Value string
}

func (img ItemString) Identifier() collection.CollectionType { return collection.TYPE_IMAGE }

const TYPE_STRING collection.CollectionType = 999

func TestCollection(t *testing.T) {
	lg := log.NewLoggerEmpty()
	wg := &sync.WaitGroup{}
	c := collection.NewCollection[ItemString](&lg, wg, TYPE_STRING)
	state := golua.NewState(golua.Options{
		SkipOpenLibs: true,
	})
	collection.CreateContext(state)

	id := c.AddItem(&lg)
	value := ""

	const expected = "test"

	c.Schedule(state, id, &collection.Task[ItemString]{
		Lib:  "test",
		Name: "first",
		Fn: func(i *collection.Item[ItemString]) {
			i.Self = &ItemString{
				Value: expected,
			}
		},
	})

	<-c.Schedule(state, id, &collection.Task[ItemString]{
		Lib:  "test",
		Name: "second",
		Fn: func(i *collection.Item[ItemString]) {
			value = i.Self.Value
		},
	})

	if value != expected {
		t.Errorf("got wrong item after task run, expected=%s got=%s", expected, value)
	}

	wg.Wait()
}

func TestCollectionNested(t *testing.T) {
	timeout := time.After(5 * time.Second)
	done := make(chan struct{})

	const expected = "test_value"
	value := ""

	go func() {
		lg := log.NewLoggerEmpty()
		wg := &sync.WaitGroup{}
		c := collection.NewCollection[ItemString](&lg, wg, TYPE_STRING)
		state := golua.NewState()
		collection.CreateContext(state)

		id := c.ScheduleAdd(state, "test", &lg, "test", "add", func(i *collection.Item[ItemString]) {
			i.Self = &ItemString{
				Value: expected,
			}
		})

		c.Schedule(state, id, &collection.Task[ItemString]{
			Lib:  "test",
			Name: "first",
			Fn: func(i *collection.Item[ItemString]) {
				inner := collection.NewThread(state, id, TYPE_STRING)
				<-c.Schedule(inner, id, &collection.Task[ItemString]{
					Lib:  "test",
					Name: "second",
					Fn: func(i *collection.Item[ItemString]) {
						value = i.Self.Value
					},
				})
			},
		})

		wg.Wait()
		done <- struct{}{}
	}()

	select {
	case <-done:
		if value != expected {
			t.Errorf("got wrong item after task run, expected=%s got=%s", expected, value)
		}
	case <-timeout:
		t.Fatal("test timed out")
	}
}
