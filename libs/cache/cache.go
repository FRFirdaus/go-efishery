package cache

type CacheService interface {
	IsExists(key string) bool
	Get(key string, v interface{}) error
	Set(key string, v interface{}) error
	Delete(key string) error
	Close() error
}
