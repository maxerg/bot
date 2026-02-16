package postgres

import (
	"context"

	"bot/internal/domain/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

func (r *UserRepo) UpsertFromTelegram(ctx context.Context, tgUserID int64, username, firstName, lastName string) (user.User, error) {
	const q = `
INSERT INTO users (tg_user_id, username, first_name, last_name)
VALUES ($1, $2, $3, $4)
ON CONFLICT (tg_user_id) DO UPDATE SET
  username   = EXCLUDED.username,
  first_name = EXCLUDED.first_name,
  last_name  = EXCLUDED.last_name,
  updated_at = now()
RETURNING id, tg_user_id, COALESCE(username,''), COALESCE(first_name,''), COALESCE(last_name,''), is_admin, created_at, updated_at;
`
	var u user.User
	err := r.pool.QueryRow(ctx, q, tgUserID, username, firstName, lastName).Scan(
		&u.ID, &u.TgUserID, &u.Username, &u.FirstName, &u.LastName, &u.IsAdmin, &u.CreatedAt, &u.UpdatedAt,
	)
	return u, err
}

func (r *UserRepo) GetByID(ctx context.Context, id int64) (user.User, bool, error) {
	const q = `
SELECT id, tg_user_id, username, first_name, last_name, is_admin, created_at, updated_at
FROM users
WHERE id = $1;
`
	var u user.User
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&u.ID, &u.TgUserID, &u.Username, &u.FirstName, &u.LastName, &u.IsAdmin, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return user.User{}, false, err
	}
	return u, true, nil
}
