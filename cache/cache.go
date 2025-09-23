package cache

import (
	"reflect"
	"time"

	"github.com/allegro/bigcache"
	"github.com/lance4117/gofuse/codec"
	"github.com/lance4117/gofuse/errs"
)

type Cache struct {
	BigCache *bigcache.BigCache
}

// NewCache 初始化带时间过期的缓存,0为永不过期
func NewCache(expTime time.Duration) (*Cache, error) {
	cache, err := bigcache.NewBigCache(bigcache.DefaultConfig(expTime))
	return &Cache{cache}, err
}

// Set 将指定的key存储到缓存中
func (i *Cache) Set(key string, value any) error {
	marshal, err := codec.MPMarshal(value)
	if err != nil {
		return err
	}
	return i.BigCache.Set(key, marshal)
}

// Get 传入key 和 v指针，返回相应类型
func (i *Cache) Get(key string, v any) error {
	// 检查 v 是否为指针类型
	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		return errs.ErrNeedPointer
	}
	bytes, err := i.BigCache.Get(key)
	if err != nil {
		return err
	}
	return codec.MPUnmarshal(bytes, v)
}

// Delete 删除
func (i *Cache) Delete(key string) error {
	return i.BigCache.Delete(key)
}

// Len 当前缓存长度
func (i *Cache) Len() int {
	return i.BigCache.Len()
}

// Reset 重置缓存
func (i *Cache) Reset() error {
	return i.BigCache.Reset()
}

// Close 关闭缓存
func (i *Cache) Close() error {
	return i.BigCache.Close()
}
