package fserror

import (
	"errors"
)

var (
	ErrNil         = errors.New("error nil")
	ErrConfigRead  = errors.New("error reading fsconfig")
	ErrConfigLoad  = errors.New("error loading fsconfig")
	ErrNeedPointer = errors.New("error must be pointer")
)
