package data

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repositories struct {
	AuthorRepository   AuthorRepositoryInterface
	CategoryRepository CategoryRepositoryInterface
}

func NewRepositories(dbpool *pgxpool.Pool) Repositories {
	return Repositories{
		AuthorRepository:   NewAuthorRepository(dbpool),
		CategoryRepository: NewCategoryRepository(dbpool),
	}
}
