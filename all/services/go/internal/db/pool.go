package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func MustConnect(dsn string) *pgxpool.Pool {
	if dsn == "" {
		panic("DB_DSN is empty")
	}

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		panic(fmt.Errorf("parse DB_DSN: %w", err))
	}

	// разумные дефолты
	cfg.MaxConns = 10
	cfg.MinConns = 1
	cfg.MaxConnLifetime = 30 * time.Minute
	cfg.MaxConnIdleTime = 5 * time.Minute
	cfg.HealthCheckPeriod = 30 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		panic(fmt.Errorf("connect db: %w", err))
	}

	if err := pool.Ping(ctx); err != nil {
		panic(fmt.Errorf("ping db: %w", err))
	}

	return pool
}
