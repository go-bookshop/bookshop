package main

import (
	"net/http"
	"time"
)

func (app *application) run() error {
	srv := &http.Server{
		Addr:         ":4000",
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app.logger.Info("starting server", "addr", srv.Addr)

	return srv.ListenAndServe()
}
