package codec

import (
	"encoding/json"
)

// JSONMarshal 将value转为bytes
func JSONMarshal(value any) ([]byte, error) {
	return json.Marshal(value)
}

// JSONUnmarshal 将bytes反序列化为value相对应的结构
func JSONUnmarshal(bytes []byte, value any) error {
	return json.Unmarshal(bytes, value)
}

// JSONUnmarshalTo 将bytes反序列化为泛型相对应的结构
func JSONUnmarshalTo[T any](bytes []byte) (T, error) {
	var value T
	if len(bytes) == 0 {
		return value, nil
	}
	err := json.Unmarshal(bytes, &value)
	return value, err
}
