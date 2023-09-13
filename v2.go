package vpk

import (
	"bufio"
	"encoding/binary"
	"io"
)

func openVPK_v2(fs FileReader, buffer []byte) (*vpk, error) {
	if _, err := fs.Read(buffer[:4*5]); err != nil {
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

	reader := bufio.NewReader(io.LimitReader(fs, int64(v.treeSize)))

	if err := treeReader(v, reader, buffer, v.addFile); err != nil {
		defer v.Close()
		return nil, err
	}

	// We should have read exactly .treeSize bytes and therefore hit EOF
	if _, err := reader.ReadByte(); err != io.EOF {
		defer v.Close()
		return nil, ErrWrongHeaderSize
	}

	return v, nil
}
