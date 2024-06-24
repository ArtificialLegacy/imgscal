package test

import (
	"testing"

	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
)

type ItemString struct {
	Value string
}

func (img ItemString) Identifier() string { return "ItemString" }

func TestCollection(t *testing.T) {
	lg := log.NewLoggerEmpty()
	c := collection.NewCollection[ItemString](&lg)

	id := c.AddItem(&lg)
	value := ""

	const expected = "test"

	c.Schedule(id, &collection.Task[ItemString]{
		Lib:  "test",
		Name: "first",
		Fn: func(i *collection.Item[ItemString]) {
			i.Self = &ItemString{
				Value: expected,
			}
		},
	})

	<-c.Schedule(id, &collection.Task[ItemString]{
		Lib:  "test",
		Name: "second",
		Fn: func(i *collection.Item[ItemString]) {
			value = i.Self.Value
		},
	})

	if value != expected {
		t.Errorf("got wrong item after task run, expected=%s got=%s", expected, value)
	}

	for c, b := c.TaskCount(); b || c > 0; {
	}
}
