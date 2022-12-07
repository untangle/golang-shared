package util

import "sync"

// Singleton implements a singleton. Allows you to wrap an object up
// in a singleton with less fuss.
//
// example:
// var myInstance *Singleton = NewSingleton(func() interface{ return NewInstance() })
// func GetMyInstance() *MyInstance {
//     return myInstance.GetInstance().(*MyInstance)
// }
//
// The GetInstance() function is threadsafe.
type Singleton struct {
	instance    interface{}
	once        sync.Once
	constructor func() interface{}
}

// GetInstance -- returns the value of constructor(), which only gets
// called once.
func (singleton *Singleton) GetInstance() interface{} {
	if singleton.instance != nil {
		return singleton.instance
	}
	singleton.once.Do(
		func() {
			singleton.instance =
				singleton.constructor()

		})
	if singleton.instance == nil {
		panic("Singleton: No instance available!")
	}
	return singleton.instance
}

// SetInstance -- workaround method for if you can't for whatever
// reason just create a singleton. This should generally not be used
// unless your singleton needs for example some variables that get
// created later, that you can't access from a constructor function.
func (singleton *Singleton) SetInstance(instance interface{}) {
	singleton.instance = instance
}

// NewSingleton -- new singleton. constructor is a func returning the
// instance GetInstance() returns as an interface.
func NewSingleton(
	constructor func() interface{}) *Singleton {
	return &Singleton{
		instance:    nil,
		constructor: constructor,
	}
}

// NewSingletonNoConstructor -- returns a singleton without a constructor.
func NewSingletonNoConstructor() *Singleton {
	return &Singleton{
		constructor: func() interface{} { return nil },
	}
}
