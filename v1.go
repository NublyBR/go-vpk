package vpk

import (
	"bufio"
	"encoding/binary"
	"io"
)

func openVPK_v1(fs FileReader, buffer []byte) (*vpk_impl, error) {
	if _, err := fs.Read(buffer[:4]); err != nil {
		return nil, err
	}

	vpk := &vpk_impl{
		stream:  fs,
		version: 1,

		treeSize: int32(binary.LittleEndian.Uint32(buffer[:4])),

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
