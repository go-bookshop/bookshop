package mock

import (
	"bookshop/internal/repository"
	"errors"
	"time"
)

func NewUserRepository() repository.UserRepositoryInterface {
	return &UserRepository{}
}

type UserRepository struct {
}

func (r *UserRepository) Insert(u *repository.User) error {
	switch u.Email {
	case DuplicateEmail:
		return repository.ErrDuplicateItem
	case UnexpectedEmail:
		return errors.New("unexpected")
	case CorruptedEmail:
		u.ID = UnexpectedUserId
	default:
		u.ID = 999
	}

	return nil
}

func (r *UserRepository) GetByEmail(email string) (*repository.User, error) {
	switch email {
	case PanicEmail:
		return &repository.User{Email: PanicEmail}, nil
	case FailedToSendEmail:
		return &repository.User{Email: FailedToSendEmail}, nil
	case NotFoundEmail:
		return nil, repository.ErrRecordNotFound
	case UnexpectedEmail:
		return nil, errors.New("empty")
	case ActivatedEmail:
		return &repository.User{Activated: true}, nil
	case SimulateFailTokenCreationEmail:
		return &repository.User{ID: UnexpectedUserId}, nil

	}
	return &repository.User{
		ID:        999,
		FirstName: "John",
		LastName:  "Doe",
		Email:     ValidEmail,
		Activated: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (r *UserRepository) Update(u *repository.User) error {
	switch u.Email {
	case UnexpectedEmail:
		return errors.New("empty")
	case ConflictEmail:
		return repository.ErrRecordEditConflict
	}
	return nil
}

func (r *UserRepository) GetByToken(scope, token string) (*repository.User, error) {
	switch token {
	case ExpiredToken:
		return nil, repository.ErrRecordNotFound
	case SimulateConflictWriteToken:
		return &repository.User{Email: ConflictEmail}, nil
	case SimulateUnexpectedWriteToken:
		return &repository.User{Email: UnexpectedEmail}, nil
	case SimulateFailToDeleteAllByUserIdToken:
		return &repository.User{ID: CorruptedUserId}, nil
	}
	return &repository.User{
		ID:        999,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "doe@mail.com",
		Activated: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}
