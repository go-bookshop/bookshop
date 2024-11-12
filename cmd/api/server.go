package main

import (
	"net/http"
	"time"
)

func (app *application) configure() http.Handler {
	routes := app.routes()

	handler := app.enableCORS(routes)

	return handler
}

func (app *application) run() error {
	srv := &http.Server{
		Addr:         ":4000",
		Handler:      app.configure(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app.logger.Info("starting server", "addr", srv.Addr)

	return srv.ListenAndServe()
}
