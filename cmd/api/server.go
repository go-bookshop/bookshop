package main

import (
	"net/http"
	"time"

	"github.com/rs/cors"
)

func (app *application) configure() http.Handler {
	routes := app.routes()

	corsOptions := cors.New(cors.Options{
		AllowedOrigins:   app.config.cors.allowedOrigins,
		AllowedHeaders:   app.config.cors.allowedHeaders,
		AllowedMethods:   app.config.cors.allowedMethods,
		AllowCredentials: true,
	})

	handler := corsOptions.Handler(routes)

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
