// internal/routes/chat.go
package routes

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/timkado/api/daisi-rest-postgres/internal/handler"
	"gitlab.com/timkado/api/daisi-rest-postgres/internal/middleware"
)

// ChatRoutes registers all /chats endpoints on the given router group
func ChatRoutes(r fiber.Router) {
	chats := r.Group("/chats")

	// GET /chats - Fetch paginated chats with filters
	// Query params:
	// - limit (int): Number of items per page (default: 20, max: 100)
	// - offset (int): Number of items to skip (default: 0)
	// - agent_id (string): Filter by agent ID
	// - assigned_to (string): Filter by contact's assigned_to field
	// - has_unread (bool): Filter by unread status (true = unread_count > 0, false = unread_count = 0)
	// - is_group (bool): Filter by group chats
	// Response: { success: true, data: [...], total: X }
	chats.Get("/", middleware.Cache(), handler.FetchChats)

	// GET /chats/range - Fetch chats by range for infinite scroll
	// Query params:
	// - start (int): Start index (inclusive, default: 0)
	// - end (int): End index (inclusive, default: start)
	// - agent_id (string): Filter by agent ID
	// - assigned_to (string): Filter by contact's assigned_to field
	// - has_unread (bool): Filter by unread status
	// - is_group (bool): Filter by group chats
	// Response: { success: true, data: [...] }
	chats.Get("/range", middleware.Cache(), handler.FetchRangeChats)

	// GET /chats/search - Search chats and contacts
	// Query params:
	// - q (string): Search query (required) - searches in phone_number, push_name, group_name, custom_name
	// - agent_id (string): Optional filter by agent ID
	// Response: { success: true, data: [...], total: X }
	chats.Get("/search", middleware.Cache(), handler.SearchChats)
}
