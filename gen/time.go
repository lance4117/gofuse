package gen

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

// DateRange 在给定区间生成随机时间
func DateRange(start, end time.Time) time.Time {
	return gofakeit.DateRange(start, end)
}

// NowRecent 生成目前的近 duration 内的时间
func NowRecent(dur time.Duration) time.Time {
	return time.Now().Add(-time.Duration(gofakeit.Number(0, int(dur.Milliseconds()))) * time.Millisecond)
}
