package data

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"testing"
)

type AuthorRepoTestSuite struct {
	pgContainer *PostgresContainer
	repository  *AuthorRepository
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
	assertNoError(t, err)

	s.repository = &AuthorRepository{DBPool: pool}
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
	err := s.repository.Insert(&author)
	assertNoError(t, err)
	if author.ID == 0 {
		t.Error("id has not been set from db after author creation")
	}
	if author.CreatedAt.IsZero() {
		t.Error("createdAt has not been set from db after author creation")
	}
	if author.UpdatedAt.IsZero() {
		t.Error("createdAt has not been set from db after author creation")
	}
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

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("received unexpected error, %v", err)
	}
}
