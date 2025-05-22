package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	fibercache "github.com/gofiber/fiber/v2/middleware/cache"
)

// Cache returns a Fiber middleware that caches GET responses in-memory.
// - 5 minute expiration
// - honors Cache-Control headers
// - bypass with ?refresh=true
func Cache() fiber.Handler {
	return fibercache.New(fibercache.Config{
		Expiration:   5 * time.Minute,
		CacheControl: true,
		Next: func(c *fiber.Ctx) bool {
			return c.Query("refresh") == "true"
		},
	})
}
