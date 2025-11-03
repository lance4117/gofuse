package kvs

type KVStore interface {
	Put(key string, val []byte) error
	Get(key string) ([]byte, error)
	Del(key string) error
	Has(key string) (bool, error)

	PutSync(key string, val []byte) error
	DelSync(key string) error

	// NewBatch 批量操作
	NewBatch() (Batch, error)
	// NewIterator 迭代器相关
	NewIterator(opt IterOption) (Iterator, error)

	Close() error
}

type IterOption struct {
	Prefix []byte // 如果设置，遍历该前缀
	Start  []byte // inclusive
	End    []byte // exclusive
}

type Batch interface {
	Put(key string, val []byte) error
	Del(key string) error

	PutSync(key string, val []byte) error
	DelSync(key string) error

	Commit() error
	Cancel() error
}

type Iterator interface {
	Valid() bool
	Next() bool
	Key() []byte
	Value() ([]byte, error)
	Close() error
}
