package main

import (
	"bookshop/internal/httputil"
	"net/http"
)

type healthCheckReponse struct {
	Status string `json:"status"`
}

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	err := httputil.WriteJSON(w, http.StatusOK, "", healthCheckReponse{Status: "ok"}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
