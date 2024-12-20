package mock

import (
	"bookshop/internal/models"
	"bookshop/internal/repository"
	"errors"
	"time"
)

func NewTokenRepository() repository.TokenRepositoryInterface {
	return &TokenRepository{}
}

type TokenRepository struct {
}

func (t TokenRepository) New(userID int64, ttl time.Duration, scope string) (*models.Token, error) {
	switch userID {
	case UnexpectedUserId:
		return nil, errors.New("unexpected")
	}
	return &models.Token{UserID: userID}, nil
}

func (t TokenRepository) Insert(token *models.Token) error {
	return nil
}

func (t TokenRepository) DeleteAllByUserID(scope string, userID int64) error {
	switch userID {
	case CorruptedUserId:
		return errors.New("empty")
	}
	return nil
}
