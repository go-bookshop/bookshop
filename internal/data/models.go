package data

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Models struct {
	AuthorModel *AuthorModel
}

func NewModels(dbpool *pgxpool.Pool) Models {
	return Models{
		AuthorModel: &AuthorModel{DBPool: dbpool},
	}
}
