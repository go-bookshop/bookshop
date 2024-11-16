package main

import (
	"bookshop/internal/assert"
	"net/http"
	"testing"
)

func TestOpenAPISpecHandler(t *testing.T) {
	app := newTestApplication()
	ts := newTestServer(app.routes())
	defer ts.Close()

	gotCode, gotHeader, gotBody := ts.get(t, "/v1/doc/")
	assert.Equal(t, gotCode, http.StatusOK)
	assert.Equal(t, gotHeader.Get("Content-Type"), "application/json")
	assert.StringContains(t, gotBody, "openapi")
}

func TestRedocHandler(t *testing.T) {
	app := newTestApplication()
	ts := newTestServer(app.routes())
	defer ts.Close()

	gotCode, gotHeader, gotBody := ts.get(t, "/v1/doc/ui")
	assert.Equal(t, gotCode, http.StatusOK)
	assert.Equal(t, gotHeader.Get("Content-Type"), "text/html; charset=utf-8")
	assert.StringContains(t, gotBody, "Redoc")
}
