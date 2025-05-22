package routes

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/timkado/api/daisi-rest-postgres/internal/handler"
	"gitlab.com/timkado/api/daisi-rest-postgres/internal/middleware"
)

// ChatRoutes registers all /chats endpoints on the given router group.
func ChatRoutes(r fiber.Router) {
	chats := r.Group("/chats")

	// GET  /chats?limit=...&offset=...&<filters>
	chats.Get("/", middleware.Cache(), handler.FetchChats)

	// GET  /chats/range?chatIdd=...&start=...&end=...
	chats.Get("/range", middleware.Cache(), handler.FetchRangeChats)
}
