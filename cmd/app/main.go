package main

import (
	"app/internal/config"
	"app/internal/storage/sqlite"
	elog "app/pkg/lib/logger"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
		//todo delete
		newBanner := sqlite.Banner{
			FeatureID: 1,
			TagIDs:    []int{1, 2, 3},
			Content:   "",
			IsActive:  true,
			CreatedAt: "2022-12-31 23:59:59",
			UpdatedAt: "2022-12-31 23:59:59",
		}

		err = storage.AddBanner(newBanner)
		if err != nil {
			fmt.Println("Error adding banner:", err)
			return
		}
	}

	//роутер
	router := chi.NewRouter()
	//MW
	router.Use(middleware.RequestID) //ID запроса
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer) // поднять после паники в обраточике
	router.Use(middleware.URLFormat)

	//todo handlers

	// serv
	srv := http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// strat server
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server")
		}
	}()

	log.Info("server started")
	<-done

	// stop server
	log.Info("stopping server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", elog.Err(err))
		return
	}

	log.Info("server was shutfown")
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
