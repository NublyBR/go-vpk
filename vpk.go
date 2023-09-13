package vpk

import (
	"os"
)

type vpk struct {
	stream     FileReader
	indexes    []FileReader
	version    int
	headerSize int

	// The size, in bytes, of the directory tree
	treeSize int32

	// How many bytes of file content are stored in this VPK file (0 in CSGO)
	fileDataSectionSize int32

	// The size, in bytes, of the section containing MD5 checksums for external archive content
	archiveMD5SectionSize int32

	// The size, in bytes, of the section containing MD5 checksums for content in this file (should always be 48)
	otherMD5SectionSize int32

	// The size, in bytes, of the section containing the public key and signature. This is either 0 (CSGO & The Ship) or 296 (HL2, HL2:DM, HL2:EP1, HL2:EP2, HL2:LC, TF2, DOD:S & CS:S)
	signatureSectionSize int32

	files   []*entry
	pathMap map[string]*entry
}

func (v *vpk) addFile(e *entry) {
	v.files = append(v.files, e)
	v.pathMap[e.Filename()] = e
	e.parent = v
}

func (v *vpk) Open(path string) (FileReader, error) {
	entry, ok := v.pathMap[path]
	if !ok {
		return nil, os.ErrNotExist
	}

	return entry.Open()
}

func (v *vpk) Entries() []Entry {
	ret := make([]Entry, len(v.files))
	for i, f := range v.files {
		ret[i] = f
	}
	return ret
}

func (v *vpk) Find(path string) (Entry, bool) {
	entry, ok := v.pathMap[path]

	return entry, ok
}

func (v *vpk) Close() error {
	v.stream.Close()

	for _, f := range v.indexes {
		f.Close()
	}

	return nil
}
