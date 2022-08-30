package cacher

// Interface for structs implementing basic caching functionality
type Cacher interface {
	Get(string) (interface{}, bool)
	Put(string, interface{})
	Clear()
	Remove(string)

	// An iterator
	//GetIterator() func() (string, interface{}, bool)

	// Runs a given function on each cache element
	ForEach(func(string, interface{}) bool)
}
