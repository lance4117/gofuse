package conv

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/lance4117/gofuse/errs"
)

// Int64ToBytes 将 int64 转换为大端序字节数组
func Int64ToBytes(v int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

// BytesToInt64 将大端序字节数组转换为 int64
func BytesToInt64(b []byte) (int64, error) {
	if len(b) != 8 {
		return 0, errs.ErrBigEndianLength
	}
	return int64(binary.BigEndian.Uint64(b)), nil
}

// IntToStr 将整数转换为字符串
func IntToStr(v int) string { return strconv.Itoa(v) }

// Int64ToStr 将64位整数转换为字符串
func Int64ToStr(v int64) string { return strconv.FormatInt(v, 10) }

// Uint64ToStr 将64位无符号整数转换为字符串
func Uint64ToStr(v uint64) string { return strconv.FormatUint(v, 10) }

// FloatToStr 将浮点数转换为字符串
// prec 参数指定小数点后保留的位数
func FloatToStr(v float64, prec int) string {
	return strconv.FormatFloat(v, 'f', prec, 64)
}

// BoolToStr 将布尔值转换为字符串("true"或"false")
func BoolToStr(v bool) string { return strconv.FormatBool(v) }

// StrToInt 将字符串解析为整数
func StrToInt(s string) (int, error) {
	return strconv.Atoi(s)
}

// StrToInt64 将字符串解析为64位整数
func StrToInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// StrToUint64 将字符串解析为64位无符号整数
func StrToUint64(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}

// StrToFloat 将字符串解析为浮点数
func StrToFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

// StrToBool 将字符串解析为布尔值
// 支持的字符串: "true", "false", "1", "0", "t", "f", "TRUE", "FALSE"等
func StrToBool(s string) (bool, error) {
	return strconv.ParseBool(s)
}

// BytesToHex 将字节数组转换为十六进制字符串
func BytesToHex(b []byte) string { return hex.EncodeToString(b) }

// HexToBytes 将十六进制字符串解析为字节数组
func HexToBytes(s string) ([]byte, error) { return hex.DecodeString(s) }

// AnyToString 会根据类型优化常见路径，其他走 fmt.Sprintf
// 使用泛型支持任意类型到字符串的转换
func AnyToString[T any](v T) string {
	switch val := any(v).(type) {
	case int:
		return strconv.Itoa(val)
	case int64:
		return strconv.FormatInt(val, 10)
	case uint64:
		return strconv.FormatUint(val, 10)
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(val)
	case string:
		return val
	default:
		return fmt.Sprintf("%v", val)
	}
}
