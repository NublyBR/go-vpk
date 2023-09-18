package vpk

import (
	"crypto/md5"
	"hash"
	"io"
	"os"
)

type entryReader struct {
	fs FileReader

	closed bool

	offset   int64
	size     int64
	position int64
}

func (e *entryReader) Read(p []byte) (int, error) {
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

func (e *entryReader) ReadAt(p []byte, off int64) (int, error) {
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

func (e *entryReader) Seek(offset int64, whence int) (int64, error) {
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

func (e *entryReader) Close() error {
	if e.closed {
		return os.ErrClosed
	}
	e.closed = true
	return nil
}

func hashReader(r io.Reader) (io.Reader, hash.Hash) {
	hasher := md5.New()
	cReader := io.TeeReader(r, hasher)
	return cReader, hasher
}

// nullBufferedRead provides a way to "read away" the provided number of bytes from a
// reader using a pre-defined buffer (to save on memory). The bytes that are read will
// be discarded.
func nullBufferedRead(r io.Reader, buffer []byte, size int) error {
	count := 0

	for count < size {
		limit := len(buffer)

		if size-count < len(buffer) {
			limit = size - count
		}

		_, err := io.ReadFull(r, buffer[:limit])
		if err != nil {
			return err
		}

		count += limit
	}

	// reset buffer to leave out any remnants of previous reads
	buffer = buffer[:0]

	return nil
}
