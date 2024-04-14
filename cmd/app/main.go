package main

import (
	"app/internal/config"
	mw "app/internal/middleware"
	"app/internal/middleware/auth"
	mapcache "app/internal/storage/cache"
	"app/internal/storage/postgresql"
	elog "app/pkg/lib/logger"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/golang-migrate/migrate/v4"
	"github.com/joho/godotenv"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}

	cfg := config.Config{}
	err = env.Parse(&cfg)
	if err != nil {
		panic("failed parse config")
	}
	fmt.Printf("Parsed Config: %+v", cfg)
	// log
	log := setupLogger(envLocal)

	log.Info("starting application")
	dbURL := cfg.PGLDsetination()
	fmt.Println("db url:", dbURL)

	mig, err := migrate.New("file://"+cfg.Migrations, dbURL)
	if err != nil {
		// log.Info("try force migration")
		// err = mig.Force(1) // sorry T_T
		// if err != nil {
		// 	log.Error("Could not force migrate: %v", err)
		// }
		log.Error("failed to migrate storage", elog.Err(err))
	}

	err = postgresql.ApplyMigrations(mig)
	if err != nil {
		log.Error("failed to apply migration", elog.Err(err))

	}
	log.Info("migration success")

	pool, err := postgresql.GetPgxPool(cfg.PGLDsetination(), cfg.MaxAttempts)
	if err != nil {
		log.Error("failed to get pool", elog.Err(err))
	}
	log.Info("connected to pool successfully")

	//кэш

	mcache := mapcache.New(5*time.Minute, 10*time.Minute)

	go func() {
		for {
			mcache.UpdateUserBannerCache(pool)
			time.Sleep(5 * time.Minute)
		}
	}()
	//роутер
	router := chi.NewRouter()
	//MW
	router.Use(middleware.RequestID) //ID запроса
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer) // поднять после паники в обраточике
	router.Use(middleware.URLFormat)

	// handlers
	router.Route("/user_banner", func(r chi.Router) {
		r.Use(auth.Authorization)
		r.Use(auth.CheckRole(auth.UserRightsRequired, pool))

		r.Get("/", mw.GetUserBanner(pool, mcache))
	})

	router.Route("/banner", func(r chi.Router) {
		r.Use(auth.Authorization)
		r.Use(auth.CheckRole(auth.AdminRightsRequired, pool))

		r.Get("/", mw.GetAllBannerByFeatureAndTag(pool))

		r.Post("/", mw.CreateBanner(pool))

		r.Patch("/{id}", mw.PatchBanner(pool))

		r.Delete("/{id}", mw.DeleteBanner(pool))
	})

	// serv
	srv := http.Server{
		Addr:         cfg.HTTPAddress,
		Handler:      router,
		ReadTimeout:  cfg.HTTPTimeout,
		WriteTimeout: cfg.HTTPTimeout,
		IdleTimeout:  cfg.HTTPIdleTimeout,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// start server

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

	log.Info("server was shutdown")

}

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
