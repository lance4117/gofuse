package fileio

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/lance4117/gofuse/errs"
	"github.com/lance4117/gofuse/times"
)

var DefaultJSONFileName = fmt.Sprintf("writer-%d.json", times.NowMilli())

// JSONWriter JSON文件写入器
type JSONWriter struct {
	file     *os.File
	encoder  *json.Encoder
	filename string
	indent   string // 缩进字符串，默认为 "  " (两个空格)
}

// JSONWriterOption JSON写入器配置选项
type JSONWriterOption func(*JSONWriter)

// WithIndent 设置JSON缩进
func WithIndent(indent string) JSONWriterOption {
	return func(w *JSONWriter) {
		w.indent = indent
	}
}

// NewJSONWriter 创建一个新的JSON写入器实例
// pathAndName: JSON文件路径 eg: ./path/to/filename.json
// 返回 JSONWriter 指针和错误
func NewJSONWriter(pathAndName string, opts ...JSONWriterOption) (*JSONWriter, error) {
	if pathAndName == "" {
		pathAndName = DefaultJSONFileName
	}

	filename := ensureJSONExtension(pathAndName)
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	w := &JSONWriter{
		file:     file,
		encoder:  json.NewEncoder(file),
		filename: filename,
		indent:   "  ", // 默认两个空格缩进
	}

	// 应用配置选项
	for _, opt := range opts {
		opt(w)
	}

	// 设置编码器缩进
	w.encoder.SetIndent("", w.indent)

	return w, nil
}

// WriteObject 写入单个对象
func (w *JSONWriter) WriteObject(obj any) error {
	if w.encoder == nil {
		return errs.ErrFileWriteNotInitialized
	}
	return w.encoder.Encode(obj)
}

// WriteArray 写入数组
func (w *JSONWriter) WriteArray(objs []any) error {
	if w.encoder == nil {
		return errs.ErrFileWriteNotInitialized
	}
	return w.encoder.Encode(objs)
}

// Write 通用写入接口实现
func (w *JSONWriter) Write(data any) error {
	if w.encoder == nil {
		return errs.ErrFileWriteNotInitialized
	}
	if err := w.encoder.Encode(data); err != nil {
		return err
	}
	return w.Flush()
}

// Flush 刷新缓冲区（JSON编码器会自动刷新，此处用于接口一致性）
func (w *JSONWriter) Flush() error {
	if w.file != nil {
		return w.file.Sync()
	}
	return nil
}

// Close 关闭JSON文件
func (w *JSONWriter) Close() error {
	var firstErr error

	// 刷新缓冲区
	if err := w.Flush(); err != nil {
		firstErr = err
	}

	if w.file != nil {
		if err := w.file.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}

	return firstErr
}

// JSONReader JSON文件读取器
type JSONReader struct {
	file     *os.File
	decoder  *json.Decoder
	filename string
}

// NewJSONReader 创建一个新的JSON读取器实例
// pathAndName: JSON文件路径 eg: ./path/to/filename.json
// 返回 JSONReader 指针和错误
func NewJSONReader(pathAndName string) (*JSONReader, error) {
	filename := ensureJSONExtension(pathAndName)
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	return &JSONReader{
		file:     file,
		decoder:  json.NewDecoder(file),
		filename: filename,
	}, nil
}

// ReadObject 读取单个对象到指定的结构体
func (r *JSONReader) ReadObject(obj any) error {
	if r.decoder == nil {
		return errs.ErrFileReaderNotInitialized
	}
	return r.decoder.Decode(obj)
}

// ReadArray 读取数组到指定的切片
func (r *JSONReader) ReadArray(objs any) error {
	if r.decoder == nil {
		return errs.ErrFileReaderNotInitialized
	}
	return r.decoder.Decode(objs)
}

// Read 通用读取接口实现，读取为 map[string]any 或 []any
func (r *JSONReader) Read() (any, error) {
	if r.decoder == nil {
		return nil, errs.ErrFileReaderNotInitialized
	}
	var result any
	if err := r.decoder.Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

// Close 关闭JSON文件
func (r *JSONReader) Close() error {
	if r.file != nil {
		return r.file.Close()
	}
	return nil
}

func ensureJSONExtension(filename string) string {
	if !strings.HasSuffix(filename, ".json") {
		filename += ".json"
	}
	return filename
}
