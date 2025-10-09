package codec

import (
	"encoding/base64"
	"fmt"
	"io"

	"github.com/lance4117/gofuse/logger"
)

// 1. URL 安全 & Padding
//   RawURLEncoding 默认不带 =，与很多 Web 组件/JWT 的约定一致。
//   如果你要与外部系统交互，请确认对方是否需要 = padding。
//
// 2. 大小写敏感
//   Base64 是大小写敏感的，避免在链路上被错误地大小写转换。
//
// 3. 性能
//
//   结构体→Base64 尽量一次完成（json.Marshal + EncodeToString）。
//   大数据走流式，避免 []byte 巨大分配。
//
// 4. 安全
//   Base64 不是加密。涉及敏感信息请先加密（例如 AES-GCM）再 Base64。

// ========== 基础：标准 Base64（带 '=' padding） ==========

// B64Encode 将二进制编码为 base64 文本
func B64Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

// B64Decode 将 base64 文本解码为二进制
func B64Decode(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

// ========== 标准字母表、无 '='、含 '+' '/'（RawStdEncoding）==========

// B64RawEncode 将二进制编码为 base64 文本
func B64RawEncode(b []byte) string { //
	return base64.RawStdEncoding.EncodeToString(b)
}

// B64RawDecode 将 base64 文本解码为二进制
func B64RawDecode(s string) ([]byte, error) {
	return base64.RawStdEncoding.DecodeString(s)
}

// ========== URL 安全：不含 '+' '/' '='（RawURLEncoding）==========

// B64URLEncode 无填充、URL 安全的 base64
func B64URLEncode(b []byte) string {
	return base64.RawURLEncoding.EncodeToString(b)
}

// B64URLDecode 解码 URL 安全 base64
func B64URLDecode(s string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(s)
}

// ========== JSON + Base64 组合：一步到位 ==========

// B64EncodeJSON 将任意结构体 -> json -> base64
func B64EncodeJSON(v any) (string, error) {
	b, err := JSONMarshal(v)
	if err != nil {
		return "", err
	}
	return B64Encode(b), nil
}

// B64DecodeJSON 将 base64 -> json -> 结构体（泛型版）
func B64DecodeJSON[T any](s string) (T, error) {
	var v T
	if s == "" {
		return v, nil
	}
	b, err := B64Decode(s)
	if err != nil {
		return v, err
	}
	err = JSONUnmarshal(b, &v)
	return v, err
}

// B64URLEncodeJSON URL 安全版（无 '='）
func B64URLEncodeJSON(v any) (string, error) {
	b, err := JSONMarshal(v)
	if err != nil {
		return "", err
	}
	return B64URLEncode(b), nil
}

// B64URLDecodeJSON URL 安全版（泛型）
func B64URLDecodeJSON[T any](s string) (T, error) {
	var v T
	if s == "" {
		return v, nil
	}
	b, err := B64URLDecode(s)
	if err != nil {
		return v, err
	}
	err = JSONUnmarshal(b, &v)
	return v, err
}

// ========== 流式：适合大对象/文件 ==========

// B64EncodeStream 从 r 读原始二进制，写 base64 文本到 w
func B64EncodeStream(w io.Writer, r io.Reader) (int64, error) {
	enc := base64.NewEncoder(base64.StdEncoding, w)
	n, err := io.Copy(enc, r)
	cerr := enc.Close()
	if err != nil {
		return n, err
	}
	return n, cerr
}

// B64DecodeStream 从 r 读 base64 文本，写解码后的二进制到 w
func B64DecodeStream(w io.Writer, r io.Reader) (int64, error) {
	dec := base64.NewDecoder(base64.StdEncoding, r)
	return io.Copy(w, dec)
}

// B64URLEncodeStream URL 安全流式（可选）
func B64URLEncodeStream(w io.Writer, r io.Reader) error {
	enc := base64.NewEncoder(base64.RawURLEncoding, w)
	if _, err := io.Copy(enc, r); err != nil {
		_ = enc.Close()
		return err
	}
	return enc.Close()
}

func B64URLDecodeStream(w io.Writer, r io.Reader) error {
	dec := base64.NewDecoder(base64.RawURLEncoding, r)
	_, err := io.Copy(w, dec)
	return err
}

// ========== 小工具/安全兜底 ==========

// MustB64EncodeJSON 便捷方法（panic on error，谨慎用于初始化场景）
func MustB64EncodeJSON(v any) string {
	s, err := B64EncodeJSON(v)
	if err != nil {
		logger.Panic(fmt.Errorf("MustB64EncodeJSON: %w", err))
	}
	return s
}
