package repository

import (
	"bookshop/internal/models"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TokenRepositoryInterface interface {
	New(userID int64, ttl time.Duration, scope string) (*models.Token, error)
	Insert(token *models.Token) error
	DeleteAllByUserID(scope string, userID int64) error
}

func NewTokenRepository(DBPool *pgxpool.Pool) TokenRepositoryInterface {
	return &TokenRepository{DBPool: DBPool}
}

func generateToken(userID int64, ttl time.Duration, scope string) (*models.Token, error) {
	token := &models.Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token, nil
}

type TokenRepository struct {
	DBPool *pgxpool.Pool
}

func (r *TokenRepository) New(userID int64, ttl time.Duration, scope string) (*models.Token, error) {
	token, err := generateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = r.Insert(token)
	return token, err
}

func (r *TokenRepository) Insert(token *models.Token) error {
	query := `
insert into tokens (hash, user_id, expiry, scope)
values ($1, $2, $3, $4)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{token.Hash, token.UserID, token.Expiry, token.Scope}
	_, err := r.DBPool.Exec(ctx, query, args...)
	return err
}

func (r *TokenRepository) DeleteAllByUserID(scope string, userID int64) error {
	query := `
delete from tokens
where scope = $1 and user_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := r.DBPool.Exec(ctx, query, scope, userID)
	return err
}
