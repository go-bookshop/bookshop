package data

import (
	"bookshop/internal/validator"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	AuthorNameMaxLength = 1000
	AuthorBioMaxLength  = 2000
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
	v.Check(strings.TrimSpace(a.Name) != "", "name", "must be provided")
	v.Check(len(a.Name) <= AuthorNameMaxLength, "name", fmt.Sprintf("must be less than %d bytes long", AuthorNameMaxLength))

	v.Check(strings.TrimSpace(a.Bio) != "", "bio", "must be provided")
	v.Check(len(a.Bio) <= AuthorNameMaxLength, "bio", fmt.Sprintf("must be less than %d bytes long", AuthorBioMaxLength))
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
