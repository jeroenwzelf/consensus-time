package main

import (
	"ConsensusTime/server"
	"context"
	"flag"
	"os"
	"os/signal"
	"time"
)

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

	serverInstance := server.RunServer(port)

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
