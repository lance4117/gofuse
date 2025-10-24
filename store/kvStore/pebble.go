package kvStore

import (
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
	var options *pebble.Options
	options = options.EnsureDefaults()
	return PebbleConfig{dirname, options}
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
	get, closer, err := kv.db.Get([]byte(key))
	defer closer.Close()
	return get, err
}
func (kv *PebbleKV) Del(key string) error {
	return kv.db.Delete([]byte(key), kv.options)
}
