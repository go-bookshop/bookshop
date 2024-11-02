package data

import (
	"bookshop/internal/assert"
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"testing"
)

type AuthorRepoTestSuite struct {
	pgContainer *PostgresContainer
	repository  AuthorRepositoryInterface
	ctx         context.Context
}

func (s *AuthorRepoTestSuite) Setup(t *testing.T) {
	s.ctx = context.Background()

	pgContainer, err := CreatePostgresContainer(s.ctx)
	if err != nil {
		t.Fatal(err)
	}
	s.pgContainer = pgContainer

	pool, err := pgxpool.New(s.ctx, pgContainer.ConnectionString)
	assert.NoError(t, err)

	s.repository = NewAuthorRepository(pool)
}

func (s *AuthorRepoTestSuite) TearDown(t *testing.T) {
	if err := s.pgContainer.Terminate(s.ctx); err != nil {
		t.Fatalf("failed to terminate postgres container: %v", err)
	}
}

func (s *AuthorRepoTestSuite) TestAuthorRepository_Insert(t *testing.T) {
	author := Author{
		Name: "Evhen Petrenko",
		Bio:  "Born in Kyiv in 1913",
	}
	t.Run("New Author", func(t *testing.T) {
		err := s.repository.Insert(&author)
		assert.NoError(t, err)
		assert.NonZero(t, author.ID, "id")
		assert.NonZero(t, author.CreatedAt, "createdAt")
		assert.NonZero(t, author.CreatedAt, "updatedAt")
	})

	t.Run("Duplicate Author", func(t *testing.T) {
		err := s.repository.Insert(&author)
		if !errors.Is(err, ErrDuplicateAuthorName) {
			t.Errorf("should return %q error on Insert with duplicate name", ErrDuplicateAuthorName)
		}
	})
}

func TestAuthorRepoTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("repositories: skipping integration test")
	}
	s := &AuthorRepoTestSuite{}
	s.Setup(t)
	defer s.TearDown(t)

	t.Run("TestAuthorRepository_Insert", s.TestAuthorRepository_Insert)
}
