package user

import "context"

type Repo interface {
	UpsertFromTelegram(ctx context.Context, tgUserID int64, username, firstName, lastName string) (User, error)
	GetByID(ctx context.Context, id int64) (User, bool, error)
}
