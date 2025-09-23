package errs

import (
	"errors"
)

var (
	ErrNil                      = errors.New("error nil")
	ErrFileReaderNotInitialized = errors.New("error file reader not initialized")
	ErrFileWriteNotInitialized  = errors.New("error file writer not initialized")
	ErrConfigRead               = errors.New("error reading config")
	ErrConfigLoad               = errors.New("error loading config")
	ErrNeedPointer              = errors.New("error must be pointer")
	ErrNewStoreEngineFail       = errors.New("error init store engine fail")
)
