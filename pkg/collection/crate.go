package collection

import (
	"fmt"
	"image"

	imageutil "github.com/ArtificialLegacy/imgscal/pkg/image_util"
	"github.com/charmbracelet/lipgloss"
	"gitlab.com/microo8/blackcl"
)

type CrateItem[T any] struct {
	Self *T
}

type RefItem[T any] struct {
	Value T
}

type StyleItem struct {
	Style *lipgloss.Style
}

type CachedImageItem struct {
	Image image.Image
	Model imageutil.ColorModel
}

type ShaderItem struct {
	Device        *blackcl.Device
	BuffersImage  []*blackcl.Image
	BuffersVector []*blackcl.Vector
	BuffersBytes  []*blackcl.Bytes
	Kernels       []*blackcl.KernelCall
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
