package repository

import (
	"bookshop/internal/models"
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepositoryInterface interface {
	Insert(u *models.User) error
	GetByEmail(email string) (*models.User, error)
	Update(u *models.User) error
	GetByToken(scope, token string) (*models.User, error)
}

func NewUserRepository(DBPool *pgxpool.Pool) UserRepositoryInterface {
	return &UserRepository{DBPool: DBPool}
}

type UserRepository struct {
	DBPool *pgxpool.Pool
}

func (r *UserRepository) Insert(u *models.User) error {
	query := `
insert into users(first_name, last_name, email, password, activated)
values ($1, $2, $3, $4, $5)
returning id, created_at, updated_at;`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{u.FirstName, u.LastName, u.Email, u.Password.GetHash(), u.Activated}
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

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	query := `
select u.id, u.first_name, u.last_name, u.email, u.password, u.activated, u.created_at, u.updated_at
from users u
where u.email = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var u models.User
	var hash []byte
	err := r.DBPool.QueryRow(ctx, query, email).Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&hash,
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

	u.Password.SetHash(hash)

	return &u, nil
}

func (r *UserRepository) Update(u *models.User) error {
	query := `
update users
set first_name = $1, last_name = $2, email = $3, password = $4, activated = $5, updated_at = $6
where id = $7 and updated_at = $8
returning created_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{u.FirstName, u.LastName, u.Email, u.Password.GetHash(), u.Activated, time.Now(), u.ID, u.UpdatedAt}
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

func (r *UserRepository) GetByToken(scope, token string) (*models.User, error) {
	query := `
select u.id, u.first_name, u.last_name, u.email, u.password, u.activated, u.created_at, u.updated_at
from users u
inner join tokens t on u.id = t.user_id
where t.hash = $1 and t.scope = $2 and t.expiry > $3`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tokenHash := sha256.Sum256([]byte(token))
	var user models.User
	var hash []byte
	err := r.DBPool.QueryRow(ctx, query, tokenHash[:], scope, time.Now()).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&hash,
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
	user.Password.SetHash(hash)

	return &user, nil
}
