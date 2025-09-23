package gen

import (
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
		return 0, nil
	}
	return snowflake.NextID()
}
