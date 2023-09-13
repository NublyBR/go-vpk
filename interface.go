package vpk

import "io"

type FileReader interface {
	io.Reader
	io.ReaderAt
	io.Seeker
	io.Closer
}

type VPK interface {
	// Opens the entry at the given path
	Open(path string) (FileReader, error)

	// Find the entry at the given path
	Find(path string) (Entry, bool)

	// All entries in the VPK
	Entries() []Entry

	// Closes the VPK
	Close() error
}

type Entry interface {
	// Filename of VPK entry
	Filename() string

	// Filename without path
	Basename() string

	// Path of VPK entry
	Path() string

	// Length of VPK entry
	Length() uint32

	// Is filename and path safe to extract in linux/unix?
	FilenameSafeUnix() bool

	// Is filename and path safe to extract in windows?
	FilenameSafeWindows() bool

	// Open VPK entry for reading
	Open() (FileReader, error)

	// CRC of VPK entry
	CRC() uint32
}
