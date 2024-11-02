package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", app.healthCheckHandler)

	mux.HandleFunc("POST /v1/authors", app.createAuthorHandler)

	return mux
}
