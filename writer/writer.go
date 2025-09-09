package writer

// Writer 定义输出器接口
type Writer interface {
	Init(headers []string) error
	Write(values []string) error
	Close() error
}
