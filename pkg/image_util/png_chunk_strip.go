package imageutil

import (
	"encoding/binary"
	"io"
)

const ancillaryBit = 0b0010_0000 << 24
const magic = "\x89PNG\x0D\x0A\x1A\x0A"
const checksumLength = 4

type PNGChunkStripper struct {
	Reader io.Reader

	err     error
	ErrList []error
	magic   bool

	chunkLength uint32
	buffer      [8]byte
	start       int
	end         int
}

func (r *PNGChunkStripper) Read(p []byte) (int, error) {
	pass := false

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
			chunkType := binary.BigEndian.Uint32(r.buffer[4:])

			pass = (chunkType & ancillaryBit) != 0
			if pass {
				r.start = r.end
			}
		}

		r.magic = true
	}
}
