package main

import (
	"car-mobile-project/config"
	uhttp "car-mobile-project/internal/user/delivery/http"
	"car-mobile-project/internal/user/repository"
	"car-mobile-project/internal/user/usecase"
	car_postgres "car-mobile-project/pkg/db/postgres"
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Println("Starting api server")
	cfgFileViper, err := config.LoadConfig("./config/config-local")

	if err != nil {
		log.Fatalf("LoadConfigError: %v", err)
	}

	cfg, err := config.ParseConfig(cfgFileViper)

	if err != nil {
		log.Fatalf("ParseConfigError: %v", err)
	}

	db, err := car_postgres.NewPsqlDB(cfg)

	if err != nil {
		log.Fatalf("PostgreSQLDBError: %v", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("PostgreSQLDBCloseError: %v", err)
		}
	}()

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})

	if err != nil {
		log.Fatalf("PostgreSQLInstanceError: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://./migrations", "postgres", driver)

	if err != nil {
		log.Fatalf("NewWithDatabaseInstanceError: %v", err)
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("migrate.Up: %v", err)
	} else {
		log.Println("migrate.Up success")
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	repo := repository.NewUserRepository(db)
	uc := usecase.NewUserUseCase(cfg, repo)
	uh := uhttp.NewUserHandler(cfg, uc)

	uhttp.MapUserRoutes(r, uh, cfg)

	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Printf("HttpServerError: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	// Attempt a graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting gracefully")

	if err = m.Down(); err != nil {
		log.Fatalf("Down err: %v", err)
	} else {
		log.Println("Down success")
	}
}
