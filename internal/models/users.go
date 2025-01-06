package models

import (
	"bookshop/internal/validator"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

const (
	UserFirstNameMaxLength = 1000
	UserLastNameMaxLength  = 1000

	PasswordMinLength = 8
	PasswordMaxLength = 1000
)

type User struct {
	ID        int64     `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Activated bool      `json:"activated"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type password struct {
	plaintext *string
	hash      []byte
}

func (p *password) Set(pwd string) error {
	hashed, err := bcrypt.GenerateFromPassword(encodePassword(pwd), 12)
	if err != nil {
		return err
	}
	p.plaintext = &pwd
	p.hash = hashed
	return nil
}

func (p *password) Matches(pwd string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, encodePassword(pwd))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func (p *password) GetHash() []byte {
	return p.hash
}

func (p *password) SetHash(hash []byte) {
	p.hash = hash
}

func encodePassword(pwd string) []byte {
	hash := sha256.Sum256([]byte(pwd))
	encoded := base64.StdEncoding.EncodeToString(hash[:])
	return []byte(encoded)
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(strings.TrimSpace(email) != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRegex), "email", "must be valid")
}

func ValidatePassword(v *validator.Validator, pwd string) {
	v.Check(strings.TrimSpace(pwd) != "", "password", "must be provided")
	v.Check(len(pwd) <= PasswordMaxLength, "password", fmt.Sprintf("must be less than %d bytes", PasswordMaxLength))
	v.Check(len(pwd) >= PasswordMinLength, "password", fmt.Sprintf("must be at least %d bytes long", PasswordMinLength))
	v.Check(strings.ContainsFunc(pwd, unicode.IsUpper), "password", "must contain at least one uppercase letter")
	v.Check(strings.ContainsFunc(pwd, unicode.IsLower), "password", "must contain at least one lowercase letter")
	v.Check(strings.ContainsFunc(pwd, unicode.IsDigit), "password", "must contain at least one digit")
	v.Check(strings.ContainsFunc(pwd, func(r rune) bool { return unicode.IsSymbol(r) || unicode.IsPunct(r) }), "password",
		"must contain at least one symbol")
}

func ValidateUser(v *validator.Validator, u *User) {
	v.Check(strings.TrimSpace(u.FirstName) != "", "first_name", "must be provided")
	v.Check(len(u.FirstName) <= UserFirstNameMaxLength, "first_name", fmt.Sprintf("must be less than %d bytes", UserFirstNameMaxLength))

	v.Check(strings.TrimSpace(u.LastName) != "", "last_name", "must be provided")
	v.Check(len(u.LastName) <= UserLastNameMaxLength, "last_name", fmt.Sprintf("must be less than %d bytes", UserLastNameMaxLength))

	ValidateEmail(v, u.Email)

	if u.Password.plaintext != nil {
		ValidatePassword(v, *u.Password.plaintext)
	}
}
