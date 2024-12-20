package repository

import (
	"bookshop/internal/models"
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepositoryInterface interface {
	Insert(c *models.Category) error
}

func NewCategoryRepository(DBPool *pgxpool.Pool) CategoryRepositoryInterface {
	return &CategoryRepository{DBPool: DBPool}
}

type CategoryRepository struct {
	DBPool *pgxpool.Pool
}

func (r *CategoryRepository) Insert(c *models.Category) error {
	query := `
insert into categories(name, description)
values ($1, $2)
returning id, created_at, updated_at;`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.DBPool.QueryRow(ctx, query, c.Name, c.Description).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return ErrDuplicateItem
		}
	}
	return err
}
