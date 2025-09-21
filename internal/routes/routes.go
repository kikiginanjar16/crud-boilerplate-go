package routes

import (
	"embed"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/gorm"

	"github.com/yourname/go-fiber-gorm-todo-auth-swagger/internal/config"
	"github.com/yourname/go-fiber-gorm-todo-auth-swagger/internal/handlers"
	"github.com/yourname/go-fiber-gorm-todo-auth-swagger/internal/middleware"
	"github.com/yourname/go-fiber-gorm-todo-auth-swagger/internal/repository"
	"github.com/yourname/go-fiber-gorm-todo-auth-swagger/internal/service"
)

//go:embed openapi/openapi.json
var openapiFS embed.FS

//go:embed openapi/redoc.html
var redocFS embed.FS

func NewFiberApp(cfg *config.Config, db *gorm.DB) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      "Go Fiber GORM TODO + JWT + OpenAPI",
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})

	app.Use(recover.New())
	app.Use(logger.New())

	// Serve uploaded files
	app.Static("/uploads", cfg.UploadDir)

	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"time":   time.Now().UTC(),
		})
	})

	// Serve OpenAPI spec and Redoc
	app.Get("/openapi.json", func(c *fiber.Ctx) error {
		b, err := openapiFS.ReadFile("openapi/openapi.json")
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		c.Set("Content-Type", "application/json")
		return c.Send(b)
	})
	app.Get("/docs", func(c *fiber.Ctx) error {
		b, err := redocFS.ReadFile("openapi/redoc.html")
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		c.Set("Content-Type", "text/html; charset=utf-8")
		return c.Send(b)
	})

	// DI
	userRepo := repository.NewUserRepository(db)
	authSvc := service.NewAuthService(cfg, userRepo)
	authHandler := handlers.NewAuthHandler(authSvc)

	userSvc := service.NewUserService(userRepo)
	profileHandler := handlers.NewProfileHandler(cfg, userSvc)

	todoRepo := repository.NewTodoRepository(db)
	todoSvc := service.NewTodoService(todoRepo)
	todoHandler := handlers.NewTodoHandler(todoSvc)

	api := app.Group("/api/v1")

	// Auth routes (public)
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	// Protected routes
	protected := api.Group("/", middleware.JWT(cfg))

	// Profile routes
	protected.Get("/me", profileHandler.Me)
	protected.Put("/me", profileHandler.Update)
	protected.Patch("/me/password", profileHandler.ChangePassword)
	protected.Post("/me/avatar", profileHandler.UploadAvatar)

	// Todos for authenticated users
	todos := protected.Group("/todos")
	todos.Get("/", todoHandler.List)
	todos.Get("/:id", todoHandler.Get)
	todos.Post("/", todoHandler.Create)
	todos.Put("/:id", todoHandler.Update)
	todos.Patch("/:id/toggle", todoHandler.Toggle)

	// Admin-only delete
	admin := todos.Use(middleware.RequireRoles("admin"))
	admin.Delete("/:id", todoHandler.Delete)

	return app
}
