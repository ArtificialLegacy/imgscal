package imageutil

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"image"
	"image/png"
	"io"

	golua "github.com/yuin/gopher-lua"
)

const ancillaryBit = 0b0010_0000 << 24
const magic = "\x89PNG\x0D\x0A\x1A\x0A"
const checksumLength = 4

const dataChunk uint32 = 0x69_73_63_4C // iscL
const dataKeyLength = 8

type PNGDataError string

func (e PNGDataError) Error() string {
	return string(e)
}

func NewPNGDataError(err string) PNGDataError {
	return PNGDataError(err)
}

func NewPNGDataErrorf(format string, args ...interface{}) PNGDataError {
	return PNGDataError(fmt.Sprintf(format, args...))
}

func PNGDataChunkWrite(key string, data string) []byte {
	db := []byte(data)
	ln := len(db) + dataKeyLength

	b := make([]byte, 12+ln) // 8 bytes for header, 4 bytes for crc32

	binary.BigEndian.PutUint32(b[0:4], uint32(ln))
	binary.BigEndian.PutUint32(b[4:8], dataChunk)

	copy(b[8:], key[:])
	copy(b[8+dataKeyLength:], db)

	crc := crc32.NewIEEE()
	crc.Write(b[4:8])
	crc.Write([]byte(key))
	crc.Write(db)
	binary.BigEndian.PutUint32(b[8+ln:], crc.Sum32())

	return b
}

func PNGDataChunkRead(p []byte) (*PNGDataChunk, error) {
	dataLn := binary.BigEndian.Uint32(p[0:4])
	head := binary.BigEndian.Uint32(p[4:8])

	if head != dataChunk {
		return nil, NewPNGDataErrorf("invalid chunk type %s", p[4:8])
	}

	key := p[8 : 8+dataKeyLength]
	data := p[8+dataKeyLength : 8+dataLn]

	crc := crc32.NewIEEE()
	crc.Write(p[4:8])
	crc.Write(key)
	crc.Write(data)
	sum := crc.Sum32()

	if sum != binary.BigEndian.Uint32(p[8+dataLn:]) {
		return nil, NewPNGDataErrorf("invalid checksum %d, expected %d", binary.BigEndian.Uint32(p[8+dataLn:]), sum)
	}

	return &PNGDataChunk{
		Key:  string(key),
		Data: string(data),
	}, nil
}

type PNGDataChunk struct {
	Key  string
	Data string
}

func NewPNGDataChunk(key string, data string) *PNGDataChunk {
	if len(key) < dataKeyLength {
		key = fmt.Sprintf("%*s", dataKeyLength, key)
	} else if len(key) > dataKeyLength {
		key = key[:dataKeyLength]
	}

	return &PNGDataChunk{Key: key, Data: data}
}

func DataChunkToTable(chunk *PNGDataChunk, state *golua.LState) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("key", golua.LString(chunk.Key))
	t.RawSetString("data", golua.LString(chunk.Data))

	return t
}

func TableToDataChunk(t *golua.LTable) *PNGDataChunk {
	key := t.RawGetString("key")
	data := t.RawGetString("data")

	return NewPNGDataChunk(string(key.(golua.LString)), string(data.(golua.LString)))
}

func PNGDataChunkEncode(w io.WriteSeeker, img image.Image, chunks []*PNGDataChunk) error {
	err := png.Encode(w, img)
	if err != nil {
		return err
	}

	_, err = w.Seek(-12, io.SeekCurrent)
	if err != nil {
		return err
	}

	for _, chunk := range chunks {
		dc := PNGDataChunkWrite(chunk.Key, chunk.Data)

		_, err = w.Write(dc)
		if err != nil {
			return err
		}
	}

	// IEND
	w.Write([]byte{
		0x0, 0x0, 0x0, 0x0, // length
		0x49, 0x45, 0x4E, 0x44, // type
		0xAE, 0x42, 0x60, 0x82, // crc32
	})

	return nil
}

func PNGDataChunkDecode(r io.Reader) (image.Image, []*PNGDataChunk, error) {
	cs := &PNGChunkStripper{
		Reader: r,
	}

	img, err := png.Decode(cs)
	if err != nil {
		return nil, nil, err
	}

	return img, cs.DataChunks, nil
}

type PNGChunkStripper struct {
	Reader     io.Reader
	DataChunks []*PNGDataChunk

	err     error
	ErrList []error
	magic   bool

	chunkType   uint32
	chunkLength uint32
	buffer      [8]byte
	start       int
	end         int
}

func (r *PNGChunkStripper) Read(p []byte) (int, error) {
	pass := false
	data := false

	for {
		if len(p) == 0 {
			return 0, nil
		}

		if r.err != nil {
			n := copy(p, r.buffer[r.start:r.end])
			r.start += n

			if r.start < r.end {
				return n, nil
			}

			return n, r.err
		}

		if r.start < r.end {
			n := copy(p, r.buffer[r.start:r.end])
			r.start += n

			return n, nil
		}

		if data {
			b := make([]byte, r.chunkLength+8)

			binary.BigEndian.PutUint32(b[0:4], r.chunkLength-4)
			binary.BigEndian.PutUint32(b[4:8], r.chunkType)

			_, err := io.ReadFull(r.Reader, b[8:])
			if err != nil {
				r.err = err
				r.ErrList = append(r.ErrList, err)
				continue
			}

			dc, err := PNGDataChunkRead(b)
			if err != nil {
				r.err = err
				r.ErrList = append(r.ErrList, err)
				continue
			}

			r.DataChunks = append(r.DataChunks, dc)

			r.chunkLength = 0
			data = false

			continue
		}

		for r.chunkLength > 0 {
			if uint32(len(p)) > r.chunkLength {
				p = p[:r.chunkLength]
			}

			n, err := r.Reader.Read(p)
			if err != nil {
				r.ErrList = append(r.ErrList, err)
				r.err = err
			}

			r.chunkLength -= uint32(n)

			if pass {
				continue
			}

			return n, err
		}

		end, err := io.ReadFull(r.Reader, r.buffer[:8])
		r.start = 0
		r.end = end

		if err != nil {
			r.err = err
			r.ErrList = append(r.ErrList, err)
			continue
		}

		if !r.magic && string(r.buffer[:8]) != magic {
			panic("cannot strip chunks: invalid png")
		}

		if r.magic {
			r.chunkLength = binary.BigEndian.Uint32(r.buffer[:4]) + checksumLength
			r.chunkType = binary.BigEndian.Uint32(r.buffer[4:])

			if r.chunkType == dataChunk {
				data = true
				r.start = r.end
			} else {
				pass = (r.chunkType & ancillaryBit) != 0
				if pass {
					r.start = r.end
				}
			}
		}

		r.magic = true
	}
}
