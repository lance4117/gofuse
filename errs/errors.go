package errs

import (
	"errors"
)

var (
	ErrNeedPointer = errors.New(" must be pointer ")
)

// fileio
var (
	ErrFileReaderNotInitialized = errors.New(" file reader not initialized ")
	ErrFileWriteNotInitialized  = errors.New(" file writer not initialized ")
)

// config
var (
	ErrConfigRead = errors.New(" reading config ")
	ErrConfigLoad = errors.New(" reading loading ")
)

// store
var (
	ErrNewStoreEngineFail = errors.New(" init storage engine fail ")
	ErrKeyNotFound        = errors.New(" key not found ")
)

// chain
var (
	ErrNoBalance    = errors.New(" no balance ")
	ErrGrpcConnFail = errors.New(" ping grpc connection fail ")
	ErrNoAmount     = errors.New(" no amount ")
)

var (
	ErrAESKeyLength = errors.New(" key length must be 16,24,32")
)
