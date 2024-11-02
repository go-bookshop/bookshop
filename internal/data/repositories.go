package data

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repositories struct {
	AuthorRepository *AuthorRepository
}

func NewRepositories(dbpool *pgxpool.Pool) Repositories {
	return Repositories{
		AuthorRepository: &AuthorRepository{DBPool: dbpool},
	}
}
