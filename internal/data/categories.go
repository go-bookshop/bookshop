package data

import (
	"bookshop/internal/validator"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	CategoryNameMaxLength        int = 100
	CategoryDescriptionMaxLength int = 1000
)

type CategoryRepositoryInterface interface {
	Insert(c *Category) error
}

func NewCategoryRepository(DBPool *pgxpool.Pool) CategoryRepositoryInterface {
	return &CategoryRepository{DBPool: DBPool}
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
	v.Check(len(c.Name) <= CategoryNameMaxLength, "name", fmt.Sprintf("must be less than %d bytes long", CategoryNameMaxLength))

	v.Check(strings.TrimSpace(c.Description) != "", "description", "must be provided")
	v.Check(len(c.Description) <= CategoryDescriptionMaxLength, "description", fmt.Sprintf("must be less than %d bytes long", CategoryDescriptionMaxLength))
}

type CategoryRepository struct {
	DBPool *pgxpool.Pool
}

func (r *CategoryRepository) Insert(c *Category) error {
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
