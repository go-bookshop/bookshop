package main

import (
	"bookshop/internal/httputil"
	"net/http"
)

type healthCheckResponse struct {
	Status string `json:"status"`
}

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	err := httputil.WriteJSON(w, http.StatusOK, healthCheckResponse{Status: "ok"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
