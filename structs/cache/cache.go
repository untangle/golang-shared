package cache

type CacheInterface interface {
	Get(string) (interface{}, bool)
	Put()
	Clear()
	Remove()
}
