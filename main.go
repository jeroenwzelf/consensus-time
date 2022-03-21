package main

import (
	"ConsensusTime/server"
	"ConsensusTime/voting"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

func NewIndexHtmlTemplateExecutor() *template.Template {
	return template.Must(template.New("index.gohtml").Funcs(template.FuncMap{
		"GetConsensusDateDifferenceMillis": func() int64 {
			return int64(voting.GetConsensusTimeDifference().Milliseconds())
		},
	}).ParseFiles("html/index.gohtml"))
}

func router() *mux.Router {
	router := server.NewDefaultRouter()
	templateExecutor := NewIndexHtmlTemplateExecutor()

	// Root shows static html template `index.gohtml`
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := templateExecutor.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

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

	return router
}

func flags() (time.Duration, string) {
	var gracefulTimeoutDuration time.Duration
	flag.DurationVar(&gracefulTimeoutDuration, "graceful-timeout", time.Second*15,
		"the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")

	var port string
	flag.StringVar(&port, "port", "8081",
		"the port which the server runs on locally")

	flag.Parse()
	return gracefulTimeoutDuration, port
}

func main() {
	gracefulTimeoutDuration, port := flags()

	serverInstance := server.RunServer(port, router())

	// Wait until SIGINT received
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Gracefully stop (until gracefulTimeoutDuration is reached)
	ctx, cancel := context.WithTimeout(context.Background(), gracefulTimeoutDuration)
	defer cancel()
	serverInstance.Shutdown(ctx)
	os.Exit(0)
}
