package routes

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/timkado/api/daisi-rest-postgres/internal/handler"
	"gitlab.com/timkado/api/daisi-rest-postgres/internal/middleware"
)

// MessageRoutes registers all /messages endpoints.
func MessageRoutes(r fiber.Router) {
	messages := r.Group("/messages")
	// GET  /messages?agent_id=...&chat_id=...&limit=...
	messages.Get("/", middleware.Cache(), handler.FetchMessagesByChatId)

	// GET  /messages/range?agent_id=...&chat_id=...&start=...&end=...
	messages.Get("/range", middleware.Cache(), handler.FetchRangeMessagesByChatId)
}
