package middleware

import (
	"github.com/gofiber/fiber/v2"
	idempotency "github.com/gofiber/fiber/v2/middleware/idempotency"
)

// Idempotency returns Fiber’s built-in idempotency middleware.
// You can pass in an optional idempotency.Config to tweak store, header names, etc.
func Idempotency(cfg ...idempotency.Config) fiber.Handler {
	return idempotency.New(cfg...)
}
