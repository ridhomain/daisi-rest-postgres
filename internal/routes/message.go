// internal/routes/message.go
package routes

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/timkado/api/daisi-rest-postgres/internal/handler"
	"gitlab.com/timkado/api/daisi-rest-postgres/internal/middleware"
)

// MessageRoutes registers all /messages endpoints
func MessageRoutes(r fiber.Router) {
	messages := r.Group("/messages")

	// GET /messages - Fetch paginated messages for a specific chat
	// Query params:
	// - agent_id (string): Agent ID (required)
	// - chat_id (string): Chat ID (required)
	// - limit (int): Number of messages per page (default: 20, max: 100)
	// - offset (int): Number of messages to skip (default: 0)
	// Response: { success: true, data: [...], total: X }
	// Messages are sorted by message_timestamp DESC (newest first)
	messages.Get("/", middleware.Cache(), handler.FetchMessagesByChatId)

	// GET /messages/range - Fetch messages by range for infinite scroll
	// Query params:
	// - agent_id (string): Agent ID (required)
	// - chat_id (string): Chat ID (required)
	// - start (int): Start index (inclusive, default: 0)
	// - end (int): End index (inclusive, default: start)
	// Response: { success: true, data: [...] }
	// Maximum range size: 100 messages
	messages.Get("/range", middleware.Cache(), handler.FetchRangeMessagesByChatId)
}
