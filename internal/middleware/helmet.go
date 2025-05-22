package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
)

// Helmet applies common security headers (CSP, HSTS, XSS protection, etc.).
func Helmet() fiber.Handler {
	return helmet.New()
}
