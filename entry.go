package vpk

import "strings"

type entry_impl struct {
	parent *vpk_impl

	ext  string
	path string
	file string

	// A 32bit crc of the file's data.
	crc uint32

	// The number of bytes contained in the index file.
	preloadBytes uint16

	// A zero based index of the archive this file's data is contained in.
	// If 0x7fff, the data follows the directory.
	archiveIndex uint16

	// If ArchiveIndex is 0x7fff, the offset of the file data relative to the end of the directory (see the header for more details).
	// Otherwise, the offset of the data from the start of the specified archive.
	entryOffset uint32

	// If zero, the entire file is stored in the preload data.
	// Otherwise, the number of bytes stored starting at EntryOffset.
	entryLength uint32
}

func (e *entry_impl) Filename() string {
	var parts []string

	if e.path != " " {
		parts = append(parts, e.path)
	}

	if e.file != " " {
		if e.ext != " " {
			parts = append(parts, e.file+"."+e.ext)
		} else {
			parts = append(parts, e.file)
		}
	} else if e.ext != " " {
		parts = append(parts, "."+e.ext)
	}

	return strings.Join(parts, "/")
}

func (e *entry_impl) Basename() string {
	if e.file != " " {
		if e.ext != " " {
			return e.file + "." + e.ext
		} else {
			return e.file
		}
	} else if e.ext != " " {
		return "." + e.ext
	}
	return ""
}

func (e *entry_impl) Path() string {
	if e.path == " " {
		return ""
	}
	return e.path
}

func (e *entry_impl) Length() uint32 {
	return e.entryLength
}

func (e *entry_impl) Open() (FileReader, error) {
	if e.archiveIndex == 0x7fff {
		return &entryReader{
			fs:     e.parent.stream,
			offset: int64(e.parent.headerSize) + int64(e.parent.treeSize) + int64(e.entryOffset),
			size:   int64(e.entryLength),
		}, nil
	}

	if e.archiveIndex >= uint16(len(e.parent.indexes)) {
		return nil, ErrInvalidArchiveIndex
	}

	return &entryReader{
		fs:     e.parent.indexes[e.archiveIndex],
		offset: int64(e.entryOffset),
		size:   int64(e.entryLength),
	}, nil
}
