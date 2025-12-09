package gen

import (
	"fmt"
	"math"
	"math/rand/v2"
	"time"

	"github.com/lance4117/gofuse/codec"
	"github.com/lance4117/gofuse/once"
	sf "github.com/sony/sonyflake/v2"
)

var getDefault = once.DoWithErr(func() (*sf.Sonyflake, error) {
	return sf.New(sf.Settings{})
})

// NewId 生成一个新的雪花ID
func NewId() (int64, error) {
	snowflake, err := getDefault()
	if err != nil {
		return 0, err
	}
	id, err := snowflake.NextID()
	if err != nil {
		return 0, err
	}
	if id > math.MaxInt64 {
		return 0, fmt.Errorf("snowflake id overflow: %d", id)
	}
	return id, nil
}

// ShortID 生成可读性强的 22~26 位 base62 短 ID（时间戳 + 随机）
func ShortID() string {
	b := make([]byte, 20)
	// 8 字节时间戳（毫秒） + 12 字节随机
	ts := time.Now().UnixMilli()
	for i := 7; i >= 0; i-- {
		b[i] = byte(ts & 0xff)
		ts >>= 8
	}
	for i := 8; i < 20; i++ {
		b[i] = byte(rand.IntN(256))
	}
	// base62 编码
	return codec.B62Encode(b)
}
