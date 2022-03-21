package server

import (
	"ConsensusTime/voting"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func AddApiRoutes(router *mux.Router) {
	// POST request to vote
	router.HandleFunc("/vote", func(w http.ResponseWriter, r *http.Request) {
		vote, err := voting.NewVoteFromJSON(json.NewDecoder(r.Body))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		voting.AddVote(vote)
		fmt.Fprintf(w, "You as %s guessed that the current time was %s (this was off by %s compared to UTC time)", *vote.User, vote.GuessedDate, vote.Difference)
	}).Methods("POST")

	// GET request to see all votes
	router.HandleFunc("/votes", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(voting.Votes())
	}).Methods("GET")
}
