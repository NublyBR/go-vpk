package vpk

import (
	"bufio"
	"encoding/binary"
	"io"
)

func openVPK_v2(fs FileReader, buffer []byte) (*vpk_impl, error) {
	if _, err := fs.Read(buffer[:4*5]); err != nil {
		return nil, err
	}

	vpk := &vpk_impl{
		stream:  fs,
		version: 2,

		treeSize:              int32(binary.LittleEndian.Uint32(buffer[:4])),
		fileDataSectionSize:   int32(binary.LittleEndian.Uint32(buffer[4:8])),
		archiveMD5SectionSize: int32(binary.LittleEndian.Uint32(buffer[8:12])),
		otherMD5SectionSize:   int32(binary.LittleEndian.Uint32(buffer[12:16])),
		signatureSectionSize:  int32(binary.LittleEndian.Uint32(buffer[16:20])),

		pathMap: make(map[string]*entry_impl),
	}

	reader := bufio.NewReader(io.LimitReader(fs, int64(vpk.treeSize)))

	if err := treeReader(vpk, reader, buffer, vpk.addFile); err != nil {
		defer vpk.Close()
		return nil, err
	}

	// We should have read exactly .treeSize bytes and therefore hit EOF
	if _, err := reader.ReadByte(); err != io.EOF {
		defer vpk.Close()
		return nil, ErrWrongHeaderSize
	}

	return vpk, nil
}
