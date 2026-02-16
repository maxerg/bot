package middleware

import (
	"context"
	"net/http"

	"bot/internal/auth"
	"bot/internal/domain/user"
)

// отдельный ключ, чтобы НЕ конфликтовать с request_id.go
type userCtxKeyType struct{}

var userCtxKey = userCtxKeyType{}

func UserFromContext(ctx context.Context) (user.User, bool) {
	u, ok := ctx.Value(userCtxKey).(user.User)
	return u, ok
}

func AuthRequired(sessions auth.SessionStore, users user.Repo) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := r.Cookie("sid")
			if err != nil || c.Value == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			sess, ok := sessions.Get(r.Context(), c.Value)
			if !ok {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			u, ok, err := users.GetByID(r.Context(), sess.UserID)
			if err != nil || !ok {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userCtxKey, u)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
