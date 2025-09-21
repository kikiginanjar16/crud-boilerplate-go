package database

import (
	"fmt"

	"github.com/yourname/go-fiber-gorm-todo-auth-swagger/internal/config"
	"github.com/yourname/go-fiber-gorm-todo-auth-swagger/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBPort, cfg.DBSSLMode, cfg.DBTimezone,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}

	// Auto-migrate models
	if err := db.AutoMigrate(&models.User{}, &models.Todo{}); err != nil {
		return nil, err
	}

	return db, nil
}
