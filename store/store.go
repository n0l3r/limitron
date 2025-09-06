package store

type Store interface {
	Get(key string) (int64, error)
	Set(key string, value int64, expiry int64) error
	Incr(key string, n int64, expiry int64) (int64, error)
}
