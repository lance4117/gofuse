package fileio

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/lance4117/gofuse/errs"
	"github.com/lance4117/gofuse/times"
)

var DefaultCSVFileName = fmt.Sprintf("writer-%d.csv", times.NowMilli())

// CSVWriter CSV文件写入器
type CSVWriter struct {
	file     *os.File
	writer   *csv.Writer
	filename string
}

// NewCSVWriter 创建一个新的CSV写入器实例
// pathAndName: CSV文件路径 eg: ./path/to/filename.csv
// 返回 CSVWriter 指针和错误
func NewCSVWriter(pathAndName string) (*CSVWriter, error) {
	if pathAndName == "" {
		pathAndName = DefaultCSVFileName
	}

	filename := ensureCSVExtension(pathAndName)
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	return &CSVWriter{
		file:     file,
		writer:   csv.NewWriter(file),
		filename: filename,
	}, nil
}

// WriteHeader 写入CSV表头
func (w *CSVWriter) WriteHeader(headers []string) error {
	return w.WriteRow(headers)
}

// WriteRow 写入一行数据
func (w *CSVWriter) WriteRow(values []string) error {
	if w.writer == nil {
		return errs.ErrFileWriteNotInitialized
	}
	return w.writer.Write(values)
}

// Write 通用写入接口实现，支持 []string 或 [][]string
func (w *CSVWriter) Write(data any) error {
	switch v := data.(type) {
	case []string:
		if err := w.WriteRow(v); err != nil {
			return err
		}
		return w.Flush()
	case [][]string:
		for _, row := range v {
			if err := w.WriteRow(row); err != nil {
				return err
			}
		}
		return w.Flush()
	default:
		return errs.ErrUnsupportedDataType
	}
}

// Flush 刷新缓冲区
func (w *CSVWriter) Flush() error {
	if w.writer != nil {
		w.writer.Flush()
		return w.writer.Error()
	}
	return nil
}

// Close 关闭CSV文件并刷新缓冲区
func (w *CSVWriter) Close() error {
	var firstErr error

	if w.writer != nil {
		w.writer.Flush()
		if err := w.writer.Error(); err != nil {
			firstErr = err
		}
	}

	if w.file != nil {
		if err := w.file.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}

	return firstErr
}

// CSVReader CSV文件读取器
type CSVReader struct {
	file     *os.File
	reader   *csv.Reader
	filename string
}

// NewCSVReader 创建一个新的CSV读取器实例
// pathAndName: CSV文件路径 eg: ./path/to/filename.csv
// 返回 CSVReader 指针和错误
func NewCSVReader(pathAndName string) (*CSVReader, error) {
	filename := ensureCSVExtension(pathAndName)
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	return &CSVReader{
		file:     file,
		reader:   csv.NewReader(file),
		filename: filename,
	}, nil
}

// ReadHeader 读取CSV表头（第一行）
func (r *CSVReader) ReadHeader() ([]string, error) {
	if r.reader == nil {
		return nil, errs.ErrFileReaderNotInitialized
	}
	return r.reader.Read()
}

// ReadRow 逐行读取CSV数据
func (r *CSVReader) ReadRow() ([]string, error) {
	if r.reader == nil {
		return nil, errs.ErrFileReaderNotInitialized
	}
	return r.reader.Read()
}

// ReadAll 读取CSV文件的所有数据
func (r *CSVReader) ReadAll() ([][]string, error) {
	if r.reader == nil {
		return nil, errs.ErrFileReaderNotInitialized
	}
	return r.reader.ReadAll()
}

// Read 通用读取接口实现，返回所有数据
func (r *CSVReader) Read() (any, error) {
	return r.ReadAll()
}

// Close 关闭CSV文件
func (r *CSVReader) Close() error {
	if r.file != nil {
		return r.file.Close()
	}
	return nil
}

func ensureCSVExtension(filename string) string {
	if !strings.HasSuffix(filename, ".csv") {
		filename += ".csv"
	}
	return filename
}
