package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/yourname/go-fiber-gorm-todo-auth-swagger/internal/models"
	"github.com/yourname/go-fiber-gorm-todo-auth-swagger/internal/service"
	"github.com/yourname/go-fiber-gorm-todo-auth-swagger/pkg/response"
)

type AuthHandler struct {
	svc service.AuthService
}

func NewAuthHandler(s service.AuthService) *AuthHandler {
	return &AuthHandler{svc: s}
}

// @Summary Register
// @Tags Auth
// @Accept json
// @Produce json
// @Param payload body map[string]interface{} true "Register body"
// @Success 201 {object} map[string]interface{}
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var body struct {
		Name     string       `json:"name"`
		Email    string       `json:"email"`
		Password string       `json:"password"`
		Role     models.Role  `json:"role"` // optional, default user
	}
	if err := c.BodyParser(&body); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}
	if body.Role == "" {
		body.Role = models.RoleUser
	}
	u, err := h.svc.Register(body.Name, body.Email, body.Password, body.Role)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}
	u.PasswordHash = ""
	return response.Created(c, u)
}

// @Summary Login
// @Tags Auth
// @Accept json
// @Produce json
// @Param payload body map[string]interface{} true "Login body"
// @Success 200 {object} map[string]interface{}
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}
	token, u, err := h.svc.Login(body.Email, body.Password)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, err.Error())
	}
	u.PasswordHash = ""
	return response.OK(c, fiber.Map{
		"token": token,
		"user":  u,
	})
}
