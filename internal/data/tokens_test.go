package data

import (
	"bookshop/internal/assert"
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"path/filepath"
	"testing"
	"time"
)

type TokenRepoTestSuite struct {
	pgContainer *PostgresContainer
	repository  TokenRepositoryInterface
	ctx         context.Context
	pool        *pgxpool.Pool
}

func (s *TokenRepoTestSuite) Setup(t *testing.T) {
	s.ctx = context.Background()

	pgContainer, err := CreatePostgresContainer(s.ctx)
	if err != nil {
		t.Fatal(err)
	}
	s.pgContainer = pgContainer

	pool, err := pgxpool.New(s.ctx, pgContainer.ConnectionString)
	assert.NoError(t, err)
	s.pool = pool

	files := []string{
		filepath.Join("../..", "testdata", "users.sql"),
		filepath.Join("../..", "testdata", "tokens.sql"),
	}
	err = RunTestData(pool, files...)
	assert.NoError(t, err)

	s.repository = NewTokenRepository(pool)
}

func (s *TokenRepoTestSuite) TearDown(t *testing.T) {
	if err := s.pgContainer.Terminate(s.ctx); err != nil {
		t.Fatalf("failed to terminate postgres container: %v", err)
	}
}

func (s *TokenRepoTestSuite) TestTokenRepository_New(t *testing.T) {
	t.Run("New Token", func(t *testing.T) {
		got, err := s.repository.New(1, 24*time.Hour, ScopeActivation)
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.NonZero(t, got.Plaintext, "Plaintext")
		assert.NotNil(t, got.Hash)
		assert.NotNil(t, got)
		assert.Equal(t, got.UserID, 1)
		assert.Equal(t, got.Scope, ScopeActivation)
	})

	t.Run("Token with invalid User", func(t *testing.T) {
		_, err := s.repository.New(999, 24*time.Hour, ScopeActivation)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			assert.Equal(t, pgErr.Code, "23503")
		}
	})
}

func (s *TokenRepoTestSuite) TestTokenRepository_Insert(t *testing.T) {
	token := &Token{
		Plaintext: "A6IGKHLC25PMUOBTBJTEOGLX5U",
		Hash:      []byte{},
		UserID:    2,
		Expiry:    time.Now(),
		Scope:     ScopeAuthentication,
	}
	t.Run("New Token", func(t *testing.T) {
		err := s.repository.Insert(token)
		assert.NoError(t, err)
		assert.NotNil(t, token)
		assert.NonZero(t, token.Plaintext, "Plaintext")
		assert.NotNil(t, token.Hash)
		assert.NotNil(t, token)
		assert.Equal(t, token.UserID, 2)
		assert.Equal(t, token.Scope, ScopeAuthentication)
	})
}

func (s *TokenRepoTestSuite) TestTokenRepository_DeleteAllByUserID(t *testing.T) {
	countBefore := countTokensByUserId(t, s.pool, 4)
	err := s.repository.DeleteAllByUserID(ScopeActivation, 4)
	assert.NoError(t, err)
	countAfter := countTokensByUserId(t, s.pool, 4)
	assert.NotEqual(t, countBefore, countAfter)
	assert.Equal(t, countAfter, 0)
}

func countTokensByUserId(t *testing.T, pool *pgxpool.Pool, userID int64) int {
	t.Helper()
	query := "select count(*) from tokens t where t.user_id = $1"
	var countOfTokens int
	err := pool.QueryRow(context.Background(), query, userID).Scan(&countOfTokens)
	assert.NoError(t, err)
	return countOfTokens
}

func TestTokenRepoTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("repositories: skipping integration test")
	}
	s := &TokenRepoTestSuite{}
	s.Setup(t)
	defer s.TearDown(t)

	t.Run("TestTokenRepository_New", s.TestTokenRepository_New)
	t.Run("TestTokenRepository_Insert", s.TestTokenRepository_Insert)
	t.Run("TestTokenRepository_DeleteAllByUserID", s.TestTokenRepository_DeleteAllByUserID)
}
