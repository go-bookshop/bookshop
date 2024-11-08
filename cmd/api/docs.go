package main

import (
	"bookshop/openapi"
	"fmt"
	"net/http"
)

func (app *application) openAPISpecHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(*app.config.doc.json)
	if err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("failed to write OpenAPI JSON response: %w", err))
	}
}

func (app *application) redocHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFileFS(w, r, openapi.Files, "ui/redoc.html")
}
