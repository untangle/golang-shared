package cache

type CacheInterface interface {
	Get(string) (interface{}, bool)
	Put(string, interface{})
	Clear()
	Remove(string)
}
