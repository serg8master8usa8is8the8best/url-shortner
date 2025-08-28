package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"sergey/url-shortner/internal/config"
	"sergey/url-shortner/internal/http-server/handlers/redirect"
	"sergey/url-shortner/internal/http-server/handlers/url/save"
	mwLogger "sergey/url-shortner/internal/http-server/middleware/logger"
	"sergey/url-shortner/internal/lib/logger/handlers/slogpretty"
	"sergey/url-shortner/internal/lib/logger/sl"
	"sergey/url-shortner/internal/storage/sqlite"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// TODO

// refactor on postgress sql

// refactor on fiber

// add migration

// configure logger

// add delete and other handlers

// add unit tests

func main() {
	fmt.Println("hello")

	// config - cleanenv или viper
	cfg := config.MustLoad()

	// logger - slog

	log := setupLogger(cfg.Env)

	log.Info("start app", slog.String("env", cfg.Env))
	log.Debug("debug")
	log.Error("Erro message ")

	// storage - sqllite

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed load storage", sl.Err(err))
		os.Exit(1)
	}

	_ = storage

	// router - chi, chi render  или fiber

	router := chi.NewRouter()

	// middleware

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/url", save.New(log, storage))
	router.Get("/{alias}", redirect.New(log, storage))

	log.Info("starting server", slog.String("address", cfg.Addres))

	srv := &http.Server{
		Addr:         cfg.Addres,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")

	// server
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
