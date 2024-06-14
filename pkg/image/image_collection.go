package image

import (
	"fmt"
	"image"
	"sync"

	"github.com/ArtificialLegacy/imgscal/pkg/log"
)

type ImageCollection struct {
	images []*Image
	lg     *log.Logger
}

func NewImageCollection(lg *log.Logger) *ImageCollection {
	return &ImageCollection{
		images: []*Image{},
		lg:     lg,
	}
}

func (ic *ImageCollection) Image(id int) (*Image, error) {
	img := ic.images[id]

	return img, nil
}

func (ic *ImageCollection) Schedule(id int, task *ImageTask) {
	ic.lg.Append(fmt.Sprintf("task scheduled for %d", id), log.LEVEL_INFO)

	img := ic.images[id]
	img.TaskQueue <- task
}

func (ic *ImageCollection) AddImage(name string) int {
	img := NewImage(name, ic.lg)
	id := len(ic.images)

	ic.images = append(ic.images, img)

	return id
}

func (ic *ImageCollection) Collect() {
	wg := sync.WaitGroup{}

	for id := range ic.images {
		wg.Add(1)
		idHere := id
		ic.images[idHere].collect = true

		ic.lg.Append(fmt.Sprintf("image %d collection queued", idHere), log.LEVEL_INFO)
		ic.Schedule(id, &ImageTask{
			Lib:  "internal",
			Name: "collect",
			Fn: func(i *Image) {
				ic.lg.Append(fmt.Sprintf("image %d collected", idHere), log.LEVEL_INFO)
				i.Img = nil
				i.cleaned = true
				wg.Done()
			},
		})
	}

	wg.Wait()
	ic.lg.Append("all images collected", log.LEVEL_INFO)
}

func (ic *ImageCollection) CollectImage(id int) {
	ic.images[id].collect = true

	ic.lg.Append(fmt.Sprintf("image %d collection queued", id), log.LEVEL_INFO)
	ic.Schedule(id, &ImageTask{
		Lib:  "internal",
		Name: "collect_image",
		Fn: func(i *Image) {
			ic.lg.Append(fmt.Sprintf("image %d collected", id), log.LEVEL_INFO)
			i.Img = nil
		},
	})
}

type ImageTask struct {
	Lib  string
	Name string
	Fn   func(i *Image)
}

type Image struct {
	Img  image.Image
	Name string

	Lg *log.Logger

	collect bool
	cleaned bool

	TaskQueue chan *ImageTask
}

func NewImage(name string, lg *log.Logger) *Image {
	i := &Image{
		Img:       nil,
		Name:      name,
		Lg:        lg,
		TaskQueue: make(chan *ImageTask, 32),
		collect:   false,
	}

	go i.process()

	return i
}

func (i *Image) process() {
	for {
		task := <-i.TaskQueue
		i.Lg.Append(fmt.Sprintf("%s.%s task called", task.Lib, task.Name), log.LEVEL_INFO)
		task.Fn(i)
		i.Lg.Append(fmt.Sprintf("%s.%s task finished", task.Lib, task.Name), log.LEVEL_INFO)

		if i.cleaned {
			i.Lg.Append(fmt.Sprintf("image %s cleaned", i.Name), log.LEVEL_INFO)
			break
		}
	}
}
