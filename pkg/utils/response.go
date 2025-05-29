package utils

import "github.com/gofiber/fiber/v2"

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Total   int64       `json:"total,omitempty"`
}

func Success(c *fiber.Ctx, data interface{}) error {
	return c.JSON(APIResponse{Success: true, Data: data})
}

func SuccessWithTotal(c *fiber.Ctx, data interface{}, total int64) error {
	return c.JSON(APIResponse{Success: true, Data: data, Total: total})
}

func Error(c *fiber.Ctx, status int, errMsg string) error {
	return c.Status(status).JSON(APIResponse{Success: false, Error: errMsg})
}
