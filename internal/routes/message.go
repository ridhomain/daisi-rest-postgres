package routes

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/timkado/api/daisi-rest-postgres/internal/handler"
	"gitlab.com/timkado/api/daisi-rest-postgres/internal/middleware"
)

// MessageRoutes registers all /messages endpoints.
func MessageRoutes(r fiber.Router) {
	messages := r.Group("/messages")
	// GET  /messages?agentId=...&chatId=...&limit=...
	messages.Get("/", middleware.Cache(), handler.FetchMessagesByChatId)

	// GET  /messages/range?agentId=...&chatId=...&start=...&end=...
	messages.Get("/range", middleware.Cache(), handler.FetchRangeMessagesByChatId)
}
