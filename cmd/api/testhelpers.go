package main

import (
	"bookshop/internal/assert"
	"bookshop/internal/data"
	"bookshop/internal/mock"
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestApplication() *application {
	cfg := config{}
	json := []byte(`{"openapi":"test"}`)
	cfg.doc.json = &json

	return &application{
		config:       cfg,
		logger:       slog.New(slog.NewTextHandler(io.Discard, nil)),
		repositories: newMockRepositories(),
	}
}

func newMockRepositories() data.Repositories {
	return data.Repositories{
		AuthorRepository: mock.NewAuthorRepository(),
	}
}

type testServer struct {
	*httptest.Server
}

func newTestServer(h http.Handler) *testServer {
	ts := httptest.NewServer(h)
	return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, endpoint string) (int, http.Header, string) {
	rs, err := ts.Client().Get(ts.URL + endpoint)
	if err != nil {
		assert.NoError(t, err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		assert.NoError(t, err)
	}
	body = bytes.TrimSpace(body)
	return rs.StatusCode, rs.Header, string(body)
}

func (ts *testServer) post(t *testing.T, endpoint string, requestBody string) (int, http.Header, string) {
	rs, err := ts.Client().Post(ts.URL+endpoint, "application/json", bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		assert.NoError(t, err)
	}

	defer rs.Body.Close()
	rsBody, err := io.ReadAll(rs.Body)
	if err != nil {
		assert.NoError(t, err)
	}
	rsBody = bytes.TrimSpace(rsBody)
	return rs.StatusCode, rs.Header, string(rsBody)
}
