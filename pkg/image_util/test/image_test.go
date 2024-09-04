package image_util_test

import (
	"image"
	"testing"

	imageutil "github.com/ArtificialLegacy/imgscal/pkg/image_util"
)

func TestImageCompare(t *testing.T) {
	img1 := image.NewRGBA(image.Rect(0, 0, 100, 100))
	img2 := image.NewRGBA(image.Rect(10, 10, 110, 110))
	img3 := image.NewNRGBA(image.Rect(0, 0, 100, 100))
	img4 := image.NewRGBA(image.Rect(0, 0, 99, 99))

	if !imageutil.ImageCompare(img1, img2) {
		t.Error("Images should be equal")
	}

	if !imageutil.ImageCompare(img1, img3) {
		t.Error("Images should be equal")
	}

	if imageutil.ImageCompare(img1, img4) {
		t.Error("Images should not be equal")
	}
}

func Test8BitTo16Bit(t *testing.T) {
	c8 := uint8(0x80)
	c16 := imageutil.Color8BitTo16Bit(c8)

	if c16 != 0x8080 {
		t.Errorf("Expected 0x8000, got %x", c16)
	}
}

func Test16BitTo8Bit(t *testing.T) {
	c16 := uint16(0x8080)
	c8 := imageutil.Color16BitTo8Bit(c16)

	if c8 != 0x80 {
		t.Errorf("Expected 0x80, got %x", c8)
	}
}
