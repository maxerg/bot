package app

import (
	"log"
	"net/http"
	"time"

	"bot/internal/auth"
	"bot/internal/config"
	"bot/internal/db"
	"bot/internal/repository/postgres"
	transport "bot/internal/transport/http"
	"bot/internal/transport/http/middleware"
)

type App struct {
	Config config.Config
	Server *http.Server
}

func New(cfg config.Config) *App {
	// DB
	pool := db.MustConnect(cfg.DBDsn)

	// Repos
	userRepo := postgres.NewUserRepo(pool)

	// Sessions store in Postgres
	sessions := auth.NewPostgresStore(pool)

	h := transport.NewRouter(transport.Deps{
		Cfg:      cfg,
		Sessions: sessions,
		Users:    userRepo,
	})

	handler := middleware.Chain(h,
		middleware.Recover,
		middleware.RequestID,
		middleware.Logger,
		middleware.CORS(cfg.CORSAllowedOrigins),
	)

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	return &App{Config: cfg, Server: srv}
}

func (a *App) Run() error {
	log.Printf("api listening on %s env=%s", a.Config.HTTPAddr, a.Config.Env)
	return a.Server.ListenAndServe()
}
