package vpk

import "errors"

var (
	ErrInvalidVPKVersion   = errors.New("invalid VPK version")
	ErrInvalidVPKSignature = errors.New("invalid VPK signature")
	ErrWrongHeaderSize     = errors.New("wrong header size")
	ErrInvalidArchiveIndex = errors.New("invalid archive index")
	ErrInvalidPath         = errors.New("invalid path")
)
