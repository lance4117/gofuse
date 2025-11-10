package fileio

import "io"

// Writer 定义通用文件写入接口
type Writer interface {
	Write(data any) error
	Flush() error
	io.Closer
}

// Reader 定义通用文件读取接口
type Reader interface {
	Read() (any, error)
	io.Closer
}

// TableWriter 定义表格类文件（CSV/Excel等）的写入接口
type TableWriter interface {
	WriteHeader(headers []string) error
	WriteRow(values []string) error
	Write(data any) error
	Flush() error
	io.Closer
}

// TableReader 定义表格类文件的读取接口
type TableReader interface {
	ReadAll() ([][]string, error)
	ReadHeader() ([]string, error)
	ReadRow() ([]string, error)
	Read() (any, error)
	io.Closer
}

// StructuredWriter 定义结构化文件（JSON/XML等）的写入接口
type StructuredWriter interface {
	WriteObject(obj any) error
	WriteArray(objs []any) error
	Write(data any) error
	Flush() error
	io.Closer
}

// StructuredReader 定义结构化文件的读取接口
type StructuredReader interface {
	ReadObject(obj any) error
	ReadArray(objs any) error
	Read() (any, error)
	io.Closer
}
