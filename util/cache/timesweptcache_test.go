package cache

import (
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type MockCache struct {
	elements  map[string]interface{}
	cacheLock sync.RWMutex
}

func (cache *MockCache) Get(key string) (interface{}, bool) {
	cache.cacheLock.RLock()
	defer cache.cacheLock.RUnlock()
	val, ok := cache.elements[key]
	return val, ok
}

func (cache *MockCache) Put(key string, value interface{}) {
	cache.cacheLock.Lock()
	defer cache.cacheLock.Unlock()
	cache.elements[key] = value
}

func (cache *MockCache) Clear() {
	cache.cacheLock.Lock()
	defer cache.cacheLock.Unlock()
	cache.elements = make(map[string]interface{})
}

func (cache *MockCache) Remove(key string) {
	cache.cacheLock.Lock()
	defer cache.cacheLock.Unlock()
	cache.removeNoLock(key)
}

func (cache *MockCache) removeNoLock(key string) {
	delete(cache.elements, key)
}

func (cache *MockCache) ForEach(cleanUp func(string, interface{}) bool) {
	cache.cacheLock.Lock()
	defer cache.cacheLock.Unlock()
	for key, val := range cache.elements {
		if cleanUp(key, val) {
			cache.removeNoLock(key)
		}
	}
}

func NewMockCache() *MockCache {
	return &MockCache{elements: make(map[string]interface{})}
}

type TestTimeSweptCache struct {
	suite.Suite
	timeSweptCache TimeSweptCache

	sweepInterval int
}

func TestTimeSweptCacheSuite(t *testing.T) {
	suite.Run(t, &TestTimeSweptCache{})
}

func (suite *TestTimeSweptCache) SetupTest() {
	suite.sweepInterval = 1
	suite.timeSweptCache = *NewTimeSweptCache(NewMockCache(), time.Duration(suite.sweepInterval))

	for i := 0; i < 5; i++ {
		suite.timeSweptCache.Put(strconv.Itoa(int(i)), i)
	}
}

// Test that the provided cleanup function to StartSweeping() is being ran
func (suite *TestTimeSweptCache) TestCleanupTaskRan() {

	// Remove elements that aren't equal to 4
	suite.timeSweptCache.StartSweeping(func(s string, i interface{}) bool {
		deleteElement := false
		if i.(int) != 4 {
			deleteElement = true
		}

		return deleteElement
	})

	time.Sleep((time.Duration(suite.sweepInterval) + 1) * time.Second)

	_, ok := suite.timeSweptCache.Get("1")
	suite.False(ok, "The cleanup task was not run as expected")

	_, ok = suite.timeSweptCache.Get("4")
	suite.True(ok, "The cleanup task removed an unexpected cache element")

}

func (suite *TestTimeSweptCache) TearDownTest() {
	suite.timeSweptCache.StopSweeping()
}
