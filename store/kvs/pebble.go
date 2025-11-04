package kvs

import (
	"errors"

	"github.com/cockroachdb/pebble"
	"github.com/lance4117/gofuse/errs"
	"github.com/lance4117/gofuse/logger"
)

// PebbleConfig 配置 Pebble 数据库的结构体。
type PebbleConfig struct {
	DirName string
	*pebble.Options
}

// PebbleKV 封装了 pebble.DB 实例，提供键值存储的基本操作接口。
type PebbleKV struct {
	db *pebble.DB
}

// PebbleBatch 封装了 pebble.Batch 实例，支持批量写入操作。
type PebbleBatch struct {
	*pebble.Batch
}

// PebbleIterator 封装了 pebble.Iterator 实例，用于遍历数据库中的键值对。
type PebbleIterator struct {
	it *pebble.Iterator
}

// DefaultPebbleConfig 创建一个默认配置的 PebbleConfig 实例。
func DefaultPebbleConfig(dirname string) PebbleConfig {
	opts := &pebble.Options{}
	opts.EnsureDefaults()
	return PebbleConfig{dirname, opts}
}

// NewPebbleKV 根据给定的配置创建一个新的 PebbleKV 实例。
func NewPebbleKV(config PebbleConfig) (*PebbleKV, error) {
	// 获取实例
	kv, err := pebble.Open(config.DirName, config.Options)
	if err != nil {
		logger.Error(err, errs.ErrNewStoreEngineFail)
		return nil, err
	}
	return &PebbleKV{kv}, nil
}

// Put 使用NoSync策略向数据库中插入或更新指定 key 对应的 value 值。
func (kv *PebbleKV) Put(key string, val []byte) error {
	return kv.db.Set([]byte(key), val, pebble.NoSync)
}

// Get 根据 key 查询对应的 value。
func (kv *PebbleKV) Get(key string) ([]byte, error) {
	b, closer, err := kv.db.Get([]byte(key))
	if err != nil {
		if errors.Is(err, pebble.ErrNotFound) {
			return nil, errs.ErrKeyNotFound
		}
		return nil, err
	}
	defer closer.Close()
	// 必须复制
	out := make([]byte, len(b))
	copy(out, b)
	return out, nil
}

// Del 删除数据库中与 key 相关的数据项。
func (kv *PebbleKV) Del(key string) error {
	return kv.db.Delete([]byte(key), pebble.NoSync)
}

// Has 判断数据库中是否存在某个 key。
func (kv *PebbleKV) Has(key string) (bool, error) {
	_, closer, err := kv.db.Get([]byte(key))
	if err != nil {
		if errors.Is(err, pebble.ErrNotFound) {
			return false, nil
		}
		return false, err
	}
	_ = closer.Close()
	return true, nil
}

// PutSync 类似于 Put 方法，但在写入时使用 Sync 策略确保数据持久化。
func (kv *PebbleKV) PutSync(key string, val []byte) error {
	return kv.db.Set([]byte(key), val, pebble.Sync)
}

// DelSync 类似于 Del 方法，但在删除时使用 Sync 策略保证操作被刷盘。
func (kv *PebbleKV) DelSync(key string) error {
	return kv.db.Delete([]byte(key), pebble.Sync)
}

// NewBatch 创建一个新的批处理对象，用于批量写入操作。
func (kv *PebbleKV) NewBatch() (Batch, error) {
	batch := kv.db.NewBatch()
	return &PebbleBatch{batch}, nil

}

// Close 关闭当前数据库连接。
func (kv *PebbleKV) Close() error {
	return kv.db.Close()
}

// Put 在批处理中设置 key-value 映射关系。
func (b *PebbleBatch) Put(key string, val []byte) error {
	return b.Batch.Set([]byte(key), val, pebble.NoSync)
}

// Del 从批处理中移除指定 key 的映射。
func (b *PebbleBatch) Del(key string) error {
	return b.Batch.Delete([]byte(key), pebble.NoSync)
}

// PutSync 类似 Put 方法，但会强制将变更同步到磁盘。
func (b *PebbleBatch) PutSync(key string, val []byte) error {
	return b.Batch.Set([]byte(key), val, pebble.Sync)
}

// DelSync 类似 Del 方法，但会强制将变更同步到磁盘。
func (b *PebbleBatch) DelSync(key string) error {
	return b.Batch.Delete([]byte(key), pebble.Sync)
}

// Commit 提交批处理中的所有修改，并将其应用到底层数据库。
func (b *PebbleBatch) Commit() error {
	return b.Batch.Commit(pebble.Sync)
}

// Cancel 放弃本次批处理的所有更改。
func (b *PebbleBatch) Cancel() error {
	return b.Batch.Close()
}

// NewIterator 创建一个新的迭代器用于遍历数据库内容。
func (kv *PebbleKV) NewIterator(opt IterOption) (Iterator, error) {
	var ro *pebble.IterOptions
	if len(opt.Prefix) > 0 {
		upper := nextPrefix(opt.Prefix)
		ro = &pebble.IterOptions{LowerBound: opt.Prefix, UpperBound: upper}
	}
	it, err := kv.db.NewIter(ro)
	if err != nil {
		return nil, err
	}
	if ro != nil && ro.LowerBound != nil {
		it.SeekGE(ro.LowerBound)
	} else {
		it.First()
	}
	return &PebbleIterator{it: it}, nil
}

// Valid 判断当前迭代器是否处于有效位置（即指向一个存在的键值对）。
func (i *PebbleIterator) Valid() bool {
	return i.it.Valid()
}

// Next 移动迭代器至下一个键值对。
func (i *PebbleIterator) Next() bool {
	return i.it.Next()
}

// Key 获取当前迭代器所指键的内容副本。
func (i *PebbleIterator) Key() []byte {
	k := i.it.Key()
	out := make([]byte, len(k))
	copy(out, k)
	return out
}

// Value 获取当前迭代器所指值的内容副本。
func (i *PebbleIterator) Value() ([]byte, error) {
	v, err := i.it.ValueAndErr()
	if err != nil {
		return nil, err
	}
	out := make([]byte, len(v))
	copy(out, v)
	return out, nil
}

// Close 释放迭代器占用的资源。
func (i *PebbleIterator) Close() error {
	return i.it.Close()
}

// nextPrefix 计算出严格大于所有以输入 p 为前缀的字典序最小字节数组。
// 主要逻辑是对最后一个非 0xFF 字节加一并截断后续部分，
// 特殊情况下若全部为 0xFF 则返回 nil 表示无上限。
// 参数 p 为原前缀数组。
// 返回计算得到的新边界数组。
func nextPrefix(p []byte) []byte {
	// 返回按字典序严格大于所有以 p 为前缀的最小切片（常见做法：末位 +1，进位）
	out := append([]byte{}, p...)
	for i := len(out) - 1; i >= 0; i-- {
		if out[i] != 0xFF {
			out[i]++
			return out[:i+1]
		}
	}
	// 全是 0xFF，无更大上界，返回 nil 代表无上界
	return nil
}
