package cache

type Cache interface {
	Get()
	Put()
	Remove()
	Clear()
	New()
}
