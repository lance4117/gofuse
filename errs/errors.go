package errs

import (
	"errors"
)

var (
	ErrNil        = errors.New("error nil")
	ErrConfigRead = errors.New("error reading config")
	ErrConfigLoad = errors.New("error loading config")
)
