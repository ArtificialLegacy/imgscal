package test

import (
	"testing"

	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
)

func TestCollection(t *testing.T) {
	lg := log.NewLoggerEmpty()
	c := collection.NewCollection[string](&lg)

	id := c.AddItem("test_item", &lg)
	value := ""

	c.Schedule(id, &collection.Task[string]{
		Fn: func(i *collection.Item[string]) {
			str := "test"
			i.Self = &str
		},
	})

	<-c.Schedule(id, &collection.Task[string]{
		Fn: func(i *collection.Item[string]) {
			value = *i.Self
		},
	})

	if value != "test" {
		t.Errorf("got wrong item after task run, expected=test_item got=%s", value)
	}

	for c, b := c.TaskCount(); b || c > 0; {
	}
}
