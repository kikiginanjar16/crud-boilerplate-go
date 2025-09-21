package response

import "github.com/gofiber/fiber/v2"

type Meta struct {
	Limit int `json:"limit,omitempty"`
	Page  int `json:"page,omitempty"`
	Total int64 `json:"total,omitempty"`
}

func OK(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": true, "data": data})
}

func List(c *fiber.Ctx, data interface{}, meta Meta) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": true,
		"data": data,
		"meta": meta,
	})
}

func Created(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": true, "data": data})
}

func NoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

func Error(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{"status": false, "error": message})
}
