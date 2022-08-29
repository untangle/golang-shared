package cache

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/untangle/golang-shared/structs/cache/cacher"
	"github.com/untangle/golang-shared/structs/cache/sweeper"
)

type TestSweptCache struct {
	suite.Suite
	sweptCache SweptCache

	sweepInterval int
}

func TestSweptCacheSuite(t *testing.T) {
	suite.Run(t, &TestSweptCache{})
}

func (suite *TestSweptCache) SetupTest() {
	suite.sweepInterval = 1
	suite.sweptCache = *NewSweptCache(cacher.NewRandomReplacementCache(5, "sweptCacheTest"), sweeper.NewSweepOnTime(time.Duration(suite.sweepInterval)*time.Second))

	for i := 0; i < 5; i++ {
		newVal := i
		suite.sweptCache.Put(strconv.Itoa(int(i)), &newVal)
	}
}

func (suite *TestSweptCache) TestElementDeletionFunction() {

	// Remove elements with a value less than 3. Run every second
	suite.sweptCache.StartSweeping(func(s string, i interface{}) bool {
		deleteElement := false
		if *(*(i.(*interface{}))).(*int) < 3 {
			deleteElement = true
		}

		return deleteElement
	})

	time.Sleep((time.Duration(suite.sweepInterval) + 1) * time.Second)

	next := suite.sweptCache.GetIterator()

	for key, val, ok := next(); ok; key, val, ok = next() {
		suite.True(*(*(val.(*interface{}))).(*int) >= 3, "The key %s was not swept as expected", key)
	}

}

func (suite *TestSweptCache) TestElementMutationFunction() {
	suite.sweptCache.StartSweeping(func(s string, i interface{}) bool {
		if *(*(i.(*interface{}))).(*int) < 4 {
			*(*(i.(*interface{}))).(*int) = 4
		}

		return false
	})

	time.Sleep((time.Duration(suite.sweepInterval) + 2) * time.Second)

	next := suite.sweptCache.GetIterator()
	for key, val, ok := next(); ok; key, val, ok = next() {

		suite.True(*(*(val.(*interface{}))).(*int) == 4, "The key %s was not altered as expected", key)
	}

}

func (suite *TestSweptCache) TearDownTest() {
	suite.sweptCache.StopSweeping()
}
