package main

import (
	"bookshop/internal/assert"
	"bookshop/internal/models"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestAuthors_createAuthorHandler(t *testing.T) {
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
			endpoint:     "/v1/authors",
			wantCode:     http.StatusCreated,
			wantBody:     `{"id":999,"name":"Igor Olympic","bio":"Born today"}`,
			wantLocation: "/v1/authors/999",
			requestBody:  `{"name":"Igor Olympic","bio":"Born today"}`,
		},
		{
			name:        "Invalid JSON",
			endpoint:    "/v1/authors",
			wantCode:    http.StatusBadRequest,
			wantBody:    "body contains badly-formed JSON",
			requestBody: "invalid",
		},
		{
			name:        "Invalid Author",
			endpoint:    "/v1/authors",
			wantCode:    http.StatusUnprocessableEntity,
			wantBody:    "must be provided",
			requestBody: `{"bio":"Born today"}`,
		},
		{
			name:        "Invalid Author Name Whitespace",
			endpoint:    "/v1/authors",
			wantCode:    http.StatusUnprocessableEntity,
			wantBody:    "must be provided",
			requestBody: `{"name": " ", "bio":"Born today"}`,
		},
		{
			name:        "Invalid Author Bio Whitespace",
			endpoint:    "/v1/authors",
			wantCode:    http.StatusUnprocessableEntity,
			wantBody:    "must be provided",
			requestBody: `{"name": "Igor Olympic", "bio":" "}`,
		},
		{
			name:        "Invalid Author Name Max Length",
			endpoint:    "/v1/authors",
			wantCode:    http.StatusUnprocessableEntity,
			wantBody:    fmt.Sprintf("must be less than %d bytes", models.AuthorNameMaxLength),
			requestBody: fmt.Sprintf(`{"name":"%s","bio":"Born today"}`, strings.Repeat("a", models.AuthorNameMaxLength+1)),
		},
		{
			name:        "Invalid Author Bio Max Length",
			endpoint:    "/v1/authors",
			wantCode:    http.StatusUnprocessableEntity,
			wantBody:    fmt.Sprintf("must be less than %d bytes", models.AuthorBioMaxLength),
			requestBody: fmt.Sprintf(`{"name":"Igor Olympic","bio":"%s"}`, strings.Repeat("a", models.AuthorBioMaxLength+1)),
		},
		{
			name:        "Unexpected error from DbPool",
			endpoint:    "/v1/authors",
			wantCode:    http.StatusInternalServerError,
			wantBody:    "Internal server error",
			requestBody: `{"name":"Unexpected","bio":"Born today"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCode, gotHeader, gotBody := ts.post(t, tt.endpoint, tt.requestBody)
			assert.Equal(t, gotCode, tt.wantCode)
			assert.Equal(t, gotHeader.Get("Content-Type"), "application/json")
			assert.StringContains(t, gotBody, tt.wantBody)
			if tt.wantLocation != "" {
				assert.Equal(t, gotHeader.Get("Location"), tt.wantLocation)
			}
		})
	}
}
