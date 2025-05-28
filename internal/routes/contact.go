package routes

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/timkado/api/daisi-rest-postgres/internal/handler"
	"gitlab.com/timkado/api/daisi-rest-postgres/internal/middleware"
)

// ContactRoutes registers all /chats endpoints on the given router group.
func ContactRoutes(r fiber.Router) {
	contacts := r.Group("/contacts")

	// GET  /contacts?limit=...&offset=...&<filters>
	contacts.Get("/", middleware.Cache(), handler.FetchContacts)

	// GET  /contacts/:id
	contacts.Get("/:id", middleware.Cache(), handler.GetContactByID)

	// PATCH /contacts/:id
	contacts.Patch("/:id", handler.UpdateContact)
}
