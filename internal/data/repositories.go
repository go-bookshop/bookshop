package data

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repositories struct {
	AuthorRepository     AuthorRepositoryInterface
	CategoriesRepository CategoriesRepositoryInterface
}

func NewRepositories(dbpool *pgxpool.Pool) Repositories {
	return Repositories{
		AuthorRepository:     NewAuthorRepository(dbpool),
		CategoriesRepository: NewCategoriesRepository(dbpool),
	}
}
