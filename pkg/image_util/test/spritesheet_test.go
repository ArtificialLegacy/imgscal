package image_util_test

import (
	"image"
	"image/color"
	"testing"

	imageutil "github.com/ArtificialLegacy/imgscal/pkg/image_util"
)

func createSpritesheet(colors []color.RGBA, width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, colors[x%len(colors)])
		}
	}

	return img
}

var c1 = color.RGBA{255, 0, 0, 255}
var c2 = color.RGBA{0, 255, 0, 255}
var c3 = color.RGBA{0, 0, 255, 255}
var c4 = color.RGBA{255, 255, 0, 255}

func TestSpritesheetToFrames(t *testing.T) {
	ss := createSpritesheet([]color.RGBA{
		c1, c2, c3, c4,
	}, 4, 4)

	frames := imageutil.SpritesheetToFrames(ss, true, 4, 1, 1, 4, 0, 0, 0, 0, 4, 0, 0)

	if len(frames) != 4 {
		t.Errorf("Expected 4 frames, got %d", len(frames))
	}

	for i, frame := range frames {
		if frame.Bounds().Dx() != 1 || frame.Bounds().Dy() != 1 {
			t.Errorf("Frame %d has incorrect dimensions: %v", i, frame.Bounds())
		}

		c := color.RGBA{}
		switch i {
		case 0:
			c = c1
		case 1:
			c = c2
		case 2:
			c = c3
		case 3:
			c = c4
		}

		frameColor := frame.At(0, 0).(color.RGBA)

		if frameColor != c {
			t.Errorf("Frame %d has incorrect color: %v, expected %v", i, frameColor, c)
		}
	}
}

func TestFramesToSpritesheet(t *testing.T) {
	frames := []image.Image{
		image.NewRGBA(image.Rect(0, 0, 1, 1)),
		image.NewRGBA(image.Rect(0, 0, 1, 1)),
		image.NewRGBA(image.Rect(0, 0, 1, 1)),
		image.NewRGBA(image.Rect(0, 0, 1, 1)),
	}

	frames[0].(*image.RGBA).Set(0, 0, c1)
	frames[1].(*image.RGBA).Set(0, 0, c2)
	frames[2].(*image.RGBA).Set(0, 0, c3)
	frames[3].(*image.RGBA).Set(0, 0, c4)

	ss := imageutil.FramesToSpritesheet(frames, imageutil.MODEL_RGBA, 4, 1, 1, 4, 0, 0, 0, 0, 0, 0, 0)

	if ss.Bounds().Dx() != 4 || ss.Bounds().Dy() != 1 {
		t.Errorf("Spritesheet has incorrect dimensions: %v", ss.Bounds())
	}

	for i, frame := range frames {
		frameColor := frame.At(0, 0).(color.RGBA)
		ssColor := ss.At(i, 0).(color.RGBA)

		if frameColor != ssColor {
			t.Errorf("Frame %d has incorrect color: %v, expected %v", i, ssColor, frameColor)
		}
	}
}
