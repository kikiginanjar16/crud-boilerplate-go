package middleware

import (
	"errors"
	"strconv"

	jwtware "github.com/gofiber/jwt/v3"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"github.com/yourname/go-fiber-gorm-todo-auth-swagger/internal/config"
)

func JWT(cfg *config.Config) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   []byte(cfg.JWTSecret),
		ContextKey:   "jwt",
		ErrorHandler: jwtError,
		TokenLookup:  "header:Authorization,cookie:token",
		AuthScheme:   "Bearer",
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": false, "error": "unauthorized"})
}

func GetUserID(c *fiber.Ctx) (uint, error) {
	token := c.Locals("jwt")
	if token == nil {
		return 0, errors.New("no token")
	}
	tk := token.(*jwt.Token)
	claims, ok := tk.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}
	switch v := claims["sub"].(type) {
	case float64:
		return uint(v), nil
	case string:
		if id64, err := strconv.ParseUint(v, 10, 64); err == nil {
			return uint(id64), nil
		}
	}
	return 0, errors.New("invalid sub")
}

func GetUserRole(c *fiber.Ctx) string {
	token := c.Locals("jwt")
	if token == nil {
		return ""
	}
	tk := token.(*jwt.Token)
	claims, ok := tk.Claims.(jwt.MapClaims)
	if !ok {
		return ""
	}
	if r, ok := claims["role"].(string); ok {
		return r
	}
	return ""
}

func RequireRoles(roles ...string) fiber.Handler {
	allowed := map[string]struct{}{}
	for _, r := range roles {
		allowed[r] = struct{}{}
	}
	return func(c *fiber.Ctx) error {
		role := GetUserRole(c)
		if _, ok := allowed[role]; !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": false, "error": "forbidden"})
		}
		return c.Next()
	}
}
