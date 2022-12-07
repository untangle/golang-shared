package util

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSingleton(t *testing.T) {
	type MyObject struct{}

	singleton := NewSingleton(func() interface{} { return &MyObject{} })
	inst := singleton.GetInstance()
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			myLocalInst := singleton.GetInstance()
			assert.Same(t, inst, myLocalInst)
			wg.Done()
		}()
	}
	wg.Wait()
}
