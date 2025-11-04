package kvs

// KVStore 定义了键值存储接口，提供基本的增删改查操作以及批量操作和迭代器功能
type KVStore interface {
	// Put 将指定的键值对存储到数据库中，异步写入
	Put(key string, val []byte) error
	// Get 根据指定的键从数据库中获取对应的值
	Get(key string) ([]byte, error)
	// Del 根据指定的键从数据库中删除对应的键值对
	Del(key string) error
	// Has 检查指定的键是否存在于数据库中
	Has(key string) (bool, error)

	// PutSync 将指定的键值对存储到数据库中，同步写入（确保数据已落盘）
	PutSync(key string, val []byte) error
	// DelSync 根据指定的键从数据库中删除对应的键值对，同步删除（确保数据已落盘）
	DelSync(key string) error

	// NewBatch 批量操作
	NewBatch() (Batch, error)
	// NewIterator 迭代器相关
	NewIterator(opt IterOption) (Iterator, error)

	// Close 关闭数据库连接并释放相关资源
	Close() error
}

// IterOption 定义了迭代器的选项配置
type IterOption struct {
	Prefix []byte // 如果设置，遍历该前缀
	Start  []byte // inclusive
	End    []byte // exclusive
}

// Batch 定义了批量操作接口，用于执行多个操作后一次性提交
type Batch interface {
	// Put 将指定的键值对添加到批处理队列中，异步写入
	Put(key string, val []byte) error
	// Del 根据指定的键将删除操作添加到批处理队列中
	Del(key string) error

	// PutSync 将指定的键值对添加到批处理队列中，同步写入
	PutSync(key string, val []byte) error
	// DelSync 根据指定的键将删除操作添加到批处理队列中，同步删除
	DelSync(key string) error

	// Commit 提交批处理队列中的所有操作
	Commit() error
	// Cancel 取消批处理队列中的所有操作
	Cancel() error
}

// Iterator 定义了迭代器接口，用于遍历数据库中的键值对
type Iterator interface {
	// Valid 检查当前迭代器位置是否有效
	Valid() bool
	// Next 将迭代器移动到下一个位置
	Next() bool
	// Key 获取当前迭代器位置的键
	Key() []byte
	// Value 获取当前迭代器位置的值
	Value() ([]byte, error)
	// Close 关闭迭代器并释放相关资源
	Close() error
}
