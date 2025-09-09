package writer

import (
	"encoding/csv"
	"fmt"
	"os"

	"gitee.com/lance4117/GoFuse/errs"
	"gitee.com/lance4117/GoFuse/times"
)

var DefaultCSVFileName = fmt.Sprintf("writer-%d", times.NowMilli())

type CSVWriter struct {
	file     *os.File
	writer   *csv.Writer
	filename string
}

func NewCSVWriter(filename string) *CSVWriter {
	if filename == "" {
		filename = DefaultCSVFileName
	}
	return &CSVWriter{file: nil, writer: nil, filename: filename}
}

func (w *CSVWriter) Init(headers []string) error {
	f, err := os.Create(w.filename + ".csv")
	if err != nil {
		return err
	}
	w.file = f
	w.writer = csv.NewWriter(f)
	return w.writer.Write(headers)
}

func (w *CSVWriter) Write(values []string) error {
	if w.writer == nil {
		return errs.ErrNotInitialized
	}
	if err := w.writer.Write(values); err != nil {
		return err
	}
	w.writer.Flush()
	return nil
}

func (w *CSVWriter) Close() error {
	if w.file != nil {
		return w.file.Close()
	}
	return nil
}
