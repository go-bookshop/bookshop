package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func (app *application) run() error {
	srv := &http.Server{
		Addr:         ":4000",
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app.logger.Info("starting server", "addr", srv.Addr)
	return srv.ListenAndServe()
}

func setupDbPool(cfg config) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	poolCfg, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	poolCfg.MaxConns = int32(cfg.db.maxConns)
	poolCfg.MinConns = int32(cfg.db.minConns)
	poolCfg.MaxConnIdleTime = cfg.db.maxIdleTime

	dbpool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, err
	}

	err = dbpool.Ping(ctx)
	if err != nil {
		return nil, errors.New(err.Error() + ": unable to connect to the database")
	}

	return dbpool, nil
}

func setupLogger(cfg config) *slog.Logger {
	var level slog.Level
	err := level.UnmarshalText([]byte(cfg.log.level))
	if err != nil {
		level = slog.LevelInfo
	}
	opts := &slog.HandlerOptions{Level: level}

	switch strings.ToLower(cfg.log.format) {
	case "text":
		return slog.New(slog.NewTextHandler(os.Stdout, opts))
	default:
		return slog.New(slog.NewJSONHandler(os.Stdout, opts))
	}
}
