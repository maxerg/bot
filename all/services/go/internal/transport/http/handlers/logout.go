package handlers

import (
	"net/http"

	"bot/internal/auth"
)

type LogoutHandler struct {
	Sessions auth.SessionStore
}

func (h LogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	c, err := r.Cookie("sid")
	if err == nil && c.Value != "" {
		_ = h.Sessions.Delete(r.Context(), c.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "sid",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	writeJSON(w, map[string]any{"ok": true})
}
