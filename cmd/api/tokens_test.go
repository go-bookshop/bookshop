package main

import (
	"bookshop/internal/assert"
	"bookshop/internal/mock"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestTokens_resendActivationTokenHandler(t *testing.T) {
	app := newTestApplication()
	logWriter := newMockLogWriter()
	defer logWriter.cleanUp()
	setLoggerInterceptor(app, logWriter)

	ts := newTestServer(app.routes())
	defer ts.Close()

	tests := []struct {
		name        string
		endpoint    string
		wantCode    int
		wantBody    string
		wantLogs    string
		requestBody string
	}{
		{
			name:        "Valid request",
			endpoint:    "/v1/users/activation/resend-token",
			wantCode:    http.StatusAccepted,
			requestBody: fmt.Sprintf(`{"email":%q}`, mock.ValidEmail),
		},
		{
			name:        "Invalid JSON",
			endpoint:    "/v1/users/activation/resend-token",
			wantCode:    http.StatusBadRequest,
			wantBody:    "contains badly-formed JSON",
			requestBody: "invalid",
		},
		{
			name:        "Invalid email",
			endpoint:    "/v1/users/activation/resend-token",
			wantCode:    http.StatusUnprocessableEntity,
			wantBody:    "must be valid",
			requestBody: `{"email":"invalid"}`,
		},
		{
			name:        "Not found email",
			endpoint:    "/v1/users/activation/resend-token",
			wantCode:    http.StatusUnprocessableEntity,
			wantBody:    "no matching email address found",
			requestBody: fmt.Sprintf(`{"email":%q}`, mock.NotFoundEmail),
		},
		{
			name:        "Unexpected error",
			endpoint:    "/v1/users/activation/resend-token",
			wantCode:    http.StatusInternalServerError,
			wantBody:    "Internal server error",
			requestBody: fmt.Sprintf(`{"email":%q}`, mock.UnexpectedEmail),
		},
		{
			name:        "Already activated",
			endpoint:    "/v1/users/activation/resend-token",
			wantCode:    http.StatusUnprocessableEntity,
			wantBody:    "user has already been activated",
			requestBody: fmt.Sprintf(`{"email":%q}`, mock.ActivatedEmail),
		},
		{
			name:        "Token creation fail",
			endpoint:    "/v1/users/activation/resend-token",
			wantCode:    http.StatusInternalServerError,
			wantBody:    "Internal server error",
			requestBody: fmt.Sprintf(`{"email":%q}`, mock.SimulateFailTokenCreationEmail),
		},
		{
			name:        "Recovers from unexpected panic in email-sending goroutine",
			endpoint:    "/v1/users/activation/resend-token",
			wantCode:    http.StatusAccepted,
			requestBody: fmt.Sprintf(`{"email":%q}`, mock.PanicEmail),
			wantLogs:    "recover me",
		},
		{
			name:        "Logs if failed to send email",
			endpoint:    "/v1/users/activation/resend-token",
			wantCode:    http.StatusAccepted,
			requestBody: fmt.Sprintf(`{"email":%q}`, mock.FailedToSendEmail),
			wantLogs:    "failed to send email",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(logWriter.cleanUp)

			gotCode, gotHeader, gotBody := ts.post(t, tt.endpoint, tt.requestBody)
			assert.Equal(t, gotCode, tt.wantCode)
			if tt.wantBody != "" {
				assert.StringContains(t, gotBody, tt.wantBody)
				assert.Equal(t, gotHeader.Get("Content-Type"), "application/json")
			} else {
				assert.Zero(t, gotHeader.Get("Content-Type"), "Content-Type")
				assert.Zero(t, gotBody, "body")
			}
			if tt.wantLogs != "" {
				time.Sleep(10 * time.Millisecond) //waiting for email-sending goroutine
				assert.StringContains(t, logWriter.logs, tt.wantLogs)
			}
		})
	}
}
