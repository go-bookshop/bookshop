package utils

import (
	"encoding/json"
	"net/http"
)

func WriteJSON[T any](w http.ResponseWriter, status int, data T, headers http.Header) error {
	w.Header().Set("Content-Type", "application/json")

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}
