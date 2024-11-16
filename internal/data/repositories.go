package data

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repositories struct {
	AuthorRepository   AuthorRepositoryInterface
	UserRepository     UserRepositoryInterface
	TokenRepository    TokenRepositoryInterface
	CategoryRepository CategoryRepositoryInterface
	BookRepository     BookRepositoryInterface
}

func NewRepositories(dbpool *pgxpool.Pool) Repositories {
	return Repositories{
		AuthorRepository:   NewAuthorRepository(dbpool),
		UserRepository:     NewUserRepository(dbpool),
		TokenRepository:    NewTokenRepository(dbpool),
		CategoryRepository: NewCategoryRepository(dbpool),
		BookRepository:     NewBookRepository(dbpool),
	}
}
