package repository

import (
	"bookshop/internal/assert"
	"bookshop/internal/models"
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepoTestSuite struct {
	pgContainer *PostgresContainer
	repository  UserRepositoryInterface
	ctx         context.Context
}

func (s *UserRepoTestSuite) Setup(t *testing.T) {
	s.ctx = context.Background()

	pgContainer, err := CreatePostgresContainer(s.ctx)
	if err != nil {
		t.Fatal(err)
	}
	s.pgContainer = pgContainer

	pool, err := pgxpool.New(s.ctx, pgContainer.ConnectionString)
	assert.NoError(t, err)

	files := []string{
		filepath.Join("../..", "testdata", "users.sql"),
		filepath.Join("../..", "testdata", "tokens.sql"),
	}
	err = RunTestData(pool, files...)
	assert.NoError(t, err)

	s.repository = NewUserRepository(pool)
}

func (s *UserRepoTestSuite) TearDown(t *testing.T) {
	if err := s.pgContainer.Terminate(s.ctx); err != nil {
		t.Fatalf("failed to terminate postgres container: %v", err)
	}
}

func (s *UserRepoTestSuite) TestUserRepository_Insert(t *testing.T) {
	table := []struct {
		name      string
		user      models.User
		password  string
		wantError error
	}{
		{
			name: "New User",
			user: models.User{
				FirstName: "Ivan",
				LastName:  "Petrenko",
				Email:     "Petr@mail.com",
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
			password: "Password123#",
		},
		{
			name: "Duplicate User",
			user: models.User{
				FirstName: "Ivan",
				LastName:  "Petrenko",
				Email:     "Petr@mail.com",
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
			password:  "Password123#",
			wantError: ErrDuplicateItem,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Password.Set(tt.password)
			assert.NoError(t, err)

			err = s.repository.Insert(&tt.user)
			if tt.wantError != nil {
				assert.Equal(t, err, tt.wantError)
			} else {
				assert.NoError(t, err)

				assert.NonZero(t, tt.user.ID, "id")
				assert.False(t, tt.user.Activated)
				assert.NonZero(t, tt.user.CreatedAt, "createdAt")
				assert.NonZero(t, tt.user.UpdatedAt, "updatedAt")

				matches, err := tt.user.Password.Matches(tt.password)
				assert.NoError(t, err)
				assert.True(t, matches)
			}
		})
	}
}

func (s *UserRepoTestSuite) TestUserRepository_GetByEmail(t *testing.T) {
	table := []struct {
		name          string
		wantEmail     string
		wantID        int64
		wantFirstName string
		wantLastName  string
		wantError     error
	}{
		{
			name:          "Existing User",
			wantEmail:     "diana@example.com",
			wantID:        4,
			wantFirstName: "Diana",
			wantLastName:  "Evans",
		},
		{
			name:      "Non-existing User",
			wantEmail: "fake@gmail.com",
			wantError: ErrRecordNotFound,
		},
	}
	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.repository.GetByEmail(tt.wantEmail)
			if tt.wantError != nil {
				assert.Equal(t, err, tt.wantError)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)

				assert.Equal(t, got.ID, tt.wantID)
				assert.Equal(t, got.FirstName, tt.wantFirstName)
				assert.Equal(t, got.LastName, tt.wantLastName)
				assert.Equal(t, got.Email, tt.wantEmail)
				assert.NonZero(t, got.CreatedAt, "createdAt")
				assert.NonZero(t, got.UpdatedAt, "updatedAt")
			}
		})
	}
}

func (s *UserRepoTestSuite) TestUserRepository_Update(t *testing.T) {
	table := []struct {
		name         string
		newFirstName string
		email        string
		newEmail     string
		wantError    error
	}{
		{
			name:         "Successful update",
			newFirstName: "Petro",
			email:        "alice@example.com",
		},
		{
			name:      "Existing email",
			email:     "alice@example.com",
			newEmail:  "bob@example.com",
			wantError: ErrDuplicateItem,
		},
		{
			name:         "Conflicting update",
			newFirstName: "Petro",
			email:        "charlie@example.com",
			wantError:    ErrRecordEditConflict,
		},
	}
	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			gotBeforeUpdate, _ := s.repository.GetByEmail(tt.email)
			updated := gotBeforeUpdate.UpdatedAt

			if tt.newFirstName != "" {
				gotBeforeUpdate.FirstName = tt.newFirstName
			}
			if tt.newEmail != "" {
				gotBeforeUpdate.Email = tt.newEmail
			}

			err := s.repository.Update(gotBeforeUpdate)
			if err != nil && tt.wantError != nil {
				assert.Equal(t, err, tt.wantError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, gotBeforeUpdate.FirstName, tt.newFirstName)

				gotAfterUpdate, _ := s.repository.GetByEmail(tt.email)
				assert.Equal(t, gotAfterUpdate.FirstName, tt.newFirstName)
				assert.NotEqual(t, gotAfterUpdate.UpdatedAt, updated)
			}

			if errors.Is(tt.wantError, ErrRecordEditConflict) {
				gotBeforeUpdate.UpdatedAt = updated
				err = s.repository.Update(gotBeforeUpdate)
				assert.Equal(t, err, ErrRecordEditConflict)
			}
		})
	}
}

func (s *UserRepoTestSuite) TestUserRepository_GetByToken(t *testing.T) {
	table := []struct {
		name       string
		scope      string
		token      string
		wantUserID int64
		wantError  error
	}{
		{
			name:       "Valid Token",
			token:      "35FM5TZVYCDCPQBCOWHVOVGE6Q",
			scope:      models.ScopeActivation,
			wantUserID: 4,
		},
		{
			name:      "Token with wrong scope",
			token:     "35FM5TZVYCDCPQBCOWHVOVGE6Q",
			scope:     models.ScopeAuthentication,
			wantError: ErrRecordNotFound,
		},
		{
			name:      "Non-existing Token",
			token:     "99FM5TZVYCDCPQBCDWHVOVGE6Q",
			scope:     models.ScopeActivation,
			wantError: ErrRecordNotFound,
		},
		{
			name:      "Expired Token",
			token:     "IS5BCEAMZJ7X5YD54Y6TW4VA2U",
			scope:     models.ScopeAuthentication,
			wantError: ErrRecordNotFound,
		},
	}
	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.repository.GetByToken(tt.scope, tt.token)
			if tt.wantError != nil {
				assert.Equal(t, err, tt.wantError)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, got.ID, tt.wantUserID)
			}
		})
	}
}

func TestUserRepoTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("repositories: skipping integration test")
	}
	s := &UserRepoTestSuite{}
	s.Setup(t)
	defer s.TearDown(t)

	t.Run("TestUserRepository_Insert", s.TestUserRepository_Insert)
	t.Run("TestUserRepository_GetByEmail", s.TestUserRepository_GetByEmail)
	t.Run("TestUserRepository_Update", s.TestUserRepository_Update)
	t.Run("TestUserRepository_GetByToken", s.TestUserRepository_GetByToken)
}
