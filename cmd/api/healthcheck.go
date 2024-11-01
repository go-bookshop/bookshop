package main

import (
	"bookshop/internal/utils"
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	err := utils.WriteJSON(w, http.StatusOK, "OK", nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
