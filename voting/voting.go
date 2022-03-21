package voting

import (
	"encoding/json"
	"errors"
	"sort"
	"time"
)

var defaultUser string = "Anonymous"

type Vote struct {
	User        *string    `json:"user"`
	GuessedTime *time.Time `json:"time"`
	ActualTime  time.Time
	Difference  time.Duration
}

var votes = []Vote{}

func NewVote(user *string, guessedTime *time.Time) (*Vote, error) {
	if user == nil {
		user = &defaultUser
	}

	if guessedTime == nil {
		return nil, errors.New("Can not create vote. A vote needs the time you are guessing.")
	}

	actualTime := time.Now().UTC()
	return &Vote{user, guessedTime, actualTime, actualTime.Sub(guessedTime.UTC())}, nil
}

func NewVoteFromJSON(decoder *json.Decoder) (*Vote, error) {
	var vote Vote
	if err := decoder.Decode(&vote); err != nil {
		return nil, err
	}

	if vote.GuessedTime == nil {
		return nil, errors.New("You need to specify the time you are guessing in the 'time' field of your JSON body.")
	}

	return NewVote(vote.User, vote.GuessedTime)
}

func AddVote(vote *Vote) {
	if vote == nil {
		return
	}

	votes = append(votes, *vote)
}

func Votes() []Vote {
	return votes
}

func VotesByUser(user string) []Vote {
	votesByUser := []Vote{}
	for _, vote := range votes {
		if *vote.User == user {
			votesByUser = append(votesByUser, vote)
		}
	}

	return votesByUser
}

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
