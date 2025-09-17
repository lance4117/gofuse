package fileio

// Files 定义文件读写接口
type Files interface {
	Create(headers []string) error
	Open() error
	Write(values []string) error
	ReadAll() ([][]string, error)
	Close() error
}
