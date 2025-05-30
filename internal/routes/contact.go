// internal/routes/contact.go
package routes

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/timkado/api/daisi-rest-postgres/internal/handler"
	"gitlab.com/timkado/api/daisi-rest-postgres/internal/middleware"
)

// ContactRoutes registers all /contacts endpoints on the given router group
func ContactRoutes(r fiber.Router) {
	contacts := r.Group("/contacts")

	// GET /contacts - Fetch paginated contacts with filters
	// Query params:
	// - limit (int): Number of items per page (default: 20, max: 100)
	// - offset (int): Number of items to skip (default: 0)
	// - sort (string): Sort field (created_at, updated_at, custom_name, phone_number, last_conversation_timestamp)
	// - order (string): Sort order (asc, desc)
	// - phone_number (string): Filter by exact phone number
	// - agent_id (string): Filter by agent ID
	// - assigned_to (string): Filter by assigned user
	// - tags (string): Filter by exact tag match (e.g., "TAG1" will match contacts with TAG1 but not TAG11)
	// - status (string): Filter by status (ACTIVE, DISABLED)
	// - origin (string): Filter by origin
	// - has_chat (bool): Filter contacts with/without associated chats
	// Response: { success: true, data: [...], total: X }
	contacts.Get("/", middleware.Cache(), handler.FetchContacts)

	// GET /contacts/search - Search contacts
	// Query params:
	// - q (string): Search query (required) - searches in phone_number, custom_name, push_name
	// - agent_id (string): Optional filter by agent ID
	// Response: { success: true, data: [...], total: X }
	contacts.Get("/search", middleware.Cache(), handler.SearchContacts)

	// GET /contacts/by-phone - Get contact by phone number and agent
	// Query params:
	// - phone_number (string): Phone number (required)
	// - agent_id (string): Agent ID (required)
	// Response: { success: true, data: {...} }
	contacts.Get("/by-phone", middleware.Cache(), handler.GetContactByPhoneAndAgent)

	// GET /contacts/:id - Get single contact by ID
	// Response: { success: true, data: {...} }
	contacts.Get("/:id", middleware.Cache(), handler.GetContactByID)

	// PATCH /contacts/:id - Update contact
	// Body: { custom_name?, assigned_to?, tags?, avatar?, notes? }
	// All fields are optional, only provided fields will be updated
	// Response: { success: true, data: {...} }
	contacts.Patch("/:id", handler.UpdateContact)
}
