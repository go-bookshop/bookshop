package models

import (
	"bookshop/internal/validator"
	"fmt"
	"time"
)

const (
	ScopeActivation     = "activation"
	ScopeAuthentication = "authentication"
)

const (
	TokenLength = 26
)

type Token struct {
	Plaintext string
	Hash      []byte
	UserID    int64
	Expiry    time.Time
	Scope     string
}

func ValidateTokenPlaintext(v *validator.Validator, token string) {
	v.Check(token != "", "token", "must be provided")
	v.Check(len(token) == TokenLength, "token", fmt.Sprintf("must be %d bytes long", TokenLength))
}
