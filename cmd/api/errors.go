package main

import (
	"bookshop/internal/utils"
	"log/slog"
	"net/http"
)

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message string) {
	if err := utils.WriteJSON(w, status, message, nil); err != nil {
		app.logger.Error(
			"Failed to write response",
			slog.String("method", r.Method),
			slog.String("path", r.URL.String()),
			slog.String("error", err.Error()),
		)
		http.Error(w, "Failed to write response.", http.StatusInternalServerError)
	}
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	msg := "Internal server error"

	app.logger.Error(
		msg,
		slog.String("method", r.Method),
		slog.String("path", r.URL.String()),
		slog.String("error", err.Error()),
	)

	app.errorResponse(w, r, http.StatusInternalServerError, msg)
}
