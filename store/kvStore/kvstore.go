package kvStore

type KVStore interface {
	Put(key string, val []byte) error
	Get(key string) ([]byte, error)
	Del(key string) error
}
