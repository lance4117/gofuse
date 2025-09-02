package fxerror

import (
	"errors"
)

var (
	ErrNil         = errors.New("error nil")
	ErrConfigRead  = errors.New("error reading fxconfig")
	ErrConfigLoad  = errors.New("error loading fxconfig")
	ErrNeedPointer = errors.New("error must be pointer")
)
