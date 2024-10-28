package main

import (
	"bookshop/internal/data"
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

type config struct {
	log struct {
		level  string
		format string
	}
	db struct {
		maxConns    int
		minConns    int
		maxIdleTime time.Duration
	}
}

type application struct {
	config config
	logger *slog.Logger
	models data.Models
}

func main() {
	var cfg config
	flag.StringVar(&cfg.log.level, "log-level", "info", "Logging level (debug|info|warning|error)")
	flag.StringVar(&cfg.log.format, "log-format", "json", "Logging format (text|json)")

	flag.IntVar(&cfg.db.maxConns, "dbpool-max-conns", 4, "Database max open connections")
	flag.IntVar(&cfg.db.minConns, "dbpool-min-conns", 1, "Database min idle connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "dbpool-max-idle-time", 15*time.Minute, "Database max connection idle time")

	flag.Parse()

	logger := setupLogger(cfg)

	dbpool, err := setupDbPool(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	logger.Info("successfully connected to the database")

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(dbpool),
	}

	err = app.dummyListenAndServe()
	if err != nil {
		os.Exit(1)
	}
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

func (app *application) dummyListenAndServe() error {
	mux := http.NewServeMux()

	ping := func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		fmt.Fprintf(w, "Thanks, %s", name)
	}

	mux.HandleFunc("GET /ping/{name}", ping)

	srv := &http.Server{
		Addr:         ":4000",
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app.logger.Info("starting server", "addr", srv.Addr)
	return srv.ListenAndServe()
}
