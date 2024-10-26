package main

import (
	"flag"
	"fmt"
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
}

type application struct {
	config config
	logger *slog.Logger
}

func main() {
	var cfg config
	flag.StringVar(&cfg.log.level, "log-level", "info", "Logging level (debug|info|warning|error)")
	flag.StringVar(&cfg.log.format, "log-format", "json", "Logging format (text|json)")
	flag.Parse()

	app := &application{
		config: cfg,
		logger: setupLogger(cfg),
	}

	err := app.dummyListenAndServe()
	if err != nil {
		os.Exit(1)
	}
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
