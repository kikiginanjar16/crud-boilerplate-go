package handlers

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/yourname/go-fiber-gorm-todo-auth-swagger/internal/config"
	"github.com/yourname/go-fiber-gorm-todo-auth-swagger/internal/middleware"
	"github.com/yourname/go-fiber-gorm-todo-auth-swagger/internal/service"
	"github.com/yourname/go-fiber-gorm-todo-auth-swagger/pkg/response"
)

type ProfileHandler struct {
	cfg *config.Config
	us  service.UserService
}

func NewProfileHandler(cfg *config.Config, us service.UserService) *ProfileHandler {
	return &ProfileHandler{cfg: cfg, us: us}
}

// @Summary Get my profile
// @Security Bearer
// @Tags Profile
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /me [get]
func (h *ProfileHandler) Me(c *fiber.Ctx) error {
	uid, _ := middleware.GetUserID(c)
	u, err := h.us.GetByID(uid)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, err.Error())
	}
	u.PasswordHash = ""
	return response.OK(c, u)
}

// @Summary Update my profile
// @Security Bearer
// @Tags Profile
// @Accept json
// @Produce json
// @Param payload body map[string]interface{} true "Profile body"
// @Success 200 {object} map[string]interface{}
// @Router /me [put]
func (h *ProfileHandler) Update(c *fiber.Ctx) error {
	var body struct {
		Name string `json:"name"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}
	uid, _ := middleware.GetUserID(c)
	u, err := h.us.UpdateProfile(uid, body.Name)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}
	u.PasswordHash = ""
	return response.OK(c, u)
}

// @Summary Change password
// @Security Bearer
// @Tags Profile
// @Accept json
// @Produce json
// @Param payload body map[string]interface{} true "Password body"
// @Success 204 {string} string "No Content"
// @Router /me/password [patch]
func (h *ProfileHandler) ChangePassword(c *fiber.Ctx) error {
	var body struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}
	if len(body.NewPassword) < 6 {
		return response.Error(c, fiber.StatusBadRequest, "new password too short (min 6 chars)")
	}
	uid, _ := middleware.GetUserID(c)
	if err := h.us.ChangePassword(uid, body.OldPassword, body.NewPassword); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}
	return response.NoContent(c)
}

// @Summary Upload avatar
// @Security Bearer
// @Tags Profile
// @Accept multipart/form-data
// @Produce json
// @Param avatar formData file true "Avatar file"
// @Success 200 {object} map[string]interface{}
// @Router /me/avatar [post]
func (h *ProfileHandler) UploadAvatar(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("avatar")
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "avatar file required")
	}
	if fileHeader.Size > 2*1024*1024 {
		return response.Error(c, fiber.StatusBadRequest, "file too large (max 2MB)")
	}
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	switch ext {
	case ".png", ".jpg", ".jpeg", ".webp":
	default:
		return response.Error(c, fiber.StatusBadRequest, "invalid file type (png|jpg|jpeg|webp)")
	}

	src, err := fileHeader.Open()
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}
	defer src.Close()

	uid, _ := middleware.GetUserID(c)
	filename := fmt.Sprintf("u%d_%d%s", uid, time.Now().UnixNano(), ext)
	dstPath := filepath.Join(h.cfg.UploadDir, filename)

	if err := os.MkdirAll(h.cfg.UploadDir, 0o755); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	dst, err := os.Create(dstPath)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	// public URL
	publicURL := "/uploads/" + filename

	u, err := h.us.UpdateAvatarURL(uid, publicURL)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}
	u.PasswordHash = ""
	return response.OK(c, fiber.Map{"avatar_url": publicURL, "user": u})
}
