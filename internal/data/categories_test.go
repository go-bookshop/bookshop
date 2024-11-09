package data

import (
	"bookshop/internal/assert"
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoriesRepoTestSuite struct {
	pgContainer *PostgresContainer
	repository  CategoriesRepositoryInterface
	ctx         context.Context
}

func (s *CategoriesRepoTestSuite) Setup(t *testing.T) {
	s.ctx = context.Background()

	pgContainer, err := CreatePostgresContainer(s.ctx)
	if err != nil {
		t.Fatal(err)
	}
	s.pgContainer = pgContainer

	pool, err := pgxpool.New(s.ctx, pgContainer.ConnectionString)
	assert.NoError(t, err)

	s.repository = NewCategoriesRepository(pool)
}

func (s *CategoriesRepoTestSuite) TearDown(t *testing.T) {
	if err := s.pgContainer.Terminate(s.ctx); err != nil {
		t.Fatalf("failed to terminate postgres container: %v", err)
	}
}

func (s *CategoriesRepoTestSuite) TestCategoriesRepository_Insert(t *testing.T) {
	category := Category{
		Name:        "Epic Adventures",
		Description: "Explore tales of heroism.",
	}
	t.Run("New Category", func(t *testing.T) {
		err := s.repository.Insert(&category)
		assert.NoError(t, err)
		assert.NonZero(t, category.ID, "id")
		assert.NonZero(t, category.CreatedAt, "createdAt")
		assert.NonZero(t, category.CreatedAt, "updatedAt")
	})

	t.Run("Duplicate Category", func(t *testing.T) {
		category.Name = "Duplicate"
		err := s.repository.Insert(&category)
		if !errors.Is(err, ErrDuplicateItem) {
			t.Errorf("should return %q error on Insert with duplicate name", ErrDuplicateItem)
		}
	})
}

func TestCategoriesRepoTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("repositories: skipping integration test")
	}
	s := &CategoriesRepoTestSuite{}
	s.Setup(t)
	defer s.TearDown(t)

	t.Run("TestCategoriesRepository_Insert", s.TestCategoriesRepository_Insert)
}
