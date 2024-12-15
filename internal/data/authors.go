package data

import (
	"bookshop/internal/models"
	"bookshop/internal/validator"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthorRepositoryInterface interface {
	Insert(a *models.Author) error
}

func NewAuthorRepository(DBPool *pgxpool.Pool) AuthorRepositoryInterface {
	return &AuthorRepository{DBPool: DBPool}
}

func ValidateAuthor(v *validator.Validator, a *models.Author) {
	v.Check(strings.TrimSpace(a.Name) != "", "name", "must be provided")
	v.Check(len(a.Name) <= models.AuthorNameMaxLength, "name", fmt.Sprintf("must be less than %d bytes long", models.AuthorNameMaxLength))

	v.Check(strings.TrimSpace(a.Bio) != "", "bio", "must be provided")
	v.Check(len(a.Bio) <= models.AuthorNameMaxLength, "bio", fmt.Sprintf("must be less than %d bytes long", models.AuthorBioMaxLength))
}

type AuthorRepository struct {
	DBPool *pgxpool.Pool
}

func (r *AuthorRepository) Insert(a *models.Author) error {
	query := `
insert into authors(name, bio)
values ($1, $2)
returning id, created_at, updated_at;`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return r.DBPool.QueryRow(ctx, query, a.Name, a.Bio).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)
}
