package kvs

import (
	"errors"

	"github.com/cockroachdb/pebble"
	"github.com/lance4117/gofuse/errs"
	"github.com/lance4117/gofuse/logger"
	"github.com/lance4117/gofuse/once"
)

type PebbleConfig struct {
	DirName string
	*pebble.Options
}

type PebbleKV struct {
	db      *pebble.DB
	options *pebble.WriteOptions
}

func DefaultPebbleConfig(dirname string) PebbleConfig {
	opts := &pebble.Options{}
	opts.EnsureDefaults()
	return PebbleConfig{dirname, opts}
}

var NewPebbleSync = once.DoWithParam(func(config PebbleConfig) *PebbleKV {
	// 获取实例
	kv, err := pebble.Open(config.DirName, config.Options)
	if err != nil {
		logger.Fatal(err, errs.ErrNewStoreEngineFail)
		return nil
	}
	return &PebbleKV{kv, pebble.Sync}
})

var NewPebbleNoSync = once.DoWithParam(func(config PebbleConfig) *PebbleKV {
	// 获取实例
	kv, err := pebble.Open(config.DirName, config.Options)
	if err != nil {
		logger.Fatal(err, errs.ErrNewStoreEngineFail)
		return nil
	}
	return &PebbleKV{kv, pebble.NoSync}
})

func (kv *PebbleKV) Put(key string, val []byte) error {
	return kv.db.Set([]byte(key), val, kv.options)
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
	return kv.db.Delete([]byte(key), kv.options)
}
