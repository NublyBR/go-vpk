package vpk

import (
	"io"
	"os"
)

type EntryReader struct {
	fs FileReader

	closed bool

	offset   int64
	size     int64
	position int64
}

func (e *EntryReader) Read(p []byte) (int, error) {
	if e.closed {
		return 0, os.ErrClosed
	}
	if e.position >= e.size {
		return 0, io.EOF
	}

	want := int64(len(p))

	if e.position+want > e.size {
		want = e.size - e.position
	}

	n, err := e.fs.ReadAt(p[:want], e.offset+e.position)
	e.position += int64(n)

	return n, err
}

func (e *EntryReader) ReadAt(p []byte, off int64) (int, error) {
	if e.closed {
		return 0, os.ErrClosed
	}
	if off < 0 || off >= e.size {
		return 0, io.EOF
	}

	want := int64(len(p))

	if off+want > e.size {
		want = e.size - off
	}

	n, err := e.fs.ReadAt(p[:want], e.offset+off)

	return n, err
}

func (e *EntryReader) Seek(offset int64, whence int) (int64, error) {
	if e.closed {
		return 0, os.ErrClosed
	}

	switch whence {
	case io.SeekStart:
		e.position = offset
	case io.SeekCurrent:
		e.position += offset
	case io.SeekEnd:
		e.position = e.size + offset
	}

	if e.position < 0 {
		e.position = 0
	} else if e.position > e.size {
		e.position = e.size
	}

	return e.position, nil
}

func (e *EntryReader) Close() error {
	if e.closed {
		return os.ErrClosed
	}
	e.closed = true
	return nil
}
