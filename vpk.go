package vpk

import (
	"os"
)

type vpk_impl struct {
	stream  FileReader
	indexes []FileReader
	version int

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

	files   []*entry_impl
	pathMap map[string]*entry_impl
}

func (v *vpk_impl) addFile(entry *entry_impl) {
	v.files = append(v.files, entry)
	v.pathMap[entry.Filename()] = entry
	entry.parent = v
}

func (v *vpk_impl) Open(path string) (FileReader, error) {
	entry, ok := v.pathMap[path]
	if !ok {
		return nil, os.ErrNotExist
	}

	return entry.Open()
}

func (v *vpk_impl) Entries() []Entry {
	ret := make([]Entry, len(v.files))
	for i, f := range v.files {
		ret[i] = f
	}
	return ret
}

func (v *vpk_impl) Find(path string) (Entry, bool) {
	entry, ok := v.pathMap[path]

	return entry, ok
}

func (v *vpk_impl) Close() error {
	v.stream.Close()

	for _, f := range v.indexes {
		f.Close()
	}

	return nil
}
