package data

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Models struct {
}

func NewModels(dbpool *pgxpool.Pool) Models {
	return Models{}
}
