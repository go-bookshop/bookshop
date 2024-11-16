package main

import (
	"bookshop/internal/assert"
	"net/http"
	"testing"
)

func TestBooks_getBooks(t *testing.T) {
	app := newTestApplication()
	ts := newTestServer(app.routes())

	defer ts.Close()

	tests := []struct {
		name     string
		endpoint string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid request",
			endpoint: "/v1/books?page=0&size=1",
			wantCode: http.StatusOK,
			wantBody: `"page":0,"size":1,"max_pages":3,"data":[{`,
		},
		{
			name:     "Valid request (empty result)",
			endpoint: "/v1/books?page=10&size=1",
			wantCode: http.StatusOK,
			wantBody: `"page":10,"size":1,"max_pages":1,"data":[]`,
		},
		{
			name:     "Invalid query params (invalid page number value)",
			endpoint: "/v1/books?page=-1&size=1",
			wantCode: http.StatusBadRequest,
			wantBody: "invalid page number",
		},
		{
			name:     "Invalid query params (invalid page number format)",
			endpoint: "/v1/books?page=one&size=1",
			wantCode: http.StatusBadRequest,
			wantBody: "failed to parse page",
		},
		{
			name:     "Invalid query params (invalid page size value)",
			endpoint: "/v1/books?page=0&size=0",
			wantCode: http.StatusBadRequest,
			wantBody: "invalid page size",
		},
		{
			name:     "Invalid query params (invalid page size format)",
			endpoint: "/v1/books?page=0&size=zero",
			wantCode: http.StatusBadRequest,
			wantBody: "failed to parse size",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCode, gotHeader, gotBody := ts.get(t, tt.endpoint)
			assert.Equal(t, gotCode, tt.wantCode)
			assert.Equal(t, gotHeader.Get("Content-Type"), "application/json")
			assert.StringContains(t, gotBody, tt.wantBody)
		})
	}
}
