package tools

import (
	"github.com/vmihailenco/msgpack/v5"
)

// MarshalAny 将value转为bytes
func MarshalAny(value any) ([]byte, error) {
	return msgpack.Marshal(value)
}

// UnmarshalAny 将bytes反序列化为value相对应的结构
func UnmarshalAny(bytes []byte, value any) error {
	return msgpack.Unmarshal(bytes, &value)
}

// UnmarshalTo 将bytes反序列化为泛型相对应的结构
func UnmarshalTo[T any](bytes []byte) (T, error) {
	var value T
	if len(bytes) == 0 {
		return value, nil
	}
	err := msgpack.Unmarshal(bytes, &value)
	return value, err
}
