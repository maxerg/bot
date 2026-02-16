package auth

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct {
	pool *pgxpool.Pool
}

func NewPostgresStore(pool *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{pool: pool}
}

func (s *PostgresStore) Create(ctx context.Context, sess Session) error {
	const q = `
INSERT INTO sessions (id, user_id, created_at, expires_at)
VALUES ($1, $2, $3, $4);
`
	_, err := s.pool.Exec(ctx, q, sess.ID, sess.UserID, sess.CreatedAt, sess.ExpiresAt)
	return err
}

func (s *PostgresStore) Get(ctx context.Context, id string) (Session, bool) {
	const q = `
SELECT id, user_id, created_at, expires_at
FROM sessions
WHERE id = $1;
`
	var out Session
	err := s.pool.QueryRow(ctx, q, id).Scan(&out.ID, &out.UserID, &out.CreatedAt, &out.ExpiresAt)
	if err != nil {
		return Session{}, false
	}
	if time.Now().After(out.ExpiresAt) {
		return Session{}, false
	}
	return out, true
}

func (s *PostgresStore) Delete(ctx context.Context, id string) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM sessions WHERE id = $1`, id)
	return err
}
