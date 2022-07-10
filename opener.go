package vpk

import (
	"encoding/binary"
	"fmt"
	"os"
)

func OpenStream(fs FileReader) (VPK, error) {
	// Buffer to be reused for reading the file data
	buffer := make([]byte, 64)

	if _, err := fs.Read(buffer[:4]); err != nil {
		return nil, err
	}

	// Verify if the file begins with the file signature `34 12 AA 55`
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

// Opens a single VPK file.
func OpenSingle(path string) (VPK, error) {
	fs, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return OpenStream(fs)
}

// Opens either a single VPK file or a VPK directory depending on the file name.
func OpenAny(path string) (VPK, error) {
	if reDirPath.MatchString(path) {
		return OpenDir(path)
	}
	return OpenSingle(path)
}
