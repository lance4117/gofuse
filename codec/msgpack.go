package codec

import (
	"github.com/vmihailenco/msgpack/v5"
)

// MPMarshal 将value转为bytes
func MPMarshal(value any) ([]byte, error) {
	return msgpack.Marshal(value)
}

// MPUnmarshal 将bytes反序列化为value相对应的结构
func MPUnmarshal(bytes []byte, value any) error {
	if value == nil {
		return msgpack.Unmarshal(bytes, nil)
	}
	return msgpack.Unmarshal(bytes, value)
}

// MPUnmarshalTo 将bytes反序列化为泛型相对应的结构
func MPUnmarshalTo[T any](bytes []byte) (T, error) {
	var value T
	if len(bytes) == 0 {
		return value, nil
	}
	err := msgpack.Unmarshal(bytes, &value)
	return value, err
}
