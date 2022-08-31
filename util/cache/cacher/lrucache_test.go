package cacher

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/suite"
)

type LruCacheTestSuite struct {
	suite.Suite
	cache     LruCache
	capacity  uint
	cacheName string
}

// Initialize a cache and max out its capacity
func (suite *LruCacheTestSuite) SetupTest() {
	suite.capacity = 5
	suite.cacheName = "LRUUnitTest"
	suite.cache = *NewLruCache(suite.capacity, suite.cacheName)
	for i := 0; i < int(suite.capacity); i++ {
		suite.cache.Put(strconv.Itoa(int(i)), i)
	}
}

func TestLruCacheTestSuite(t *testing.T) {
	suite.Run(t, new(LruCacheTestSuite))
}

// Tests getting the most recently fetched element from the cache
func (suite *LruCacheTestSuite) TestGetMostRecentlyUsed() {
	expectedKey, expectedValue := "4", 4
	key, value := suite.cache.GetMostRecentlyUsed()

	suite.Equal(expectedKey, key)
	suite.Equal(expectedValue, value)
}

// Tests getting the least recently used element from the cache
func (suite *LruCacheTestSuite) TestGetLeastRecentlyUsed() {
	expectedKey, expectedValue := "0", 0
	key, value := suite.cache.GetLeastRecentlyUsed()

	suite.Equal(expectedKey, key)
	suite.Equal(expectedValue, value)
}

// Tests removing an element from the cache
func (suite *LruCacheTestSuite) TestRemove() {
	toRemove := "2"
	_, ok := suite.cache.Get(toRemove)
	suite.True(ok, "The element with key %s to be removed from the cache is not in the cache", toRemove)

	suite.cache.Remove(toRemove)

	_, okAfterRemoval := suite.cache.Get(toRemove)
	suite.False(okAfterRemoval, "The element with key %s was not removed from the cache", toRemove)
}

// Check if the values the cache was initialized with can be retrieved
func (suite *LruCacheTestSuite) TestGet() {
	for i := 0; i < int(suite.capacity); i++ {
		value, ok := suite.cache.Get(strconv.Itoa(int(i)))
		suite.True(ok, "The key %d did not exist in the cache", i)

		suite.Equal(value, i)

		// Check element was moved to the front of the linked-list
		key, value := suite.cache.GetMostRecentlyUsed()
		suite.Equal(strconv.Itoa(int(i)), key)
		suite.Equal(i, value)
	}
}

// Tests getting the total number of elements in the cache
func (suite *LruCacheTestSuite) TestGetCurrentCapacity() {
	suite.Equal(int(suite.capacity), suite.cache.GetTotalElements())
}

// Tests adding an element to the cache when the cache is at capacity
func (suite *LruCacheTestSuite) TestCapacityExceeded() {
	// Check that the cache has something in it to start with

	// The first element put in the cache is the least recently used element
	// so adding more elements should delete it from the queue
	toRemove := "0"

	suite.cache.Put(strconv.Itoa(int(suite.capacity)), suite.capacity)

	_, okAfterOverwritten := suite.cache.Get(toRemove)
	suite.False(okAfterOverwritten, "The element with key %s was not overwritten in the cache", toRemove)
}

// Tests writing to a value already in the cache
func (suite *LruCacheTestSuite) TestUpdatingCacheValue() {
	toUpdate := "2"
	updatedValue := 10

	_, ok := suite.cache.Get(toUpdate)
	suite.True(ok, "The element with key %s to be updated in the cache is not in the cache", toUpdate)

	suite.cache.Put(toUpdate, updatedValue)

	// Check value was updated
	value, _ := suite.cache.Get(toUpdate)
	suite.Equal(updatedValue, value)

	// Check element was moved to the front of the linked-list
	key, value := suite.cache.GetMostRecentlyUsed()
	suite.Equal(toUpdate, key)
	suite.Equal(updatedValue, value)
}

// Tests clearing the cache
func (suite *LruCacheTestSuite) TestClear() {
	// Check that the cache has something in it to start with
	suite.Equal(int(suite.capacity), suite.cache.GetTotalElements(), "The cache is missing elements. It was not setup properly by SetupTest()")

	suite.cache.Clear()

	suite.Equal(0, suite.cache.GetTotalElements(), "The cache was not successfully cleared")
}
