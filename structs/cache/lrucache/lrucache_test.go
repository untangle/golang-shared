package lrucache

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var capacity uint = 5

type LruCacheTestSuite struct {
	suite.Suite
	cache  LruCache
	assert assert.Assertions
}

// Initialize a cache and max out its capacity
func (suite *LruCacheTestSuite) SetupTest() {
	suite.cache = *NewLruCache(capacity)
	for i := 0; i < int(capacity); i++ {
		suite.cache.Put(strconv.Itoa(int(i)), i)
	}
}

func (suite *LruCacheTestSuite) TestGetMostRecentlyUsed() {
	expectedKey, expectedValue := "4", 4
	key, value := suite.cache.GetMostRecentlyUsed()

	assert.Equal(suite.T(), expectedKey, key)
	assert.Equal(suite.T(), expectedValue, value)
}

func (suite *LruCacheTestSuite) TestGetLeastRecentlyUsed() {
	expectedKey, expectedValue := "0", 0
	key, value := suite.cache.GetLeastRecentlyUsed()

	assert.Equal(suite.T(), expectedKey, key)
	assert.Equal(suite.T(), expectedValue, value)
}

func (suite *LruCacheTestSuite) TestRemove() {
	toRemove := "2"
	_, ok := suite.cache.Get(toRemove)
	assert.True(suite.T(), ok, "The element with key %s to be removed from the cache is not in the cache", toRemove)

	suite.cache.Remove(toRemove)

	_, okAfterRemoval := suite.cache.Get(toRemove)
	assert.False(suite.T(), okAfterRemoval, "The element with key %s was not removed from the cache", toRemove)
}

// Check if the values the cache was initialized with can be retrieved
func (suite *LruCacheTestSuite) TestGet() {
	for i := 0; i < int(capacity); i++ {
		value, ok := suite.cache.Get(strconv.Itoa(int(i)))
		assert.True(suite.T(), ok, "The key %d did not exist in the cache", i)

		if ok {
			assert.Equal(suite.T(), value, i)
		}
	}
}

func (suite *LruCacheTestSuite) TestGetCurrentCapacity() {
	assert.Equal(suite.T(), int(capacity), suite.cache.GetCurrentCapacity())
}

func (suite *LruCacheTestSuite) TestCapacityExceeded() {
	// Check that the cache has something in it to start with

	// The first element put in the cache is the least recently used element
	// so adding more elements should delete it from the queue
	toRemove := "0"

	suite.cache.Put(strconv.Itoa(int(capacity)), capacity)

	_, okAfterOverwritten := suite.cache.Get(toRemove)
	assert.False(suite.T(), okAfterOverwritten, "The element with key %s was not overwritten in the cache", toRemove)
}

func (suite *LruCacheTestSuite) TestClear() {
	// Check that the cache has something in it to start with
	assert.Equal(suite.T(), int(capacity), suite.cache.GetCurrentCapacity(), "The cache is missing elements. It was not setup properly by SetupTest()")

	suite.cache.Clear()

	assert.Equal(suite.T(), 0, suite.cache.GetCurrentCapacity(), "The cache was not successfully cleared")
}

func TestLruCacheTestSuite(t *testing.T) {
	suite.Run(t, new(LruCacheTestSuite))
}