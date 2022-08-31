package cacher

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RRCacheTestSuite struct {
	suite.Suite
	cache     RandomReplacementCache
	capacity  uint
	cacheName string
}

func TestRRCacheTestSuite(t *testing.T) {
	suite.Run(t, new(RRCacheTestSuite))
}

// Initialize a cache and max out its capacity
func (suite *RRCacheTestSuite) SetupTest() {
	suite.capacity = 5
	suite.cacheName = "RRUnitTest"
	suite.cache = *NewRandomReplacementCache(suite.capacity, suite.cacheName)

	for i := 0; i < int(suite.capacity); i++ {
		suite.cache.Put(strconv.Itoa(int(i)), i)
	}
}

// The suite is not being used since
// a pointer is required as a cache value
// if any alteration to the value are going to
// be successful
func TestForEachElementMutation(t *testing.T) {
	capacity := 5
	cacheName := "cacheMutationTest"
	testCache := *NewRandomReplacementCache(uint(capacity), cacheName)

	for i := 0; i < capacity; i++ {
		// Create a copy of i so all elements don't
		// point to the same int
		newVal := i
		testCache.Put(strconv.Itoa(int(i)), &newVal)
	}

	assert.Equal(t, testCache.GetCurrentCapacity(), capacity)

	mutateElement := func(s string, i interface{}) bool {
		deleteElement := false
		if *(i).(*int) != 3 {
			*(i).(*int) = 3
		}

		return deleteElement
	}
	testCache.ForEach(mutateElement)

	for key, val := range testCache.elements {
		assert.Equal(t, 3, *(val.value.(*int)), "The key %s was not altered as expected", key)
	}
}

func (suite *RRCacheTestSuite) TestForEachElementDeletion() {
	suite.cache.ForEach(func(s string, i interface{}) bool {
		deleteElement := false
		if i.(int) < 4 {
			deleteElement = true
		}

		return deleteElement
	})

	for key, val := range suite.cache.elements {
		suite.Equal(4, val.value.(int), "The key %s was not altered as expected", key)
	}
}

// func (suite *RRCacheTestSuite) TestNextElement() {
// 	next := suite.cache.GetIterator()
// 	var key string
// 	var isNext bool

// 	count := 0
// 	for key, _, isNext = next(); isNext == true; key, _, isNext = next() {
// 		_, ok := suite.cache.elements[key]

// 		suite.True(ok, "The iterator retrieved a value not in the cache")
// 		count += 1
// 	}

// 	suite.Equal(suite.capacity, uint(count), "The iterator did not iterate over ever element in the cache")

// }

func (suite *RRCacheTestSuite) TestGet() {
	for i := 0; i < int(suite.capacity); i++ {
		value, ok := suite.cache.Get(strconv.Itoa(i))
		suite.True(ok, "The key %d did not exist in the cache", i)

		suite.Equal(value, i)
	}
}

func (suite *RRCacheTestSuite) TestUpdatingCacheValue() {
	toUpdate := "2"
	updatedValue := 10

	_, ok := suite.cache.Get(toUpdate)
	suite.True(ok, "The element with key %s to be updated in the cache is not in the cache", toUpdate)

	suite.cache.Put(toUpdate, updatedValue)

	// Check value was updated
	value, _ := suite.cache.Get(toUpdate)
	suite.Equal(updatedValue, value)
}

func (suite *RRCacheTestSuite) TestClear() {
	suite.Equal(int(suite.capacity), suite.cache.GetCurrentCapacity(), "The cache is missing elements. It was not setup properly by SetupTest()")

	suite.cache.Clear()

	suite.Equal(0, suite.cache.GetCurrentCapacity(), "The cache was not successfully cleared")
}

func (suite *RRCacheTestSuite) TestCapacityExceeded() {
	keysBeforePut := getMapKeys(&suite.cache.elements)
	newKey := "6"
	newVal := 6

	suite.cache.Put("6", 6)
	keysAfterPut := getMapKeys(&suite.cache.elements)

	suite.Equal(int(suite.capacity), len(suite.cache.elements))

	val, ok := suite.cache.Get("6")
	suite.True(ok, "The cache did not contain the newly added value with key %s", newKey)
	suite.Equal(newVal, val, "The key %s did not have the expected value of %d", newKey, newVal)

	suite.NotEqual(keysAfterPut, keysBeforePut)
}

func (suite *RRCacheTestSuite) TestGetCurrentCapacity() {
	suite.Equal(int(suite.capacity), suite.cache.GetCurrentCapacity())
}

func (suite *RRCacheTestSuite) TestRemove() {
	keyToRemove := "2"
	_, ok := suite.cache.Get(keyToRemove)
	suite.True(ok, "The key -- %s -- going to be removed wasn't in the cache at the start of the test", keyToRemove)

	suite.cache.Remove("2")

	_, ok = suite.cache.Get(keyToRemove)
	suite.False(ok, "The key -- %s -- remained in the cache after being removed", keyToRemove)
	suite.Equal(suite.cache.maxCapacity-1, suite.cache.totalElements)
}
