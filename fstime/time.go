package fstime

import "time"

const (
	Nanosecond  int64 = 1
	Microsecond       = 1000 * Nanosecond
	Millisecond       = 1000 * Microsecond
	Second            = 1000 * Millisecond
	Minute            = 60 * Second
	Hour              = 60 * Minute
)

// NowMilli 当前时间的毫秒时间戳
func NowMilli() int64 {
	return time.Now().UnixMilli()
}

// NowDateTime 当前日期时间
func NowDateTime() string {
	return time.Now().Format(time.DateTime)
}

// NowDateOnly 当前日期，没有时间
func NowDateOnly() string {
	return time.Now().Format(time.DateOnly)
}

// NowAfter 当前时间变化后的毫秒时间戳，负数是以前
func NowAfter(dur int64) int64 {
	return time.Now().Add(time.Duration(dur)).UnixMilli()
}

// ToDateTime 将毫秒时间戳变为时间日期格式
func ToDateTime(milliseconds int64) string {
	return time.UnixMilli(milliseconds).Format(time.DateTime)
}

// ToDateOnly 将毫秒时间戳变为年月日格式
func ToDateOnly(milliseconds int64) string {
	return time.UnixMilli(milliseconds).Format(time.DateOnly)
}
