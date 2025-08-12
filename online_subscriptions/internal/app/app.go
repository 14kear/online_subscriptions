package app

import (
	"github.com/14kear/effective_mobile/online_subscriptions/internal/config"
	"github.com/14kear/effective_mobile/online_subscriptions/internal/handlers"
	"github.com/14kear/effective_mobile/online_subscriptions/internal/repository"
	"github.com/14kear/effective_mobile/online_subscriptions/internal/routes"
	"github.com/14kear/effective_mobile/online_subscriptions/internal/services"
	"github.com/14kear/effective_mobile/online_subscriptions/internal/storage"
	"github.com/gin-gonic/gin"
	"log"
	"log/slog"
)

func NewApp(cfg *config.Config) (*gin.Engine, error) {
	logger := slog.Default()

	database, err := storage.InitDB(cfg)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	repo := repository.NewRepository(database)

	service := services.NewRecordService(logger, repo)

	handler := handlers.NewRecordHandler(service)

	r := gin.Default()
	api := r.Group("/api")
	routes.RegisterRoutes(api, handler)

	return r, nil
}
