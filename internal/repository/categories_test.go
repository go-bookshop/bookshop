package repository

import (
	"bookshop/internal/assert"
	"bookshop/internal/models"
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepoTestSuite struct {
	pgContainer *PostgresContainer
	repository  CategoryRepositoryInterface
	ctx         context.Context
}

func (s *CategoryRepoTestSuite) Setup(t *testing.T) {
	s.ctx = context.Background()

	pgContainer, err := CreatePostgresContainer(s.ctx)
	if err != nil {
		t.Fatal(err)
	}
	s.pgContainer = pgContainer

	pool, err := pgxpool.New(s.ctx, pgContainer.ConnectionString)
	assert.NoError(t, err)

	s.repository = NewCategoryRepository(pool)
}

func (s *CategoryRepoTestSuite) TearDown(t *testing.T) {
	if err := s.pgContainer.Terminate(s.ctx); err != nil {
		t.Fatalf("failed to terminate postgres container: %v", err)
	}
}

func (s *CategoryRepoTestSuite) TestCategoryRepository_Insert(t *testing.T) {
	category := models.Category{
		Name:        "Epic Adventures",
		Description: "Explore tales of heroism.",
	}
	t.Run("New Category", func(t *testing.T) {
		err := s.repository.Insert(&category)
		assert.NoError(t, err)
		assert.NonZero(t, category.ID, "id")
		assert.NonZero(t, category.CreatedAt, "createdAt")
		assert.NonZero(t, category.UpdatedAt, "updatedAt")
	})

	t.Run("Duplicate Category", func(t *testing.T) {
		err := s.repository.Insert(&category)
		if !errors.Is(err, ErrDuplicateItem) {
			t.Errorf("should return %q error on Insert with duplicate name", ErrDuplicateItem)
		}
	})
}

func TestCategoryRepoTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("repositories: skipping integration test")
	}
	s := &CategoryRepoTestSuite{}
	s.Setup(t)
	defer s.TearDown(t)

	t.Run("TestCategoriesRepository_Insert", s.TestCategoryRepository_Insert)
}
