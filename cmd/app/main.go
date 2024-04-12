package main

import (
	"app/internal/config"
	"app/internal/handlers"
	mw "app/internal/middleware"
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

	"github.com/caarlos0/env"
	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/golang-migrate/migrate"
)

func main() {
	// config
	cfg := config.MustLoad()
	cfgPGL := config.Config{}
	if err := env.Parse(&cfg); err != nil {
		panic("failed parse config")
	}
	// log
	log := setupLogger(cfg.Env)
	log.Info("starting application", slog.String("env", cfg.Env))

	mig, err := migrate.New("file://"+cfgPGL.MigrationPath, cfg.PgSQL.PGLDsetination())
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
	_ = hand
	//роутер
	router := chi.NewRouter()
	//MW
	router.Use(middleware.RequestID) //ID запроса
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer) // поднять после паники в обраточике
	router.Use(middleware.URLFormat)

	//todo handlers
	r := gin.Default()
	r.GET("/user_banner", mw.GetUserBanner(pool))
	r.GET("/banner", mw.GetBanner(pool))

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
