package cache

type CacheInterface interface {
	Get() (interface{}, bool)
	Put()
	Clear()
	Remove()
}
