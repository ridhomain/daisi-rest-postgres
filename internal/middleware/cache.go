package middleware

import (
	"fmt"
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
		Expiration:   1 * time.Minute,
		CacheControl: true,
		// Optional manual bypass
		Next: func(c *fiber.Ctx) bool {
			return c.Query("refresh") == "true"
		},

		// Key includes companyId from token
		KeyGenerator: func(c *fiber.Ctx) string {
			companyId, ok := c.Locals("companyId").(string)
			if !ok || companyId == "" {
				// fallback to unauthenticated cache (or avoid caching)
				return "cache:unknown:" + c.OriginalURL()
			}
			return fmt.Sprintf("cache:%s:%s", companyId, c.OriginalURL())
		},
	})
}
