package voting

import (
	"sort"
	"time"
)

// Calculate Median Consensus Time by taking the median of the time difference between `GuessedTime` and `ActualTime` of all votes
func GetConsensusTime() time.Time {
	return time.Now().Add(GetConsensusTimeDifference())
}

func GetConsensusTimeDifference() time.Duration {
	if len(votes) == 0 {
		return 0
	}

	times := []time.Duration{}
	for _, vote := range votes {
		times = append(times, vote.Difference)
	}

	// Sort UTCTime fields by time
	sort.Slice(times, func(i, j int) bool { return times[i] < times[j] })

	// Get the median difference
	return times[len(times)/2]
}
