package middleware

import (
	"github.com/gofiber/fiber/v2"
	fiberRecover "github.com/gofiber/fiber/v2/middleware/recover"
)

// Recover catches panics in handlers and returns a 500 instead of crashing.
func Recover() fiber.Handler {
	return fiberRecover.New()
}
