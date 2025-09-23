package fileio

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/lance4117/gofuse/errs"
	"github.com/lance4117/gofuse/times"
)

var DefaultCSVFileName = fmt.Sprintf("writer-%d", times.NowMilli())

// CSVFileIO CSV文件读写器结构体
type CSVFileIO struct {
	File     *os.File
	Writer   *csv.Writer
	Reader   *csv.Reader
	Filename string
}

// NewCSVFileIO 创建一个新的CSV写入器实例
// pathAndName: CSV文件路径 eg: ./path/to/filename
// 返回CSVWriter指针
func NewCSVFileIO(pathAndName string) *CSVFileIO {
	if pathAndName == "" {
		pathAndName = DefaultCSVFileName
	}
	return &CSVFileIO{Filename: pathAndName}
}

// Create 创建CSV文件并可选择写入表头(读模式)
// headers: 可选的表头行数据
func (w *CSVFileIO) Create(headers []string) error {
	var err error
	w.File, err = os.Create(w.Filename + ".csv")
	if err != nil {
		return err
	}
	w.Writer = csv.NewWriter(w.File)
	if headers != nil {
		return w.Write(headers)
	}
	return nil
}

// Open 打开已存在的CSV文件用于读取(写模式)
func (w *CSVFileIO) Open() error {
	var err error
	w.File, err = os.Open(w.Filename + ".csv")
	if err != nil {
		return err
	}
	w.Reader = csv.NewReader(w.File)
	return nil
}

// Write 向CSV文件写入一行数据
// values: 要写入的字符串切片
func (w *CSVFileIO) Write(values []string) error {
	if w.Writer == nil {
		return errs.ErrFileWriteNotInitialized
	}
	if err := w.Writer.Write(values); err != nil {
		return err
	}
	w.Writer.Flush()
	return nil
}

// ReadAll 读取CSV文件的所有数据
// 返回二维字符串切片包含所有行数据
func (w *CSVFileIO) ReadAll() ([][]string, error) {
	if w.Reader == nil {
		return nil, errs.ErrFileReaderNotInitialized
	}
	return w.Reader.ReadAll()
}

// Close 关闭CSV文件并刷新缓冲区
func (w *CSVFileIO) Close() error {
	if w.Writer != nil {
		w.Writer.Flush()
	}
	if w.File != nil {
		return w.File.Close()
	}
	return nil
}
