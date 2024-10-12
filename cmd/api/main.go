package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error(".env missing")
	}
	mux := http.NewServeMux()

	ping := func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		fmt.Fprintf(w, "Thanks, %s, ping has worked", name)
	}

	mux.HandleFunc("GET /ping/{name}", ping)

	srv := &http.Server{
		Addr:         ":4000",
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err = srv.ListenAndServe()
	if err != nil {
		os.Exit(1)
	}
}
