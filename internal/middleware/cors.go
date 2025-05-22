package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// CORS enables Cross-Origin Resource Sharing with default settings.
// Tweak Config if you need to restrict origins, methods, headers, etc.
func CORS() fiber.Handler {
	return cors.New()
}
