package errs

import (
	"errors"
	"log"
)

var (
	ErrNil        = errors.New("error nil")
	ErrConfigRead = errors.New("error reading config")
)

func FatalErr(errType, err error) {
	log.Fatalf("%s: %s", errType, err)
}
