package kvs

import (
	"errors"

	"github.com/cockroachdb/pebble"
	"github.com/lance4117/gofuse/errs"
	"github.com/lance4117/gofuse/logger"
)

type PebbleConfig struct {
	DirName string
	*pebble.Options
}

type PebbleKV struct {
	db *pebble.DB
}

type PebbleBatch struct {
	*pebble.Batch
}

type PebbleIterator struct {
	it *pebble.Iterator
}

func DefaultPebbleConfig(dirname string) PebbleConfig {
	opts := &pebble.Options{}
	opts.EnsureDefaults()
	return PebbleConfig{dirname, opts}
}

func NewPebbleKV(config PebbleConfig) *PebbleKV {
	// 获取实例
	kv, err := pebble.Open(config.DirName, config.Options)
	if err != nil {
		logger.Fatal(err, errs.ErrNewStoreEngineFail)
		return nil
	}
	return &PebbleKV{kv}
}

func (kv *PebbleKV) Put(key string, val []byte) error {
	return kv.db.Set([]byte(key), val, pebble.NoSync)
}

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
func (kv *PebbleKV) Del(key string) error {
	return kv.db.Delete([]byte(key), pebble.NoSync)
}

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

func (kv *PebbleKV) PutSync(key string, val []byte) error {
	return kv.db.Set([]byte(key), val, pebble.Sync)
}

func (kv *PebbleKV) DelSync(key string) error {
	return kv.db.Delete([]byte(key), pebble.Sync)
}

func (kv *PebbleKV) NewBatch() (Batch, error) {
	batch := kv.db.NewBatch()
	return &PebbleBatch{batch}, nil

}

func (kv *PebbleKV) Close() error {
	return kv.db.Close()
}

func (b *PebbleBatch) Put(key string, val []byte) error {
	return b.Batch.Set([]byte(key), val, pebble.NoSync)
}

func (b *PebbleBatch) Del(key string) error {
	return b.Batch.Delete([]byte(key), pebble.NoSync)
}

func (b *PebbleBatch) PutSync(key string, val []byte) error {
	return b.Batch.Set([]byte(key), val, pebble.Sync)
}

func (b *PebbleBatch) DelSync(key string) error {
	return b.Batch.Delete([]byte(key), pebble.Sync)
}

func (b *PebbleBatch) Commit() error {
	return b.Batch.Commit(pebble.Sync)
}

func (b *PebbleBatch) Cancel() error {
	return b.Batch.Close()
}

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

func (i *PebbleIterator) Valid() bool {
	return i.it.Valid()
}

func (i *PebbleIterator) Next() bool {
	return i.it.Next()
}

func (i *PebbleIterator) Key() []byte {
	k := i.it.Key()
	out := make([]byte, len(k))
	copy(out, k)
	return out
}

func (i *PebbleIterator) Value() ([]byte, error) {
	v, err := i.it.ValueAndErr()
	if err != nil {
		return nil, err
	}
	out := make([]byte, len(v))
	copy(out, v)
	return out, nil
}

func (i *PebbleIterator) Close() error {
	return i.it.Close()
}

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
