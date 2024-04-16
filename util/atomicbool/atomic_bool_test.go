package atomicbool

import (
	"sync"
	"testing"
)

func TestAtomicBoolConcurrent(t *testing.T) {
	ab := NewAtomicBool(false)

	numGoroutines := 10000

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			for j := 0; j < 1000; j++ {
				ab.Set(!ab.Get())
			}
			wg.Done()
		}()
	}

	wg.Wait()

	finalValue := ab.Get()
	if finalValue != true && finalValue != false {
		t.Errorf("Final value should be true or false, got %v", finalValue)
	}
}

func TestAtomicBool(t *testing.T) {
	ab := NewAtomicBool(false)

	if ab.Get() != false {
		t.Errorf("Initial value should be false, got %v", ab.Get())
	}

	ab.Set(true)

	if ab.Get() != true {
		t.Errorf("Value should be true after setting, got %v", ab.Get())
	}

	ab.Set(false)

	if ab.Get() != false {
		t.Errorf("Value should be false after setting, got %v", ab.Get())
	}
}
