package main

import (
	"bookshop/internal/assert"
	"bookshop/internal/data"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestCategories_createBooksCategoryHandler(t *testing.T) {
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
			wantLocation: "/v1/books/categories/999",
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
			name:        "Invalid Category Name Whitespace",
			endpoint:    "/v1/books/categories",
			wantCode:    http.StatusUnprocessableEntity,
			wantBody:    "must be provided",
			requestBody: `{"name": " ", "description":"Explore tales of heroism."}`,
		},
		{
			name:        "Invalid Category Description Whitespace",
			endpoint:    "/v1/books/categories",
			wantCode:    http.StatusUnprocessableEntity,
			wantBody:    "must be provided",
			requestBody: `{"name": "Epic Adventures", "description":" "}`,
		},
		{
			name:        "Invalid Category Name Max Length",
			endpoint:    "/v1/books/categories",
			wantCode:    http.StatusUnprocessableEntity,
			wantBody:    fmt.Sprintf("must be less than %d bytes", data.CategoryNameMaxLength),
			requestBody: fmt.Sprintf(`{"name":"%s","description":"Explore tales of heroism."}`, strings.Repeat("a", data.CategoryNameMaxLength+1)),
		},
		{
			name:        "Invalid Category Description Max Length",
			endpoint:    "/v1/books/categories",
			wantCode:    http.StatusUnprocessableEntity,
			wantBody:    fmt.Sprintf("must be less than %d bytes", data.CategoryDescriptionMaxLength),
			requestBody: fmt.Sprintf(`{"name":"Epic Adventures","description":"%s"}`, strings.Repeat("a", data.CategoryDescriptionMaxLength+1)),
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
			assert.StringContains(t, gotBody, tt.wantBody)
			if tt.wantLocation != "" {
				assert.Equal(t, gotHeader.Get("Location"), tt.wantLocation)
			}
		})
	}
}
