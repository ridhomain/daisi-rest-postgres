package routes

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/timkado/api/daisi-rest-postgres/internal/handler"
	"gitlab.com/timkado/api/daisi-rest-postgres/internal/middleware"
)

// AgentRoutes registers all /agents endpoints on the given router group.
func AgentRoutes(r fiber.Router) {
	agents := r.Group("/agents")

	// GET /agents?agentids=... or /agents — cached
	agents.Get("/", middleware.Cache(), handler.ListAgents)

	// GET /agents/:agent_id — cached
	agents.Get("/:agent_id", middleware.Cache(), handler.GetAgent)

	// POST  /agents
	agents.Post("/", handler.CreateAgent)

	// PATCH /agents/:id
	agents.Patch("/:id", handler.UpdateAgentName)

	// DELETE /agents/:id
	agents.Delete("/:id", handler.DeleteAgent)
}
