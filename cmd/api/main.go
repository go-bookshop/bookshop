package main

import (
	"bookshop/internal/data"
	"flag"
	"log/slog"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
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
	cfg := parseConfig()

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

	err = app.run()
	if err != nil {
		os.Exit(1)
	}
}

func parseConfig() config {
	var cfg config
	flag.StringVar(&cfg.log.level, "log-level", "info", "Logging level (debug|info|warning|error)")
	flag.StringVar(&cfg.log.format, "log-format", "json", "Logging format (text|json)")

	flag.IntVar(&cfg.db.maxConns, "dbpool-max-conns", 4, "Database max open connections")
	flag.IntVar(&cfg.db.minConns, "dbpool-min-conns", 1, "Database min idle connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "dbpool-max-idle-time", 15*time.Minute, "Database max connection idle time")

	flag.Parse()

	return cfg
}
