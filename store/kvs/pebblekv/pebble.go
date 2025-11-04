package pebblekv

import (
	"errors"

	"github.com/cockroachdb/pebble"
	"github.com/lance4117/gofuse/errs"
)

// PebbleKV 封装了 pebble.DB 实例，提供键值存储的基本操作接口。
type PebbleKV struct {
	PebbleDB *pebble.DB
}

// Put 使用NoSync策略向数据库中插入或更新指定 key 对应的 value 值。
func (kv *PebbleKV) Put(key string, val []byte) error {
	return kv.PebbleDB.Set([]byte(key), val, pebble.NoSync)
}

// Get 根据 key 查询对应的 value。
func (kv *PebbleKV) Get(key string) ([]byte, error) {
	b, closer, err := kv.PebbleDB.Get([]byte(key))
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
	return kv.PebbleDB.Delete([]byte(key), pebble.NoSync)
}

// Has 判断数据库中是否存在某个 key。
func (kv *PebbleKV) Has(key string) (bool, error) {
	_, closer, err := kv.PebbleDB.Get([]byte(key))
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
	return kv.PebbleDB.Set([]byte(key), val, pebble.Sync)
}

// DelSync 类似于 Del 方法，但在删除时使用 Sync 策略保证操作被刷盘。
func (kv *PebbleKV) DelSync(key string) error {
	return kv.PebbleDB.Delete([]byte(key), pebble.Sync)
}

// NewBatch 创建一个新的批处理对象，用于批量写入操作。
func (kv *PebbleKV) NewBatch() *pebble.Batch {
	return kv.PebbleDB.NewBatch()
}

// Close 关闭当前数据库连接。
func (kv *PebbleKV) Close() error {
	return kv.PebbleDB.Close()
}

// IterOption 定义了迭代器的选项配置
type IterOption struct {
	Prefix []byte // 如果设置，遍历该前缀
	Start  []byte // inclusive
	End    []byte // exclusive
}

// NewIterator 创建一个新的迭代器用于遍历数据库内容。
func (kv *PebbleKV) NewIterator(opt IterOption) (*pebble.Iterator, error) {
	var ro *pebble.IterOptions
	if len(opt.Prefix) > 0 {
		upper := nextPrefix(opt.Prefix)
		ro = &pebble.IterOptions{LowerBound: opt.Prefix, UpperBound: upper}
	}
	it, err := kv.PebbleDB.NewIter(ro)
	if err != nil {
		return nil, err
	}
	if ro != nil && ro.LowerBound != nil {
		it.SeekGE(ro.LowerBound)
	} else {
		it.First()
	}
	return it, nil
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
