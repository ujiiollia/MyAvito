package main

import (
	"app/internal/config"
	"app/internal/storage/sqlite"
	elog "app/pkg/lib/logger"
	"fmt"
	"log/slog"
	"os"
)

func main() {
	// config
	cfg := config.MustLoad()
	// log
	log := setupLogger(cfg.Env)
	log.Info("starting application", slog.String("env", cfg.Env))
	// DB
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", elog.Err(err))
		os.Exit(1)
	}
	log.Info("storage created", slog.String("dbPath", cfg.StoragePath))

	//_ = storage
	// Создаем новую фичу
	newFeature := sqlite.Feature{
		ID: 1,
	}
	err = storage.AddFeature(newFeature)
	if err != nil {
		fmt.Println("Ошибка при создании фичи:", err)
		return
	}
	//todo logic
	//todo start serv
}

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func setupLogger(env string) *slog.Logger {
	var log slog.Logger
	switch env {
	case envLocal:
		log = *slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = *slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = *slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return &log
}
