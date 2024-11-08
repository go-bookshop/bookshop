package data

import (
	"bookshop/internal/validator"
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoriesRepositoryInterface interface {
	Insert(bc *Category) error
}

func NewCategoriesRepository(DBPool *pgxpool.Pool) CategoriesRepositoryInterface {
	return &CategoriesRepository{DBPool: DBPool}
}

type Category struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}

func ValidateCategory(v *validator.Validator, c *Category) {
	v.Check(strings.TrimSpace(c.Name) != "", "name", "must be provided")
	v.Check(len(c.Name) <= 100, "name", "must be less than 100 bytes long")

	v.Check(strings.TrimSpace(c.Description) != "", "description", "must be provided")
	v.Check(len(c.Description) <= 1000, "description", "must be less than 1000 bytes long")
}

type CategoriesRepository struct {
	DBPool *pgxpool.Pool
}

func (r *CategoriesRepository) Insert(c *Category) error {
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
