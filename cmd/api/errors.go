package main

import (
	"bookshop/internal/httputil"
	"log/slog"
	"net/http"
)

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	if err := httputil.WriteJSON(w, status, message, nil); err != nil {
		app.logError("Failed to write response", r.Method, r.URL.String(), err)

		w.WriteHeader(http.StatusInternalServerError)

		fallbackResponse := `{"message": "The server encountered an error and could not process your request"}`
		_, _ = w.Write([]byte(fallbackResponse))
	}
}

type errorResponse struct {
	ErrMsg string `json:"errMsg"`
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	msg := "Internal server error"
	app.logError(msg, r.Method, r.URL.String(), err)
	app.errorResponse(w, r, http.StatusInternalServerError, errorResponse{msg})
}

func (app *application) logError(msg, reqMethod, reqURL string, err error) {
	app.logger.Error(
		msg,
		slog.String("method", reqMethod),
		slog.String("path", reqURL),
		slog.String("error", err.Error()),
	)
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, errorResponse{err.Error()})
}

type validationResponse struct {
	Errors map[string]string `json:"errors"`
}

func (app *application) validationErrorResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, validationResponse{Errors: errors})
}
