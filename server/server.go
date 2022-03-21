package server

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func NewDefaultRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	// Redirect any other route to root
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	return router
}

func RunServer(port string, router *mux.Router) *http.Server {
	// Configure server
	hostAddr := net.JoinHostPort("127.0.0.1", port)
	server := &http.Server{
		Addr:         hostAddr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 65,
		Handler:      router,
	}

	// Start server (non-blocking)
	go func() {
		log.Println("Running server on " + hostAddr)
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	return server
}
