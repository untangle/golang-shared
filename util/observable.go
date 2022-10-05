package util

import "sync"

// Listener type for geoip. It is called with the session and an
// interface.
type Listener func(changed ...interface{})

// ObservableInterface is an interface all Observables should adhere
// to.
type ObservableInterface interface {
	AddListener(Listener)
	NotifyListeners(...interface{})
}

// Observable type has a list of Listener objects that are called with
// a given session and 'extra information' put into an interface{}
// object when NotifyListeners are called. Call AddListener to add a
// listener.
type Observable struct {
	listeners []Listener
}

// AddListener Adds a Listener to the list of listeners that will get called by
// NotifyListeners.
func (observable *Observable) AddListener(listener Listener) {
	observable.listeners = append(observable.listeners, listener)
}

// NotifyListeners will loop through all listeners in the Observable
// and notify them -- calling them with the given session and
// 'changed' object. The type of changed depends on who is notifying.
func (observable *Observable) NotifyListeners(event ...interface{}) {
	for _, listener := range observable.listeners {
		listener(event...)
	}
}

// ThreadsafeObservable is an observable that is threadsafe. It works
// the same as Observable but wraps operations in a lock.
type ThreadsafeObservable struct {
	Observable
	lock sync.Mutex
}

// AddListener adds a listener to the list of listeners. This listener
// function gets called each time NotifyListeners() is called.
func (observable *ThreadsafeObservable) AddListener(listener Listener) {
	observable.lock.Lock()
	defer observable.lock.Unlock()
	observable.Observable.AddListener(listener)
}

// NotifyListeners notifies all listeners of the event. It makes a
// copy of the current listener list and calls all listeners on the
// copy list. This makes it so that we only have to hold the lock for
// as long as it takes to copy.
func (observable *ThreadsafeObservable) NotifyListeners(event ...interface{}) {
	observable.lock.Lock()
	listeners := make([]Listener, len(observable.listeners))
	copy(listeners, observable.listeners)
	observable.lock.Unlock()

	for _, listener := range listeners {
		listener(event...)
	}
}
