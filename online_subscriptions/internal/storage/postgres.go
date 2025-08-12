package storage

import (
	"fmt"
	"github.com/14kear/effective_mobile/online_subscriptions/internal/config"
	"github.com/14kear/effective_mobile/online_subscriptions/internal/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var db *gorm.DB

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DB.Host, cfg.DB.User, cfg.DB.Password, cfg.DB.Dbname, cfg.DB.Port, cfg.DB.Sslmode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&entity.Record{}); err != nil {
		log.Fatalf("Could not migrate table: %v", err)
	}

	return db, nil
}
