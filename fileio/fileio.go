package fileio

import "io"

// Writer 定义通用文件写入接口
type Writer interface {
	Write(data any) error
	io.Closer
}

// Reader 定义通用文件读取接口
type Reader interface {
	Read() (any, error)
	io.Closer
}

// TableWriter 定义表格类文件（CSV/Excel等）的写入接口
type TableWriter interface {
	Writer
	WriteHeader(headers []string) error
	WriteRow(values []string) error
}

// TableReader 定义表格类文件的读取接口
type TableReader interface {
	Reader
	ReadAll() ([][]string, error)
	ReadHeader() ([]string, error)
}

// StructuredWriter 定义结构化文件（JSON/XML等）的写入接口
type StructuredWriter interface {
	Writer
	WriteObject(obj any) error
	WriteArray(objs []any) error
}

// StructuredReader 定义结构化文件的读取接口
type StructuredReader interface {
	Reader
	ReadObject(obj any) error
	ReadArray(objs any) error
}
