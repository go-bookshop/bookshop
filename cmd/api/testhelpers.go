package main

import (
	"bookshop/internal/assert"
	"bookshop/internal/mock"
	"bookshop/internal/repository"
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync"
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
		mailer:       mock.NewMailTrap(),
	}
}

func newMockRepositories() repository.Repositories {
	return repository.Repositories{
		AuthorRepository:   mock.NewAuthorRepository(),
		CategoryRepository: mock.NewCategoryRepository(),
		UserRepository:     mock.NewUserRepository(),
		TokenRepository:    mock.NewTokenRepository(),
		BookRepository:     mock.NewBookRepository(),
	}
}

type mockLogWriter struct {
	logs []byte
	l    sync.RWMutex
}

func newMockLogWriter() *mockLogWriter {
	return &mockLogWriter{}
}

func (w *mockLogWriter) Write(p []byte) (n int, err error) {
	w.l.Lock()
	defer w.l.Unlock()
	w.logs = append(w.logs, p...)
	return len(p), nil
}

func (w *mockLogWriter) Logs() string {
	w.l.RLock()
	defer w.l.RUnlock()
	return string(w.logs)
}

func (w *mockLogWriter) cleanUp() {
	w.logs = make([]byte, 0)
}

func setLoggerInterceptor(app *application, w io.Writer) {
	app.logger = slog.New(slog.NewTextHandler(w, nil))
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
	assert.NoError(t, err)

	body = bytes.TrimSpace(body)
	return rs.StatusCode, rs.Header, string(body)
}

func (ts *testServer) post(t *testing.T, endpoint string, requestBody string) (int, http.Header, string) {
	rs, err := ts.Client().Post(ts.URL+endpoint, "application/json", bytes.NewBuffer([]byte(requestBody)))
	assert.NoError(t, err)

	defer rs.Body.Close()
	rsBody, err := io.ReadAll(rs.Body)
	assert.NoError(t, err)

	rsBody = bytes.TrimSpace(rsBody)
	return rs.StatusCode, rs.Header, string(rsBody)
}

func (ts *testServer) put(t *testing.T, endpoint string, requestBody string) (int, http.Header, string) {
	req, err := http.NewRequest(http.MethodPut, ts.URL+endpoint, bytes.NewBuffer([]byte(requestBody)))
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	rs, err := ts.Client().Do(req)
	assert.NoError(t, err)

	defer rs.Body.Close()
	rsBody, err := io.ReadAll(rs.Body)
	assert.NoError(t, err)

	rsBody = bytes.TrimSpace(rsBody)
	return rs.StatusCode, rs.Header, string(rsBody)
}
