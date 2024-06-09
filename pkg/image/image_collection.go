package image

import (
	"fmt"
	"image"
	"sync"
)

type Image struct {
	Mutex sync.Mutex
	Img   *image.Image

	Ready   bool
	Cleaned bool
}

func NewImage() *Image {
	return &Image{
		Mutex:   sync.Mutex{},
		Img:     nil,
		Ready:   false,
		Cleaned: false,
	}
}

type ImageCollection struct {
	images []*Image
}

func (ic *ImageCollection) Image(id int) (*Image, error) {
	img := ic.images[id]

	if img.Cleaned {
		return nil, fmt.Errorf("attempting to get a cleaned image")
	}

	if img.Ready {
		return nil, fmt.Errorf("attempting to get a non-ready image")
	}

	return img, nil
}

func (ic *ImageCollection) AddImage() (*Image, int) {
	img := NewImage()
	id := len(ic.images)

	ic.images = append(ic.images, img)

	return img, id
}
