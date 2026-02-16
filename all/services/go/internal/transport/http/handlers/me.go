package handlers

import (
	"net/http"

	"bot/internal/transport/http/middleware"
)

func Me(w http.ResponseWriter, r *http.Request) {
	u, ok := middleware.UserFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	writeJSON(w, map[string]any{
		"ok": true,
		"user": map[string]any{
			"id":        u.ID,
			"tgUserId":  u.TgUserID,
			"username":  u.Username,
			"firstName": u.FirstName,
			"lastName":  u.LastName,
			"isAdmin":   u.IsAdmin,
		},
	})
}
