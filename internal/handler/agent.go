// internal/handler/agent.go
package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/timkado/api/daisi-rest-postgres/internal/model"
	"gitlab.com/timkado/api/daisi-rest-postgres/internal/service"
)

var agentSvc service.AgentService

// RegisterAgentService wires in the AgentService implementation.
func RegisterAgentService(svc service.AgentService) {
	agentSvc = svc
}

// ListAgents handles GET /agents?agentIds=... or GET /agents
func ListAgents(c *fiber.Ctx) error {
	companyId := c.Locals("companyId").(string)
	idsParam := c.Query("agentIds")

	var (
		agents []*model.Agent
		err    error
	)

	if idsParam != "" {
		agentIds := strings.Split(idsParam, ",")
		agents, err = agentSvc.ListByAgentIDs(c.Context(), companyId, agentIds)
	} else {
		agents, err = agentSvc.ListByCompanyID(c.Context(), companyId)
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}

	if agents == nil {
		agents = make([]*model.Agent, 0)
	}
	return c.JSON(agents)
}

// GetAgent handles GET /agents/:id
func GetAgent(c *fiber.Ctx) error {
	companyId := c.Locals("companyId").(string)
	agentId := c.Params("id")

	agent, err := agentSvc.GetByAgentID(c.Context(), companyId, agentId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}
	if agent == nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": "agent not found"})
	}
	return c.JSON(agent)
}

// CreateAgent handles POST /agents
func CreateAgent(c *fiber.Ctx) error {
	companyId := c.Locals("companyId").(string)
	var in model.Agent

	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}

	created, err := agentSvc.Create(c.Context(), companyId, &in)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(created)
}

// UpdateAgentName handles PATCH /agents/:id
func UpdateAgentName(c *fiber.Ctx) error {
	companyId := c.Locals("companyId").(string)
	agentId := c.Params("id")

	var body struct {
		AgentName string `json:"agent_name"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}

	updated, err := agentSvc.UpdateName(c.Context(), companyId, agentId, body.AgentName)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}
	if updated == nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": "agent not found"})
	}
	return c.JSON(updated)
}

// DeleteAgent handles DELETE /agents/:id
func DeleteAgent(c *fiber.Ctx) error {
	companyId := c.Locals("companyId").(string)
	agentId := c.Params("id")

	if err := agentSvc.Delete(c.Context(), companyId, agentId); err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
