package image

import (
	"fmt"
	"image"
	"sync"
)

type Image struct {
	Mutex *sync.Mutex
	Img   *image.Image
	Name  string

	Ready   bool
	Cleaned bool
}

func NewImage(name string) *Image {
	return &Image{
		Mutex:   &sync.Mutex{},
		Img:     nil,
		Name:    name,
		Cleaned: false,
	}
}

type ImageCollection struct {
	images []*Image
}

func NewImageCollection() *ImageCollection {
	return &ImageCollection{
		images: []*Image{},
	}
}

func (ic *ImageCollection) Image(id int) (*Image, error) {
	img := ic.images[id]
	img.Mutex.Lock()

	if img.Cleaned {
		img.Mutex.Unlock()
		return nil, fmt.Errorf("attempting to get a cleaned image")
	}

	return img, nil
}

func (ic *ImageCollection) AddImage(name string) (*Image, int) {
	img := NewImage(name)
	id := len(ic.images)

	ic.images = append(ic.images, img)

	return img, id
}
