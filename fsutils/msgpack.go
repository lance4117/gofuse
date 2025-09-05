package fsutils

import (
	"github.com/vmihailenco/msgpack/v5"
)

func MarshalAny(value any) ([]byte, error) {
	return msgpack.Marshal(value)
}

func UnmarshalTo[T any](bytes []byte) (T, error) {
	var value T
	if len(bytes) == 0 {
		return value, nil
	}
	err := msgpack.Unmarshal(bytes, &value)
	return value, err
}
