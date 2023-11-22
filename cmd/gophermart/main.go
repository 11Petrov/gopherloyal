package main

import (
	"context"
	"net/http"

	"github.com/11Petrov/gopherloyal/internal/config"
	"github.com/11Petrov/gopherloyal/internal/handlers"
	"github.com/11Petrov/gopherloyal/internal/logger"
	"github.com/11Petrov/gopherloyal/internal/storage/postgre"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.NewConfig()
	log := logger.NewLogger()

	ctx := logger.ContextWithLogger(context.Background(), &log)

	if err := Run(cfg, ctx); err != nil {
		log.Fatal(err)
	}
}

func Run(cfg *config.Config, ctx context.Context) error {
	log := logger.LoggerFromContext(ctx)
	store, err := postgre.NewDBStore(cfg.DatabaseAddress, ctx)
	if err != nil {
		log.Errorf("store failed %s", err)
	}
	userHandler := handlers.NewUsersHandler(store)
	ordersHandler := handlers.NewOrdersHandler(store)

	r := chi.NewRouter()
	r.Use(logger.WithLogging)

	r.Post("/api/user/register", userHandler.UserRegister)
	r.Post("/api/user/login", userHandler.UserLogin)
	r.Post("/api/user/orders", ordersHandler.UploadOrder)
	r.Get("/api/user/orders", ordersHandler.GetUserOrders)

	log.Infof(
		"Running server",
		"address", cfg.ServerAddress,
		"DSN", cfg.DatabaseAddress,
	)
	return http.ListenAndServe(cfg.ServerAddress, r)
}
