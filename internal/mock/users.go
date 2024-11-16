package mock

import (
	"bookshop/internal/data"
	"errors"
	"time"
)

func NewUserRepository() data.UserRepositoryInterface {
	return &UserRepository{}
}

type UserRepository struct {
}

func (r *UserRepository) Insert(u *data.User) error {
	switch u.Email {
	case DuplicateEmail:
		return data.ErrDuplicateItem
	case UnexpectedEmail:
		return errors.New("unexpected")
	case CorruptedEmail:
		u.ID = UnexpectedUserId
	default:
		u.ID = 999
	}

	return nil
}

func (r *UserRepository) GetByEmail(email string) (*data.User, error) {
	switch email {
	case PanicEmail:
		return &data.User{Email: PanicEmail}, nil
	case FailedToSendEmail:
		return &data.User{Email: FailedToSendEmail}, nil
	case NotFoundEmail:
		return nil, data.ErrRecordNotFound
	case UnexpectedEmail:
		return nil, errors.New("empty")
	case ActivatedEmail:
		return &data.User{Activated: true}, nil
	case SimulateFailTokenCreationEmail:
		return &data.User{ID: UnexpectedUserId}, nil

	}
	return &data.User{
		ID:        999,
		FirstName: "John",
		LastName:  "Doe",
		Email:     ValidEmail,
		Activated: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (r *UserRepository) Update(u *data.User) error {
	switch u.Email {
	case UnexpectedEmail:
		return errors.New("empty")
	case ConflictEmail:
		return data.ErrRecordEditConflict
	}
	return nil
}

func (r *UserRepository) GetByToken(scope, token string) (*data.User, error) {
	switch token {
	case ExpiredToken:
		return nil, data.ErrRecordNotFound
	case SimulateConflictWriteToken:
		return &data.User{Email: ConflictEmail}, nil
	case SimulateUnexpectedWriteToken:
		return &data.User{Email: UnexpectedEmail}, nil
	case SimulateFailToDeleteAllByUserIdToken:
		return &data.User{ID: CorruptedUserId}, nil
	}
	return &data.User{
		ID:        999,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "doe@mail.com",
		Activated: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}
