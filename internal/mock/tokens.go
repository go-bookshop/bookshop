package mock

import (
	"bookshop/internal/data"
	"errors"
	"time"
)

func NewTokenRepository() data.TokenRepositoryInterface {
	return &TokenRepository{}
}

type TokenRepository struct {
}

func (t TokenRepository) New(userID int64, ttl time.Duration, scope string) (*data.Token, error) {
	switch userID {
	case UnexpectedUserId:
		return nil, errors.New("unexpected")
	}
	return &data.Token{UserID: userID}, nil
}

func (t TokenRepository) Insert(token *data.Token) error {
	return nil
}

func (t TokenRepository) DeleteAllByUserID(scope string, userID int64) error {
	switch userID {
	case CorruptedUserId:
		return errors.New("empty")
	}
	return nil
}
