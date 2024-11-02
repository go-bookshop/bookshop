package main

import (
	"bookshop/internal/httputil"
	"log/slog"
	"net/http"
)

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message string) {
	if err := httputil.WriteJSON[any](w, status, message, nil, nil); err != nil {
		logError(app.logger, "Failed to write response", r.Method, r.URL.String(), err)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to write response"))
	}
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	msg := "Internal server error"

	logError(app.logger, msg, r.Method, r.URL.String(), err)

	app.errorResponse(w, r, http.StatusInternalServerError, msg)
}

func logError(logger *slog.Logger, msg, reqMethod, reqURL string, err error) {
	logger.Error(
		msg,
		slog.String("method", reqMethod),
		slog.String("path", reqURL),
		slog.String("error", err.Error()),
	)
}
