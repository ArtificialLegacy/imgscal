package writeseeker

import "errors"

type WriteSeeker struct {
	buf []byte
	pos int
}

func NewWriteSeeker(size int) *WriteSeeker {
	return &WriteSeeker{
		buf: make([]byte, size),
	}
}

func (ws *WriteSeeker) Write(p []byte) (n int, err error) {
	extra := ws.pos + len(p) - len(ws.buf)
	if extra > 0 {
		ws.buf = append(ws.buf, make([]byte, extra)...)
	}

	copy(ws.buf[ws.pos:], p)
	ws.pos += len(p)
	return len(p), nil
}

func (ws *WriteSeeker) Seek(offset int64, whence int) (int64, error) {
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

func (ws *WriteSeeker) Bytes() []byte {
	return ws.buf
}
