package main

import (
	"app/internal/config"
	"app/internal/storage/sqlite"
	elog "app/pkg/lib/logger"
	"log/slog"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	// config
	cfg := config.MustLoad()
	// log
	log := setupLogger(cfg.Env)
	log.Info("starting application", slog.String("env", cfg.Env))
	// DB
	storage, err := sqlite.NewBanner(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", elog.Err(err))
		os.Exit(1)
	}
	log.Info("storage created", slog.String("dbPath", cfg.StoragePath))
	user, err := sqlite.NewUser(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", elog.Err(err))
		os.Exit(1)
	}

	_ = storage
	_ = user
	{ // Создаем Баннер
		// //todo delete
		// newBanner := sqlite.Banner{
		// 	FeatureID:       1,
		// 	TagIDs:          []int{1, 2, 3},
		// 	JSONData:        "",
		// 	Active:          true,
		// 	LastUpdatedTime: "2022-12-31 23:59:59",
		// }

		// err = storage.AddBanner(newBanner)
		// if err != nil {
		// 	fmt.Println("Error adding banner:", err)
		// 	return
		// }
	}

	//роутер
	router := chi.NewRouter()
	//MW
	router.Use(middleware.RequestID) //ID запроса
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer) // поднять после паники в обраточике
	router.Use(middleware.URLFormat)

	//todo handlers
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
