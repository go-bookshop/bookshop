package data

import (
	"bookshop/internal/assert"
	"bookshop/internal/httputil"
	"context"
	"path/filepath"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BookRepoTestSuite struct {
	pgContainer *PostgresContainer
	repository  BookRepositoryInterface
	ctx         context.Context
}

func (s *BookRepoTestSuite) Setup(t *testing.T) {
	s.ctx = context.Background()

	pgContainer, err := CreatePostgresContainer(s.ctx)
	if err != nil {
		t.Fatal(err)
	}
	s.pgContainer = pgContainer

	pool, err := pgxpool.New(s.ctx, pgContainer.ConnectionString)
	assert.NoError(t, err)

	files := []string{
		filepath.Join("../..", "testdata", "books.sql"),
	}
	err = RunTestData(pool, files...)
	assert.NoError(t, err)

	s.repository = NewBookRepository(pool)
}

func (s *BookRepoTestSuite) TearDown(t *testing.T) {
	if err := s.pgContainer.Terminate(s.ctx); err != nil {
		t.Fatalf("failed to terminate postgres container: %v", err)
	}
}

func (s *BookRepoTestSuite) TestBookRepository_GetBooks(t *testing.T) {
	p := &httputil.PaginationData{
		PageNumber: 0,
		PageSize:   1,
	}

	books, count, err := s.repository.GetBooks(p)
	assert.NoError(t, err)
	assert.Equal(t, count, 6)
	assert.Equal(t, len(books), p.PageSize)

	for _, b := range books {
		assert.NonZero(t, b.ID, "id")
		assert.NonZero(t, b.Title, "title")
		assert.NonZero(t, b.CreatedAt, "created_at")
		assert.NonZero(t, b.UpdatedAt, "updated_at")
		assert.NotNil(t, b.ImageUrls)
		assert.True(t, b.AvgReview >= 0)
	}
}

func TestBookRepoTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("repositories: skipping integration test")
	}
	s := &BookRepoTestSuite{}
	s.Setup(t)
	defer s.TearDown(t)

	t.Run("TestBookRepository_GetBooks", s.TestBookRepository_GetBooks)
}
