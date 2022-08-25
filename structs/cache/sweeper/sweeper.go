package sweeper

// Interface for altering/removing elements from a cache
// StartSweep should kickoff a goroutine that scans the a
// cache on a set interval.
type Sweeper interface {
	StartSweeping(func())
	StopSweeping()
}
