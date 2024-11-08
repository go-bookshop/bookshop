package data

import (
	"bookshop/internal/validator"
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthorRepositoryInterface interface {
	Insert(a *Author) error
}

func NewAuthorRepository(DBPool *pgxpool.Pool) AuthorRepositoryInterface {
	return &AuthorRepository{DBPool: DBPool}
}

type Author struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Bio       string    `json:"bio"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func ValidateAuthor(v *validator.Validator, a *Author) {
	v.Check(a.Name != "", "name", "must be provided")
	v.Check(len(a.Name) <= 1000, "name", "must be less than 1000 bytes long")

	v.Check(a.Bio != "", "bio", "must be provided")
	v.Check(len(a.Bio) <= 2000, "bio", "must be less than 2000 bytes long")
}

type AuthorRepository struct {
	DBPool *pgxpool.Pool
}

func (r *AuthorRepository) Insert(a *Author) error {
	query := `
insert into authors(name, bio)
values ($1, $2)
returning id, created_at, updated_at;`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return r.DBPool.QueryRow(ctx, query, a.Name, a.Bio).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)
}
