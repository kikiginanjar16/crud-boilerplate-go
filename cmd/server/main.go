package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/yourname/go-fiber-gorm-todo-auth-swagger/internal/config"
	"github.com/yourname/go-fiber-gorm-todo-auth-swagger/internal/database"
	"github.com/yourname/go-fiber-gorm-todo-auth-swagger/internal/routes"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	app := routes.NewFiberApp(cfg, db)

	addr := fmt.Sprintf(":%d", cfg.AppPort)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
