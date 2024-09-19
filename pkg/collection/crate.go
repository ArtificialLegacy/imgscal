package collection

import (
	"fmt"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	golua "github.com/yuin/gopher-lua"
)

type CrateItem[T any] struct {
	Self *T
}

type RefItem[T any] struct {
	Value T
}

type TeaItem struct {
	LuaModel *golua.LTable

	FnInit   *golua.LFunction
	FnUpdate *golua.LFunction
	FnView   *golua.LFunction

	Msg  *tea.Msg
	Cmds []tea.Cmd

	Spinners    map[int]*spinner.Model
	TextAreas   []*textarea.Model
	TextInputs  []*textinput.Model
	Cursors     []*cursor.Model
	FilePickers []*filepicker.Model
}

type Crate[T any] struct {
	clean func(i *CrateItem[T])
	Items []*CrateItem[T]
}

func NewCrate[T any, C CrateItem[T]]() *Crate[T] {
	return &Crate[T]{
		Items: []*CrateItem[T]{},
		clean: nil,
	}
}

func (c *Crate[T]) OnClean(clean func(i *CrateItem[T])) *Crate[T] {
	c.clean = clean
	return c
}

func (c *Crate[T]) CleanAll() {
	for i, ic := range c.Items {
		if c.clean != nil {
			c.clean(ic)
		}
		c.Items[i] = nil
	}
}

func (c *Crate[T]) Clean(index int) {
	if c.clean != nil {
		c.clean(c.Items[index])
	}

	c.Items[index] = nil
}

func (c *Crate[T]) Add(i *T) int {
	c.Items = append(c.Items, &CrateItem[T]{
		Self: i,
	})

	return len(c.Items) - 1
}

func (c *Crate[T]) Item(id int) (*T, error) {
	if id < 0 {
		return nil, fmt.Errorf("id out of range: %d < 0", id)
	}
	if id >= len(c.Items) {
		return nil, fmt.Errorf("id out of range: %d >= %d", id, len(c.Items))
	}

	i := c.Items[id]
	return i.Self, nil
}

func (c *Crate[T]) ItemFull(id int) (*CrateItem[T], error) {
	if id < 0 {
		return nil, fmt.Errorf("id out of range: %d < 0", id)
	}
	if id >= len(c.Items) {
		return nil, fmt.Errorf("id out of range: %d >= %d", id, len(c.Items))
	}

	i := c.Items[id]
	return i, nil
}
