package main

import (
	"bookshop/internal/assert"
	"net/http"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	app := newTestApplication()

	ts := newTestServer(app.routes())
	defer ts.Close()

	code, header, body := ts.get(t, "/healthz")
	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, header.Get("Content-Type"), "application/json")
	assert.Equal(t, body, `{"status":"ok"}`)
}
