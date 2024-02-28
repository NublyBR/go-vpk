package vpk

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

func openVPK_v2(fs FileReader, buffer []byte) (*vpk, error) {
	// reset to beginning of file  (required for calculating the file checksum)
	_, err := fs.Seek(0, io.SeekStart)
	if err != nil {
		panic(err)
	}

	// split reader into two - a general reader and a md5 hasher
	file, fileMD5 := hashReader(fs)

	// re-read the file header (nothing required from it so header ignored)
	if _, err := file.Read(buffer[:4*2]); err != nil {
		return nil, err
	}

	// read parts of header that haven't already been read
	if _, err := file.Read(buffer[:4*5]); err != nil {
		return nil, err
	}

	v := &vpk{
		stream:     fs,
		version:    2,
		headerSize: 4 * 7,

		treeSize:              int32(binary.LittleEndian.Uint32(buffer[:4])),
		fileDataSectionSize:   int32(binary.LittleEndian.Uint32(buffer[4:8])),
		archiveMD5SectionSize: int32(binary.LittleEndian.Uint32(buffer[8:12])),
		otherMD5SectionSize:   int32(binary.LittleEndian.Uint32(buffer[12:16])),
		signatureSectionSize:  int32(binary.LittleEndian.Uint32(buffer[16:20])),

		pathMap: make(map[string]*entry),
	}

	// create tree reader and tree md5 hasher
	tree, treeMD5 := hashReader(io.LimitReader(file, int64(v.treeSize)))
	if err := treeReader(v, bufio.NewReader(tree), buffer, v.addFile); err != nil {
		defer v.Close()
		return nil, err
	}

	// We should have read exactly .treeSize bytes and therefore hit EOF
	if _, err := tree.Read(make([]byte, 1)); err != io.EOF {
		defer v.Close()
		return nil, ErrWrongHeaderSize
	}

	// read fileDataSection
	err = nullBufferedRead(file, buffer, int(v.fileDataSectionSize))
	if err != nil {
		defer v.Close()
		return nil, err
	}

	// read fileDataChecksum
	fdhReader, fileDataMD5 := hashReader(io.LimitReader(file, int64(v.archiveMD5SectionSize)))
	err = nullBufferedRead(fdhReader, buffer, int(v.archiveMD5SectionSize))
	if err != nil {
		defer v.Close()
		return nil, err
	}

	// read tree checksum, and fileDataHashes checksum
	if _, err := io.ReadFull(file, buffer[:32]); err != nil {
		defer v.Close()
		return nil, err
	}

	// the checksum value of the tree_section
	tcsExpected := buffer[:16]
	if !bytes.Equal(treeMD5.Sum(nil), tcsExpected) {
		defer v.Close()
		return nil, fmt.Errorf("mismatched tree checksum")
	}

	// the checksum of the fileData section
	fdhExpected := buffer[16:32]
	if !bytes.Equal(fileDataMD5.Sum(nil), fdhExpected) {
		defer v.Close()
		return nil, fmt.Errorf("mismatched file data checksum")
	}

	// calculate the file checksum value (minus the last 16 bytes of file which are the
	// checksum to compare against)
	fileMD5Value := fileMD5.Sum(nil)
	if _, err := io.ReadFull(file, buffer[:16]); err != nil {
		defer v.Close()
		return nil, err
	}

	// compare file checksum
	if !bytes.Equal(fileMD5Value, buffer[:16]) {
		defer v.Close()
		return nil, fmt.Errorf("mismatched file data checksum")
	}

	return v, nil
}
