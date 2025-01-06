package repository

import (
	"bookshop/internal/models"
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthorRepositoryInterface interface {
	Insert(a *models.Author) error
}

func NewAuthorRepository(DBPool *pgxpool.Pool) AuthorRepositoryInterface {
	return &AuthorRepository{DBPool: DBPool}
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
