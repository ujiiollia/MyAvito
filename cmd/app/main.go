package main

import (
	"app/internal/config"
	"fmt"
	"log/slog"
	"os"
)

func main() {
	// config
	cfg := config.MustLoad()
	fmt.Println(cfg.Address)
	// log
	log := setupLogger(cfg.Env)
	log.Info("starting application ", slog.String("env", cfg.Env))
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
