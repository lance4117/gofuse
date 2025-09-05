package fsid

import (
	"gitee.com/lance4117/GoFuse/fsonce"
	sf "github.com/sony/sonyflake/v2"
)

var getDefault = fsonce.DoWithErr(func() (*sf.Sonyflake, error) {
	return sf.New(sf.Settings{})
})

func NewId() (int64, error) {
	snowflake, err := getDefault()
	if err != nil {
		return 0, nil
	}
	return snowflake.NextID()
}
