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
func (s *PostgresStore) Create(ctx context.Context, sess Session) (Session, error) {
    const q = `
INSERT INTO sessions (id, user_id, created_at, expires_at, user_agent, ip)
VALUES (gen_random_uuid(), $1, now(), $2, $3, $4)
RETURNING id, user_id, created_at, expires_at;
`
    var out Session
    err := s.pool.QueryRow(ctx, q, sess.UserID, sess.ExpiresAt, sess.UserAgent, sess.IP).
        Scan(&out.ID, &out.UserID, &out.CreatedAt, &out.ExpiresAt)
    if err != nil {
        return Session{}, err
    }
    return out, nil
}


func (s *PostgresStore) Get(ctx context.Context, id string) (Session, bool) {
	const q = `
SELECT id, user_id, created_at, expires_at, COALESCE(user_agent,''), COALESCE(ip::text,'')
FROM sessions
WHERE id = $1;
`
	var out Session
	err := s.pool.QueryRow(ctx, q, id).
		Scan(&out.ID, &out.UserID, &out.CreatedAt, &out.ExpiresAt, &out.UserAgent, &out.IP)
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
