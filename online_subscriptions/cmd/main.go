package main

import (
	"github.com/14kear/effective_mobile/online_subscriptions/internal/app"
	"github.com/14kear/effective_mobile/online_subscriptions/internal/config"
	"log"
	"os"
)

// @title Online Subscriptions API
// @description API для управления онлайн подписками
// @host localhost:8080
// @BasePath /api/
func main() {
	cfg := config.Load(os.Getenv("CONFIG_PATH"))

	application, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("Failed to init app: %v", err)
	}

	if err := application.Run(cfg.HTTP.Port); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
