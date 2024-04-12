package atomicbool

import (
	"sync/atomic"
)

// AtomicBool is a boolean type that supports atomic operations.
type AtomicBool struct {
	flag int32
}

// NewAtomicBool initializes a new AtomicBool with the specified initial value.
func NewAtomicBool(initialValue bool) *AtomicBool {
	var flag int32
	if initialValue {
		flag = 1
	}
	return &AtomicBool{flag: flag}
}

// Set atomically sets the value of the boolean flag.
func (b *AtomicBool) Set(value bool) {
	var newValue int32
	if value {
		newValue = 1
	}
	atomic.StoreInt32(&b.flag, newValue)
}

// Get atomically retrieves the current value of the boolean flag.
func (b *AtomicBool) Get() bool {
	return atomic.LoadInt32(&b.flag) != 0
}
