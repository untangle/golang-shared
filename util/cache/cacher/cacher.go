package cacher

// Interface for structs implementing a basic caching functionality
type Cacher interface {
	Get(string) (interface{}, bool)
	Put(string, interface{})
	Clear()
	Remove(string)

	// Runs a given function on each cache element
	ForEach(func(string, interface{}) bool)
}
