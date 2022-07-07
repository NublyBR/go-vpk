package vpk

import (
	"encoding/binary"
	"fmt"
	"os"
)

func OpenVPK(fs FileReader) (VPK, error) {
	buffer := make([]byte, 4096)

	if _, err := fs.Read(buffer[:4]); err != nil {
		return nil, err
	}

	if binary.LittleEndian.Uint32(buffer[:4]) != 0x55aa1234 {
		return nil, ErrInvalidVPKSignature
	}

	if _, err := fs.Read(buffer[:4]); err != nil {
		return nil, err
	}

	var version uint32 = binary.LittleEndian.Uint32(buffer[:4])

	switch version {
	case 1:
		return openVPK_v1(fs, buffer)
	case 2:
		return openVPK_v2(fs, buffer)
	}

	return nil, fmt.Errorf("unknown VPK version %d", version)
}

func OpenSingle(path string) (VPK, error) {
	fs, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return OpenVPK(fs)
}
