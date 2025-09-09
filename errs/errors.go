package errs

import (
	"errors"
)

var (
	ErrNil            = errors.New("error nil")
	ErrNotInitialized = errors.New("error not initialized")
	ErrConfigRead     = errors.New("error reading config")
	ErrConfigLoad     = errors.New("error loading config")
	ErrNeedPointer    = errors.New("error must be pointer")
)
