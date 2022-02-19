package main

import (
	"log"
	"net/http"
	"time"

	"github.com/ejuju/video-streaming-server/internal/video"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", video.ServeFromLocalFile)

	server := http.Server{
		Addr:              ":8080",
		Handler:           router,
		IdleTimeout:       5 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
