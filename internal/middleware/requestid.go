package middleware

import (
	"github.com/gofiber/fiber/v2"
	requestid "github.com/gofiber/fiber/v2/middleware/requestid"
)

// RequestID injects a unique X-Request-ID in each request/response header.
func RequestID() fiber.Handler {
	return requestid.New()
}
