package main

import (
	"bookshop/internal/assert"
	"net/http"
	"testing"
)

func TestCategories(t *testing.T) {
	app := newTestApplication()
	ts := newTestServer(app.routes())
	defer ts.Close()

	tests := []struct {
		name         string
		endpoint     string
		wantCode     int
		wantBody     string
		wantLocation string
		requestBody  string
	}{
		{
			name:         "Valid request",
			endpoint:     "/v1/books/categories",
			wantCode:     http.StatusCreated,
			wantBody:     `{"id":999,"name":"Epic Adventures","description":"Explore tales of heroism."}`,
			wantLocation: "/v1/books/categories",
			requestBody:  `{"name":"Epic Adventures","description":"Explore tales of heroism."}`,
		},
		{
			name:        "Invalid JSON",
			endpoint:    "/v1/books/categories",
			wantCode:    http.StatusBadRequest,
			wantBody:    "body contains badly-formed JSON",
			requestBody: "invalid",
		},
		{
			name:        "Invalid Category",
			endpoint:    "/v1/books/categories",
			wantCode:    http.StatusUnprocessableEntity,
			wantBody:    "must be provided",
			requestBody: `{"description":"Explore tales of heroism."}`,
		},
		{
			name:        "Invalid Category Whitespace",
			endpoint:    "/v1/books/categories",
			wantCode:    http.StatusUnprocessableEntity,
			wantBody:    "must be provided",
			requestBody: `{"name": " ", "description":" "}`,
		},
		{
			name:        "Duplicate Category",
			endpoint:    "/v1/books/categories",
			wantCode:    http.StatusUnprocessableEntity,
			wantBody:    "already exists",
			requestBody: `{"name":"Duplicate","description":"Explore tales of heroism."}`,
		},
		{
			name:        "Unexpected error from DbPool",
			endpoint:    "/v1/books/categories",
			wantCode:    http.StatusInternalServerError,
			wantBody:    "Internal server error",
			requestBody: `{"name":"Unexpected","description":"Explore tales of heroism."}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCode, gotHeader, gotBody := ts.post(t, tt.endpoint, tt.requestBody)
			assert.Equal(t, gotCode, tt.wantCode)
			assert.Equal(t, gotHeader.Get("Content-Type"), "application/json")
			assert.Contains(t, gotBody, tt.wantBody)
			if tt.wantLocation != "" {
				assert.Equal(t, gotHeader.Get("Location"), tt.wantLocation)
			}
		})
	}
}
