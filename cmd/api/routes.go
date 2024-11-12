package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", app.healthCheckHandler)

	mux.HandleFunc("POST /v1/authors", app.createAuthorHandler)

	mux.HandleFunc("POST /v1/users", app.registerUserHandler)
	mux.HandleFunc("PUT /v1/users/activate", app.activateUserHandler)
	mux.HandleFunc("POST /v1/users/activation/resend-token", app.resendActivationTokenHandler)

	mux.HandleFunc("GET /v1/doc/", app.openAPISpecHandler)
	mux.HandleFunc("GET /v1/doc/ui", app.redocHandler)

	mux.HandleFunc("POST /v1/books/categories", app.createBooksCategoryHandler)

	return mux
}
