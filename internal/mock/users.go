package mock

import (
	"bookshop/internal/data"
	"errors"
	"strings"
)

func NewUserRepository() data.UserRepositoryInterface {
	return &UserRepository{}
}

type UserRepository struct {
}

func (r *UserRepository) Insert(u *data.User) error {
	if localPart, _, found := strings.Cut(u.Email, "@"); found {
		switch strings.ToLower(localPart) {
		case "duplicate":
			return data.ErrDuplicateItem
		case "unexpected":
			return errors.New("unexpected")
		case "corrupted":
			u.ID = -1
		default:
			u.ID = 999
		}
	}
	return nil
}

func (r *UserRepository) GetByEmail(email string) (*data.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r *UserRepository) Update(u *data.User) error {
	//TODO implement me
	panic("implement me")
}

func (r *UserRepository) GetByToken(scope, token string) (*data.User, error) {
	//TODO implement me
	panic("implement me")
}
