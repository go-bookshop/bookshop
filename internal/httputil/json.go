package httputil

import (
	"encoding/json"
	"net/http"
)

type response[T any] struct {
	Message string `json:"message,omitempty"`
	Data    T      `json:"data,omitempty"`
}

func WriteJSON[T any](w http.ResponseWriter, status int, msg string, data T, headers http.Header) error {
	resp := response[T]{
		Message: msg,
		Data:    data,
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	for key, value := range headers {
		w.Header()[key] = value
	}

	w.WriteHeader(status)

	_, writeErr := w.Write(jsonResp)
	if writeErr != nil {
		return writeErr
	}

	return nil
}
