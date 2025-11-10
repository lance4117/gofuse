package fileio

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"

	"github.com/lance4117/gofuse/errs"
	"github.com/lance4117/gofuse/times"
)

var DefaultXMLFileName = fmt.Sprintf("writer-%d.xml", times.NowMilli())

// XMLWriter XML文件写入器
type XMLWriter struct {
	file     *os.File
	encoder  *xml.Encoder
	filename string
	indent   string // 缩进字符串，默认为 "  " (两个空格)
}

// XMLWriterOption XML写入器配置选项
type XMLWriterOption func(*XMLWriter)

// WithXMLIndent 设置XML缩进
func WithXMLIndent(indent string) XMLWriterOption {
	return func(w *XMLWriter) {
		w.indent = indent
	}
}

// NewXMLWriter 创建一个新的XML写入器实例
// pathAndName: XML文件路径 eg: ./path/to/filename.xml
// 返回 XMLWriter 指针和错误
func NewXMLWriter(pathAndName string, opts ...XMLWriterOption) (*XMLWriter, error) {
	if pathAndName == "" {
		pathAndName = DefaultXMLFileName
	}

	filename := ensureXMLExtension(pathAndName)
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	w := &XMLWriter{
		file:     file,
		encoder:  xml.NewEncoder(file),
		filename: filename,
		indent:   "  ", // 默认两个空格缩进
	}

	// 应用配置选项
	for _, opt := range opts {
		opt(w)
	}

	// 设置编码器缩进
	w.encoder.Indent("", w.indent)

	// 写入 XML 声明
	file.WriteString(xml.Header)

	return w, nil
}

// WriteObject 写入单个对象
func (w *XMLWriter) WriteObject(obj any) error {
	if w.encoder == nil {
		return errs.ErrFileWriteNotInitialized
	}
	return w.encoder.Encode(obj)
}

// WriteArray 写入数组（包装在根元素中）
func (w *XMLWriter) WriteArray(objs []any) error {
	if w.encoder == nil {
		return errs.ErrFileWriteNotInitialized
	}

	// XML需要一个根元素，这里使用 <items> 作为默认根元素
	type Items struct {
		XMLName xml.Name `xml:"items"`
		Items   []any    `xml:"item"`
	}

	wrapper := Items{Items: objs}
	return w.encoder.Encode(wrapper)
}

// Write 通用写入接口实现
func (w *XMLWriter) Write(data any) error {
	if w.encoder == nil {
		return errs.ErrFileWriteNotInitialized
	}
	if err := w.encoder.Encode(data); err != nil {
		return err
	}
	return w.Flush()
}

// Flush 刷新缓冲区
func (w *XMLWriter) Flush() error {
	if w.encoder != nil {
		if err := w.encoder.Flush(); err != nil {
			return err
		}
	}
	if w.file != nil {
		return w.file.Sync()
	}
	return nil
}

// Close 关闭XML文件
func (w *XMLWriter) Close() error {
	var firstErr error

	// 刷新缓冲区
	if w.encoder != nil {
		if err := w.encoder.Flush(); err != nil {
			firstErr = err
		}
	}

	if w.file != nil {
		// 写入换行符，使文件格式更友好
		w.file.WriteString("\n")

		if err := w.file.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}

	return firstErr
}

// XMLReader XML文件读取器
type XMLReader struct {
	file     *os.File
	decoder  *xml.Decoder
	filename string
}

// NewXMLReader 创建一个新的XML读取器实例
// pathAndName: XML文件路径 eg: ./path/to/filename.xml
// 返回 XMLReader 指针和错误
func NewXMLReader(pathAndName string) (*XMLReader, error) {
	filename := ensureXMLExtension(pathAndName)
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	return &XMLReader{
		file:     file,
		decoder:  xml.NewDecoder(file),
		filename: filename,
	}, nil
}

// ReadObject 读取单个对象到指定的结构体
func (r *XMLReader) ReadObject(obj any) error {
	if r.decoder == nil {
		return errs.ErrFileReaderNotInitialized
	}
	return r.decoder.Decode(obj)
}

// ReadArray 读取数组到指定的切片
func (r *XMLReader) ReadArray(objs any) error {
	if r.decoder == nil {
		return errs.ErrFileReaderNotInitialized
	}
	return r.decoder.Decode(objs)
}

// Read 通用读取接口实现
func (r *XMLReader) Read() (any, error) {
	if r.decoder == nil {
		return nil, errs.ErrFileReaderNotInitialized
	}
	var result any
	if err := r.decoder.Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

// Close 关闭XML文件
func (r *XMLReader) Close() error {
	if r.file != nil {
		return r.file.Close()
	}
	return nil
}

func ensureXMLExtension(filename string) string {
	if !strings.HasSuffix(filename, ".xml") {
		filename += ".xml"
	}
	return filename
}
