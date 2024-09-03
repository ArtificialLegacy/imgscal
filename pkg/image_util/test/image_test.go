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
