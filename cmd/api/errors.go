package main

import (
	"bookshop/internal/utils"
	"log/slog"
	"net/http"
)

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
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

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

type validationErrors struct {
	Errors map[string]string `json:"errors"`
}

func (app *application) validationErrorResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, validationErrors{Errors: errors})
}
