package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"bot/internal/auth"
	"bot/internal/domain/user"

	initdata "github.com/telegram-mini-apps/init-data-golang"
)

type AuthTelegramHandler struct {
	BotToken string
	Sessions auth.SessionStore
	Users    user.Repo
}

type authReq struct {
	InitData string `json:"initData"`
}
func realIP(r *http.Request) string {
	// Apache часто прокидывает это
	if v := r.Header.Get("X-Forwarded-For"); v != "" {
		// берем первый IP из списка
		for i := 0; i < len(v); i++ {
			if v[i] == ',' {
				return v[:i]
			}
		}
		return v
	}
	if v := r.Header.Get("X-Real-IP"); v != "" {
		return v
	}
	return r.RemoteAddr
}

func (h AuthTelegramHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if h.BotToken == "" || h.Users == nil || h.Sessions == nil {
		http.Error(w, "server misconfigured", http.StatusInternalServerError)
		return
	}

	var req authReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.InitData == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// 1) Валидация подписи + свежести
	const maxAge = 10 * time.Minute
	if err := initdata.Validate(req.InitData, h.BotToken, maxAge); err != nil {
		// НЕ логируем initData (там секреты), только причину + длины
		log.Printf("auth_telegram: validate failed: %v (initDataLen=%d tokenLen=%d)", err, len(req.InitData), len(h.BotToken))
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// 2) Парсим данные
	d, err := initdata.Parse(req.InitData)
	if err != nil {
		log.Printf("auth_telegram: parse failed: %v (initDataLen=%d)", err, len(req.InitData))
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Лог для обучения: насколько старые данные пришли
	ad := d.AuthDate().Unix()
	now := time.Now().Unix()
	log.Printf("auth_telegram: auth_date=%d now=%d diff=%ds", ad, now, now-ad)

	tgUser := d.User
	if tgUser.ID == 0 {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// 3) Upsert user
	u, err := h.Users.UpsertFromTelegram(r.Context(), tgUser.ID, tgUser.Username, tgUser.FirstName, tgUser.LastName)
	if err != nil {
		log.Printf("auth_telegram: users upsert failed: %v", err)
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	// 4) Создаём сессию
	sess, err := h.Sessions.Create(r.Context(), auth.Session{
		UserID:    u.ID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		UserAgent: r.UserAgent(),
		IP:        realIP(r),
	})
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "sid",
		Value:    sess.ID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400,
	})


	writeJSON(w, map[string]any{"ok": true})
}
