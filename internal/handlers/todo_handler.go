package handlers

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/yourname/go-fiber-gorm-todo-auth-swagger/internal/middleware"
	"github.com/yourname/go-fiber-gorm-todo-auth-swagger/internal/models"
	"github.com/yourname/go-fiber-gorm-todo-auth-swagger/internal/service"
	"github.com/yourname/go-fiber-gorm-todo-auth-swagger/pkg/response"
)

type TodoHandler struct {
	svc service.TodoService
}

func NewTodoHandler(s service.TodoService) *TodoHandler {
	return &TodoHandler{svc: s}
}

// @Summary List todos
// @Security Bearer
// @Tags Todos
// @Produce json
// @Param limit query int false "limit"
// @Param page query int false "page"
// @Param q query string false "search"
// @Param completed query bool false "completed"
// @Param priority query string false "low|medium|high"
// @Param sort query string false "created_asc|created_desc|due_asc|due_desc"
// @Success 200 {object} map[string]interface{}
// @Router /todos [get]
func (h *TodoHandler) List(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	search := strings.TrimSpace(c.Query("q", ""))

	var completedPtr *bool
	if v := c.Query("completed", ""); v != "" {
		b := v == "true" || v == "1"
		completedPtr = &b
	}

	var priorityPtr *models.Priority
	if v := c.Query("priority", ""); v != "" {
		p := models.Priority(v)
		priorityPtr = &p
	}

	sort := c.Query("sort", "created_asc")

	ownerID, _ := middleware.GetUserID(c)

	items, total, err := h.svc.List(limit, page, search, completedPtr, priorityPtr, sort, &ownerID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}
	return response.List(c, items, response.Meta{Limit: limit, Page: page, Total: total})
}

// @Summary Get todo
// @Security Bearer
// @Tags Todos
// @Produce json
// @Param id path int true "Todo ID"
// @Success 200 {object} map[string]interface{}
// @Router /todos/{id} [get]
func (h *TodoHandler) Get(c *fiber.Ctx) error {
	id64, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid id")
	}
	obj, err := h.svc.Get(uint(id64))
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, err.Error())
	}
	return response.OK(c, obj)
}

// @Summary Create todo
// @Security Bearer
// @Tags Todos
// @Accept json
// @Produce json
// @Param payload body map[string]interface{} true "Todo body"
// @Success 201 {object} map[string]interface{}
// @Router /todos [post]
func (h *TodoHandler) Create(c *fiber.Ctx) error {
	var input models.Todo
	if err := c.BodyParser(&input); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}
	uid, _ := middleware.GetUserID(c)
	input.OwnerID = uid
	created, err := h.svc.Create(&input)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}
	return response.Created(c, created)
}

// @Summary Update todo
// @Security Bearer
// @Tags Todos
// @Accept json
// @Produce json
// @Param id path int true "Todo ID"
// @Param payload body map[string]interface{} true "Todo body"
// @Success 200 {object} map[string]interface{}
// @Router /todos/{id} [put]
func (h *TodoHandler) Update(c *fiber.Ctx) error {
	id64, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid id")
	}
	var input models.Todo
	if err := c.BodyParser(&input); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}
	updated, err := h.svc.Update(uint(id64), &input)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}
	return response.OK(c, updated)
}

// @Summary Toggle todo
// @Security Bearer
// @Tags Todos
// @Accept json
// @Produce json
// @Param id path int true "Todo ID"
// @Param payload body map[string]interface{} true "Toggle body"
// @Success 200 {object} map[string]interface{}
// @Router /todos/{id}/toggle [patch]
func (h *TodoHandler) Toggle(c *fiber.Ctx) error {
	id64, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid id")
	}
	var body struct {
		Completed bool `json:"completed"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}
	obj, err := h.svc.ToggleComplete(uint(id64), body.Completed)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}
	return response.OK(c, obj)
}

// @Summary Delete todo (admin only)
// @Security Bearer
// @Tags Todos
// @Produce json
// @Param id path int true "Todo ID"
// @Success 204 {string} string "No Content"
// @Router /todos/{id} [delete]
func (h *TodoHandler) Delete(c *fiber.Ctx) error {
	id64, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid id")
	}
	if err := h.svc.Delete(uint(id64)); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}
	return response.NoContent(c)
}
