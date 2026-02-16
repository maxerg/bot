package http

import (
	"net/http"

	"bot/internal/auth"
	"bot/internal/config"
	"bot/internal/domain/user"
	"bot/internal/transport/http/handlers"
	"bot/internal/transport/http/middleware"
)

type Deps struct {
	Cfg      config.Config
	Sessions auth.SessionStore
	Users    user.Repo
}

func NewRouter(d Deps) http.Handler {
	mux := http.NewServeMux()

	// public
	mux.HandleFunc("/api/ping", handlers.Ping)

	mux.Handle("/api/auth/telegram", handlers.AuthTelegramHandler{
		BotToken: d.Cfg.TelegramBotToken,
		Sessions: d.Sessions,
		Users:    d.Users,
	})

	// private
	needAuth := middleware.AuthRequired(d.Sessions, d.Users)

	mux.Handle("/api/me", needAuth(http.HandlerFunc(handlers.Me)))
	mux.Handle("/api/logout", needAuth(handlers.LogoutHandler{Sessions: d.Sessions}))

	return mux
}
