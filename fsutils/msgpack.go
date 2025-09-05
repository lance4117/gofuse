package fsutils

import (
	"gitee.com/lance4117/GoFuse/fserror"
	"github.com/vmihailenco/msgpack/v5"
)

func AnyToBytes(value any) ([]byte, error) {
	if value == nil {
		return nil, fserror.ErrNil
	}
	return msgpack.Marshal(value)
}

func BytesToAny[T any](bytes []byte) (T, error) {
	var value T
	err := msgpack.Unmarshal(bytes, &value)
	return value, err
}
