package main

import (
	"bookshop/internal/assert"
	"bookshop/internal/data"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestUsers_registerUserHandler(t *testing.T) {
	app := newTestApplication()

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
			endpoint:    "/v1/users",
			wantCode:    http.StatusAccepted,
			wantBody:    `{"id":999,"firstName":"John","lastName":"Doe","email":"doe@mail.com","activated":false}`,
			requestBody: `{"firstName":"John","lastName":"Doe","email":"doe@mail.com","password":"PassWord123#"}`,
		},
		{
			name:        "Invalid JSON",
			endpoint:    "/v1/users",
			wantCode:    http.StatusBadRequest,
			wantBody:    "body contains badly-formed JSON",
			requestBody: "invalid",
		},
		{
			name:     "72+ byte password",
			endpoint: "/v1/users",
			wantCode: http.StatusAccepted,
			wantBody: `{"id":999,"firstName":"John","lastName":"Doe","email":"doe@mail.com","activated":false}`,
			requestBody: fmt.Sprintf(`{"firstName":"John","lastName":"Doe","email":"doe@mail.com","password":%q}`,
				strings.Repeat("Pa1#", 100)),
		},
		{
			name:        "Too short password",
			endpoint:    "/v1/users",
			wantCode:    http.StatusUnprocessableEntity,
			wantBody:    fmt.Sprintf("must be at least %d bytes long", data.PasswordMinLength),
			requestBody: `{"firstName":"John","lastName":"Doe","email":"doe@mail.com","password":"Pa1#"}`,
		},
		{
			name:     "Too long password",
			endpoint: "/v1/users",
			wantCode: http.StatusUnprocessableEntity,
			wantBody: fmt.Sprintf("must be less than %d bytes", data.PasswordMaxLength),
			requestBody: fmt.Sprintf(`{"firstName":"John","lastName":"Doe","email":"doe@mail.com","password":%q}`,
				strings.Repeat("Pa1#", data.PasswordMaxLength/4+5)),
		},
		{
			name:        "Empty password",
			endpoint:    "/v1/users",
			wantCode:    http.StatusUnprocessableEntity,
			wantBody:    "must be provided",
			requestBody: `{"firstName":"John","lastName":"Doe","email":"doe@mail.com","password":""}`,
		},
		{
			name:        "No uppercase password",
			endpoint:    "/v1/users",
			wantCode:    http.StatusUnprocessableEntity,
			wantBody:    "must contain at least one uppercase letter",
			requestBody: `{"firstName":"John","lastName":"Doe","email":"doe@mail.com","password":"password123#"}`,
		},
		{
			name:        "No lowercase password",
			endpoint:    "/v1/users",
			wantCode:    http.StatusUnprocessableEntity,
			wantBody:    "must contain at least one lowercase letter",
			requestBody: `{"firstName":"John","lastName":"Doe","email":"doe@mail.com","password":"PASSWORD123#"}`,
		},
		{
			name:        "No digit password",
			endpoint:    "/v1/users",
			wantCode:    http.StatusUnprocessableEntity,
			wantBody:    "must contain at least one digit",
			requestBody: `{"firstName":"John","lastName":"Doe","email":"doe@mail.com","password":"PASSWORDddd#"}`,
		},
		{
			name:        "No symbol password",
			endpoint:    "/v1/users",
			wantCode:    http.StatusUnprocessableEntity,
			wantBody:    "must contain at least one symbol",
			requestBody: `{"firstName":"John","lastName":"Doe","email":"doe@mail.com","password":"PASSWORDddd123"}`,
		},
		{
			name:        "Empty FirstName",
			endpoint:    "/v1/users",
			wantCode:    http.StatusUnprocessableEntity,
			wantBody:    "must be provided",
			requestBody: `{"firstName":"","lastName":"Doe","email":"doe@mail.com","password":"PASSWORDddd123#"}`,
		},
		{
			name:        "Whitespace FirstName",
			endpoint:    "/v1/users",
			wantCode:    http.StatusUnprocessableEntity,
			wantBody:    "must be provided",
			requestBody: `{"firstName":"     ","lastName":"Doe","email":"doe@mail.com","password":"PASSWORDddd123#"}`,
		},
		{
			name:     "Too long FirstName",
			endpoint: "/v1/users",
			wantCode: http.StatusUnprocessableEntity,
			wantBody: fmt.Sprintf("must be less than %d bytes", data.UserFirstNameMaxLength),
			requestBody: fmt.Sprintf(`{"firstName":%q,"lastName":"Doe","email":"doe@mail.com","password":"PASSWORDddd123#"}`,
				strings.Repeat("a", data.UserFirstNameMaxLength+1)),
		},
		{
			name:        "Empty LastName",
			endpoint:    "/v1/users",
			wantCode:    http.StatusUnprocessableEntity,
			wantBody:    "must be provided",
			requestBody: `{"firstName":"John","lastName":"","email":"doe@mail.com","password":"PASSWORDddd123#"}`,
		},
		{
			name:        "Whitespace LastName",
			endpoint:    "/v1/users",
			wantCode:    http.StatusUnprocessableEntity,
			wantBody:    "must be provided",
			requestBody: `{"firstName":"John","lastName":"    ","email":"doe@mail.com","password":"PASSWORDddd123#"}`,
		},
		{
			name:     "Too long LastName",
			endpoint: "/v1/users",
			wantCode: http.StatusUnprocessableEntity,
			wantBody: fmt.Sprintf("must be less than %d bytes", data.UserLastNameMaxLength),
			requestBody: fmt.Sprintf(`{"firstName":"John","lastName":%q,"email":"doe@mail.com","password":"PASSWORDddd123#"}`,
				strings.Repeat("a", data.UserLastNameMaxLength+1)),
		},
		{
			name:        "Empty Email",
			endpoint:    "/v1/users",
			wantCode:    http.StatusUnprocessableEntity,
			wantBody:    "must be provided",
			requestBody: `{"firstName":"John","lastName":"Doe","email":"","password":"PASSWORDddd123#"}`,
		},
		{
			name:        "Whitespace Email",
			endpoint:    "/v1/users",
			wantCode:    http.StatusUnprocessableEntity,
			wantBody:    "must be provided",
			requestBody: `{"firstName":"John","lastName":"Doe","email":"     ","password":"PASSWORDddd123#"}`,
		},
		{
			name:        "Invalid Email",
			endpoint:    "/v1/users",
			wantCode:    http.StatusUnprocessableEntity,
			wantBody:    "must be valid",
			requestBody: `{"firstName":"John","lastName":"Doe","email":"invalid","password":"PASSWORDddd123#"}`,
		},
		{
			name:        "Duplicate User",
			endpoint:    "/v1/users",
			wantCode:    http.StatusUnprocessableEntity,
			wantBody:    "user already exists",
			requestBody: `{"firstName":"John","lastName":"Doe","email":"duplicate@mail.com","password":"PassWord123#"}`,
		},
		{
			name:        "Unexpected error from UserRepository",
			endpoint:    "/v1/users",
			wantCode:    http.StatusInternalServerError,
			wantBody:    "Internal server error",
			requestBody: `{"firstName":"John","lastName":"Doe","email":"unexpected@mail.com","password":"PassWord123#"}`,
		},
		{
			name:        "Unexpected error from TokenRepository",
			endpoint:    "/v1/users",
			wantCode:    http.StatusInternalServerError,
			wantBody:    "Internal server error",
			requestBody: `{"firstName":"John","lastName":"Doe","email":"corrupted@mail.com","password":"PassWord123#"}`,
		},
		{
			name:        "Recovers from unexpected panic in email-sending goroutine",
			endpoint:    "/v1/users",
			wantCode:    http.StatusAccepted,
			wantBody:    `{"id":999,"firstName":"John","lastName":"Doe","email":"panic@mail.com","activated":false}`,
			requestBody: `{"firstName":"John","lastName":"Doe","email":"panic@mail.com","password":"PassWord123#"}`,
			wantLogs:    "recover me",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCode, gotHeader, gotBody := ts.post(t, tt.endpoint, tt.requestBody)
			assert.Equal(t, gotCode, tt.wantCode)
			assert.Equal(t, gotHeader.Get("Content-Type"), "application/json")
			assert.StringContains(t, gotBody, tt.wantBody)
			if tt.wantLogs != "" {
				logWriter := newMockLogWriter()
				setLoggerInterceptor(app, logWriter)
				for _, l := range logWriter.logs { //todo: it shouldn't check every log and stop on first found and not fail if not found in first string, maybe redesign logger?
					assert.StringContains(t, l, tt.wantLogs)
				}
			}
		})
	}

}
