package image_util_test

import (
	"encoding/binary"
	"hash/crc32"
	"os"
	"testing"

	"github.com/ArtificialLegacy/imgscal/pkg/byteseeker"
	imageutil "github.com/ArtificialLegacy/imgscal/pkg/image_util"
)

func TestPNGDataKeyLength(t *testing.T) {
	c1 := imageutil.NewPNGDataChunk("testkey0", "Hello, World!")
	c2 := imageutil.NewPNGDataChunk("testkey11", "Hello, World!")
	c3 := imageutil.NewPNGDataChunk("testkey", "Hello, World!")

	if len(c1.Key) != 8 {
		t.Errorf("Expected key length of 8, got %d", len(c1.Key))
	}

	if len(c2.Key) != 8 {
		t.Errorf("Expected key length of 8, got %d", len(c2.Key))
	}

	if len(c3.Key) != 8 {
		t.Errorf("Expected key length of 8, got %d", len(c3.Key))
	}
}

func TestPNGDataWrite(t *testing.T) {
	data := "Hello, World!"
	key := "testkey0"
	dataLen := make([]byte, 4)
	binary.BigEndian.PutUint32(dataLen, uint32(len(data)+8))

	b := imageutil.PNGDataChunkWrite(key, data)

	if len(b) != 20+len(data) {
		t.Errorf("Expected length of %d, got %d", 20+len(data), len(b))
	}

	if string(b[:4]) != string(dataLen) {
		t.Errorf("Expected length of %d, got %s", dataLen, b[:4])
	}

	if string(b[4:8]) != "iscL" {
		t.Errorf("Expected chunk type of iscl, got %s", string(b[4:8]))
	}

	if string(b[8:16]) != key {
		t.Errorf("Expected key of %s, got %s", key, string(b[8:16]))
	}

	if string(b[16:16+len(data)]) != data {
		t.Errorf("Expected data of %s, got %s", data, string(b[16:16+len(data)]))
	}

	crc := crc32.NewIEEE()
	crc.Write(b[4:8])
	crc.Write([]byte(key))
	crc.Write([]byte(data))
	sum := crc.Sum32()

	if binary.BigEndian.Uint32(b[16+len(data):]) != sum {
		t.Errorf("Expected checksum of %d, got %d", sum, binary.BigEndian.Uint32(b[16+len(data):]))
	}
}

func TestPNGDataRead(t *testing.T) {
	data := "Hello, World!"
	key := "testkey0"
	dataLen := make([]byte, 4)
	binary.BigEndian.PutUint32(dataLen, uint32(len(data)))

	b := imageutil.PNGDataChunkWrite(key, data)

	d, err := imageutil.PNGDataChunkRead(b)
	if err != nil {
		t.Error(err)
	}

	if d.Key != key {
		t.Errorf("Expected key of %s, got %s", key, d.Key)
	}

	if d.Data != data {
		t.Errorf("Expected data of %s, got %s", data, d.Data)
	}
}

func TestPNGDataEncode(t *testing.T) {
	img := imageutil.NewImage(8, 8, imageutil.MODEL_NRGBA)

	data := []*imageutil.PNGDataChunk{
		imageutil.NewPNGDataChunk("testkey0", "Hello, World!"),
		imageutil.NewPNGDataChunk("testkey1", "Hello, World!"),
		imageutil.NewPNGDataChunk("testkey2", "Hello, World!"),
	}

	ws := byteseeker.NewByteSeeker(500, 50)

	err := imageutil.PNGDataChunkEncode(ws, img, data)
	if err != nil {
		t.Error(err)
	}

	os.WriteFile("test.png", ws.Bytes(), 0666)
}

func TestPNGDataDecode(t *testing.T) {
	r, err := os.ReadFile("test.png")
	if err != nil {
		t.Error(err)
	}

	ws := byteseeker.NewByteSeekerFromBytes(r, 50, true)

	img, data, err := imageutil.PNGDataChunkDecode(ws)
	if err != nil {
		t.Error(err)
	}

	if img.Bounds().Dx() != 8 {
		t.Errorf("Expected width of 8, got %d", img.Bounds().Dx())
	}

	if img.Bounds().Dy() != 8 {
		t.Errorf("Expected height of 8, got %d", img.Bounds().Dy())
	}

	if len(data) != 3 {
		t.Errorf("Expected 3 data chunks, got %d", len(data))
	}

	if data[0].Key != "testkey0" {
		t.Errorf("Expected key of testkey0, got %s", data[0].Key)
	}

	if data[1].Key != "testkey1" {
		t.Errorf("Expected key of testkey1, got %s", data[1].Key)
	}

	if data[2].Key != "testkey2" {
		t.Errorf("Expected key of testkey2, got %s", data[2].Key)
	}
}
