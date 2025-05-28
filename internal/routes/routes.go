package routes

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/timkado/api/daisi-rest-postgres/internal/middleware"
)

// RegisterRoutes mounts all sub-route groups under /api/v1
func RegisterV1Routes(app *fiber.App) {
	// Create /api/v1 group
	v1 := app.Group("/api/v1", middleware.AuthenticateBearerToken())

	// Mount each resource under /api/v1
	AgentRoutes(v1)
	ChatRoutes(v1)
	MessageRoutes(v1)
	ContactRoutes(v1)
}
