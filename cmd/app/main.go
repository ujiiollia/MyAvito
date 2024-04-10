package main

import (
	"app/internal/config"
	"app/internal/handlers"
	"app/internal/services"
	"app/internal/storage/postgresql"

	elog "app/pkg/lib/logger"
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/golang-migrate/migrate"
)

func main() {
	// config
	cfg := config.MustLoad()
	// log
	log := setupLogger(cfg.Env)
	log.Info("starting application", slog.String("env", cfg.Env))
	// DB
	// storage, err := sqlite.NewBanner(cfg.StoragePath)
	// if err != nil {
	// 	log.Error("failed to init storage", elog.Err(err))
	// 	os.Exit(1)
	// }
	// log.Info("storage created", slog.String("dbPath", cfg.StoragePath))
	// _ = storage

	// user, err := sqlite.NewUser(cfg.StoragePath)
	// if err != nil {
	// 	log.Error("failed to init storage", elog.Err(err))
	// 	os.Exit(1)
	// }
	// _ = user

	//pgl
	mig, err := migrate.New("file://"+cfg.PgSQL.MigrationPath, cfg.PgSQL.PGLDsetination())
	if err != nil {
		log.Error("failed to migrate storage", elog.Err(err))
	}
	err = postgresql.ApplyMigrations(mig)
	if err != nil {
		log.Error("failed to apply migration", elog.Err(err))

	}
	log.Info("migration succsess")

	pool, err := postgresql.GetPgxPool(cfg.PGLDsetination(), cfg.MaxAttempts)
	if err != nil {
		log.Error("failed to get pool", elog.Err(err))
	}
	log.Info("connected to pool succsessfully")

	pg := postgresql.NewPostgres(pool)
	repo := services.NewBanner(pg)
	hand := handlers.NewBanner(repo)
	//роутер
	router := chi.NewRouter()
	//MW
	router.Use(middleware.RequestID) //ID запроса
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer) // поднять после паники в обраточике
	router.Use(middleware.URLFormat)

	//todo handlers
	_ = hand
	// router.Route("GET /ping", hand.Ping)
	//todo cache for banners (map [ID banner] srtuct, 5min update)
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
