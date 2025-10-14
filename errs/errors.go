package errs

import (
	"errors"
	"fmt"
)

var (
	ErrNil                      = errors.New(" error nil ")
	ErrFileReaderNotInitialized = errors.New(" error file reader not initialized ")
	ErrFileWriteNotInitialized  = errors.New(" error file writer not initialized ")
	ErrConfigRead               = errors.New(" error reading config ")
	ErrNeedPointer              = errors.New(" error must be pointer ")
	ErrNewStoreEngineFail       = errors.New(" error init store engine fail ")
)

func ErrConfigLoad(config string) error {
	return errors.New(fmt.Sprintf(" error loading config %s ", config))
}
