package mock

import (
	"bookshop/internal/models"
	"bookshop/internal/repository"
	"errors"
	"time"
)

func NewUserRepository() repository.UserRepositoryInterface {
	return &UserRepository{}
}

type UserRepository struct {
}

func (r *UserRepository) Insert(u *models.User) error {
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

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	switch email {
	case PanicEmail:
		return &models.User{Email: PanicEmail}, nil
	case FailedToSendEmail:
		return &models.User{Email: FailedToSendEmail}, nil
	case NotFoundEmail:
		return nil, repository.ErrRecordNotFound
	case UnexpectedEmail:
		return nil, errors.New("empty")
	case ActivatedEmail:
		return &models.User{Activated: true}, nil
	case SimulateFailTokenCreationEmail:
		return &models.User{ID: UnexpectedUserId}, nil

	}
	return &models.User{
		ID:        999,
		FirstName: "John",
		LastName:  "Doe",
		Email:     ValidEmail,
		Activated: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (r *UserRepository) Update(u *models.User) error {
	switch u.Email {
	case UnexpectedEmail:
		return errors.New("empty")
	case ConflictEmail:
		return repository.ErrRecordEditConflict
	}
	return nil
}

func (r *UserRepository) GetByToken(scope, token string) (*models.User, error) {
	switch token {
	case ExpiredToken:
		return nil, repository.ErrRecordNotFound
	case SimulateConflictWriteToken:
		return &models.User{Email: ConflictEmail}, nil
	case SimulateUnexpectedWriteToken:
		return &models.User{Email: UnexpectedEmail}, nil
	case SimulateFailToDeleteAllByUserIdToken:
		return &models.User{ID: CorruptedUserId}, nil
	}
	return &models.User{
		ID:        999,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "doe@mail.com",
		Activated: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}
