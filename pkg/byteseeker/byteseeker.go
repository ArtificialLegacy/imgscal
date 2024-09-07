package byteseeker

import (
	"errors"
	"math"
)

type ByteSeeker struct {
	buf  []byte
	pos  int
	grow int
}

func NewByteSeeker(size, grow int) *ByteSeeker {
	return &ByteSeeker{
		buf:  make([]byte, size),
		grow: grow,
	}
}

func NewByteSeekerFromBytes(b []byte, grow int, start bool) *ByteSeeker {
	pos := 0
	if !start {
		pos = len(b)
	}

	return &ByteSeeker{
		buf:  b,
		pos:  pos,
		grow: grow,
	}
}

func (ws *ByteSeeker) Write(p []byte) (n int, err error) {
	extra := ws.pos + len(p) - len(ws.buf)
	if extra > 0 {
		grow := int(math.Max(float64(extra), float64(ws.grow)))
		ws.buf = append(ws.buf, make([]byte, grow)...)
	}

	copy(ws.buf[ws.pos:], p)
	ws.pos += len(p)
	return len(p), nil
}

func (ws *ByteSeeker) Read(p []byte) (n int, err error) {
	if ws.pos >= len(ws.buf) {
		return 0, errors.New("EOF")
	}

	n = copy(p, ws.buf[ws.pos:])
	ws.pos += n
	return n, nil
}

func (ws *ByteSeeker) Seek(offset int64, whence int) (int64, error) {
	newPos, offs := 0, int(offset)
	switch whence {
	case 0:
		newPos = offs
	case 1:
		newPos = ws.pos + offs
	case 2:
		newPos = len(ws.buf) + offs
	}
	if newPos < 0 {
		return 0, errors.New("negative result pos")
	}
	ws.pos = newPos
	return int64(newPos), nil
}

func (ws *ByteSeeker) Bytes() []byte {
	return ws.buf[:ws.pos]
}
