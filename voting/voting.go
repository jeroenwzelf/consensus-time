package voting

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

var defaultUser string = "Anonymous"
var formats = []string{"15:04:05", "15:04", time.Kitchen}

type VoteJSON struct {
	User        *string    `json:"user"`
	GuessedDate *time.Time `json:"date"`
	GuessedTime *string    `json:"time"`
}

type Vote struct {
	User        *string
	GuessedDate *time.Time
	ActualDate  time.Time
	Difference  time.Duration
}

var votes = []Vote{}

func NewVote(user *string, guessedDate *time.Time) (*Vote, error) {
	if user == nil {
		user = &defaultUser
	}

	if guessedDate == nil {
		return nil, errors.New("Can not create vote. A vote needs the time you are guessing.")
	}

	actualTime := time.Now().UTC()
	return &Vote{user, guessedDate, actualTime, guessedDate.UTC().Sub(actualTime)}, nil
}

func NewVoteFromJSON(decoder *json.Decoder) (*Vote, error) {
	var vote VoteJSON
	if err := decoder.Decode(&vote); err != nil {
		return nil, err
	}

	if vote.GuessedDate == nil {
		if vote.GuessedTime == nil {
			return nil, errors.New("You need to specify the time or date you are guessing in the 'time' or 'date' field of your JSON body. Examples:\n" +
				"\ttime: " + strings.Join(formats, " | ") + "\n" +
				"\tdate: 2022-03-21T01:04:40Z")
		}

		date, err := TimeStringToDate(*vote.GuessedTime, formats)
		if err != nil {
			return nil, err
		}
		vote.GuessedDate = FindNearestDateForTime(date)
	}

	return NewVote(vote.User, vote.GuessedDate)
}

func AddVote(vote *Vote) {
	if vote == nil {
		return
	}

	votes = append(votes, *vote)
}

func Clear() {
	votes = []Vote{}
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
