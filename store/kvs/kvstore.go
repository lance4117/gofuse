package kvs

// KVStore 定义了键值存储接口，提供基本的增删改查操作以及批量操作和迭代器功能
type KVStore interface {
	// Put 将指定的键值对存储到数据库中，异步写入
	// key: 要存储的键
	// val: 要存储的值
	// 返回写入过程中可能发生的错误
	Put(key string, val []byte) error

	// Get 根据指定的键从数据库中获取对应的值
	// key: 要查询的键
	// 返回查询到的值和查询过程中可能发生的错误
	Get(key string) ([]byte, error)

	// Del 根据指定的键从数据库中删除对应的键值对
	// key: 要删除的键
	// 返回删除过程中可能发生的错误
	Del(key string) error

	// Has 检查指定的键是否存在于数据库中
	// key: 要检查的键
	// 返回是否存在布尔值和检查过程中可能发生的错误
	Has(key string) (bool, error)

	// PutSync 将指定的键值对存储到数据库中，同步写入（确保数据已落盘）
	// key: 要存储的键
	// val: 要存储的值
	// 返回写入过程中可能发生的错误
	PutSync(key string, val []byte) error

	// DelSync 根据指定的键从数据库中删除对应的键值对，同步删除（确保数据已落盘）
	// key: 要删除的键
	// 返回删除过程中可能发生的错误
	DelSync(key string) error

	// NewBatch 批量操作
	NewBatch() (Batch, error)

	// NewIterator 迭代器相关
	NewIterator(opt IterOption) (Iterator, error)

	// Close 关闭数据库连接并释放相关资源
	// 返回关闭过程中可能发生的错误
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
	// key: 要存储的键
	// val: 要存储的值
	// 返回添加过程中可能发生的错误
	Put(key string, val []byte) error

	// Del 根据指定的键将删除操作添加到批处理队列中
	// key: 要删除的键
	// 返回添加过程中可能发生的错误
	Del(key string) error

	// PutSync 将指定的键值对添加到批处理队列中，同步写入
	// key: 要存储的键
	// val: 要存储的值
	// 返回添加过程中可能发生的错误
	PutSync(key string, val []byte) error

	// DelSync 根据指定的键将删除操作添加到批处理队列中，同步删除
	// key: 要删除的键
	// 返回添加过程中可能发生的错误
	DelSync(key string) error

	// Commit 提交批处理队列中的所有操作
	// 返回提交过程中可能发生的错误
	Commit() error

	// Cancel 取消批处理队列中的所有操作
	// 返回取消过程中可能发生的错误
	Cancel() error
}

// Iterator 定义了迭代器接口，用于遍历数据库中的键值对
type Iterator interface {
	// Valid 检查当前迭代器位置是否有效
	// 返回当前位置是否有效的布尔值
	Valid() bool

	// Next 将迭代器移动到下一个位置
	// 返回移动后位置是否有效的布尔值
	Next() bool

	// Key 获取当前迭代器位置的键
	// 返回当前键的字节切片
	Key() []byte

	// Value 获取当前迭代器位置的值
	// 返回当前值的字节切片和获取过程中可能发生的错误
	Value() ([]byte, error)

	// Close 关闭迭代器并释放相关资源
	// 返回关闭过程中可能发生的错误
	Close() error
}
