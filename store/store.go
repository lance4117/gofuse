package store

type Store interface {
	Insert(beans ...any) (int64, error)
	Get(beans ...any) (bool, error)
	Ping() error
	Close() error
}
