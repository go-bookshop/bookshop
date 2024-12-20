package repository

import (
	"bookshop/internal/validator"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

const (
	UserFirstNameMaxLength = 1000
	UserLastNameMaxLength  = 1000

	PasswordMinLength = 8
	PasswordMaxLength = 1000
)

type UserRepositoryInterface interface {
	Insert(u *User) error
	GetByEmail(email string) (*User, error)
	Update(u *User) error
	GetByToken(scope, token string) (*User, error)
}

func NewUserRepository(DBPool *pgxpool.Pool) UserRepositoryInterface {
	return &UserRepository{DBPool: DBPool}
}

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

type UserRepository struct {
	DBPool *pgxpool.Pool
}

func (r *UserRepository) Insert(u *User) error {
	query := `
insert into users(first_name, last_name, email, password, activated)
values ($1, $2, $3, $4, $5)
returning id, created_at, updated_at;`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{u.FirstName, u.LastName, u.Email, u.Password.hash, u.Activated}
	err := r.DBPool.QueryRow(ctx, query, args...).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return ErrDuplicateItem
		}
	}
	return err
}

func (r *UserRepository) GetByEmail(email string) (*User, error) {
	query := `
select u.id, u.first_name, u.last_name, u.email, u.password, u.activated, u.created_at, u.updated_at
from users u
where u.email = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var u User
	err := r.DBPool.QueryRow(ctx, query, email).Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password.hash,
		&u.Activated,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &u, nil
}

func (r *UserRepository) Update(u *User) error {
	query := `
update users
set first_name = $1, last_name = $2, email = $3, password = $4, activated = $5, updated_at = $6
where id = $7 and updated_at = $8
returning created_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{u.FirstName, u.LastName, u.Email, u.Password.hash, u.Activated, time.Now(), u.ID, u.UpdatedAt}
	err := r.DBPool.QueryRow(ctx, query, args...).Scan(&u.CreatedAt)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return ErrDuplicateItem
		}
	}
	if errors.Is(err, sql.ErrNoRows) {
		return ErrRecordEditConflict
	}
	return err
}

func (r *UserRepository) GetByToken(scope, token string) (*User, error) {
	query := `
select u.id, u.first_name, u.last_name, u.email, u.password, u.activated, u.created_at, u.updated_at
from users u
inner join tokens t on u.id = t.user_id
where t.hash = $1 and t.scope = $2 and t.expiry > $3`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tokenHash := sha256.Sum256([]byte(token))
	var user User
	err := r.DBPool.QueryRow(ctx, query, tokenHash[:], scope, time.Now()).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}
