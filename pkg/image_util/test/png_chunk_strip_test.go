package image_util_test

import (
	"bytes"
	"image/png"
	"io"
	"os"
	"testing"

	imageutil "github.com/ArtificialLegacy/imgscal/pkg/image_util"
)

// needs files placed in this dir to pass, ideally files with known invalid checksums
func TestPNGChunkStrip(t *testing.T) {
	b, err := os.ReadFile("./Sprite_1.png")
	if err != nil {
		t.Error(err)
		return
	}

	strip := imageutil.PNGChunkStripper{
		Reader: bytes.NewReader(b),
	}

	img, err := png.Decode(io.Reader(&strip))
	if err != nil {
		t.Error(err)
		return
	}

	f, err := os.OpenFile("./output.png", os.O_CREATE, 0o666)
	if err != nil {
		t.Error(err)
		return
	}
	defer f.Close()

	err = png.Encode(f, img)
	if err != nil {
		t.Error(err)
		return
	}
}
