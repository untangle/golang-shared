package sweeper

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type SweepOnTimeTestSuite struct {
	suite.Suite
	sweeper       SweepOnTime
	sweepInterval time.Duration

	// Cache used to test if the cleanup task gets run
	cache map[string]int

	// A value used as sweep criteria
	cacheMaxValue int
}

func (suite *SweepOnTimeTestSuite) getCleanupFunc() func() {
	cleanupCache := &(suite.cache)
	elementMaxValue := suite.cacheMaxValue

	// Simple function to remove a cache element if it's over a value
	return func() {
		for key, val := range *cleanupCache {
			if val < elementMaxValue {
				delete(*cleanupCache, key)
			}
		}
	}
}

func TestSweepOnTimeTestSuite(t *testing.T) {
	suite.Run(t, new(SweepOnTimeTestSuite))
}

func (suite *SweepOnTimeTestSuite) SetupTest() {
	suite.sweepInterval = 5 * time.Millisecond
	suite.cache = map[string]int{"1": 1, "2": 2, "3": 3, "4": 4}
	suite.cacheMaxValue = 3
	suite.sweeper = *NewSweepOnTime(suite.sweepInterval)

	suite.sweeper.StartSweeping(suite.getCleanupFunc())
}

func (suite *SweepOnTimeTestSuite) TestSwept() {
	time.Sleep((suite.sweepInterval) + 1)

	for key, val := range suite.cache {
		fmt.Println(key)
		suite.True(val >= suite.cacheMaxValue, "The test cache was not swept")
	}
}

func (suite *SweepOnTimeTestSuite) TearDownTest() {
	suite.sweeper.StopSweeping()
}
