package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", app.healthCheckHandler)

	mux.HandleFunc("POST /v1/authors", app.createAuthorHandler)

	mux.HandleFunc("GET /v1/doc/", app.openAPISpecHandler)
	mux.HandleFunc("GET /v1/doc/ui", app.redocHandler)

	mux.HandleFunc("POST /v1/books/categories", app.createBooksCategoryHandler)

	return mux
}
