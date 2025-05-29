// internal/handler/agent.go
package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/timkado/api/daisi-rest-postgres/internal/model"
	"gitlab.com/timkado/api/daisi-rest-postgres/internal/service"
	"gitlab.com/timkado/api/daisi-rest-postgres/pkg/utils"
)

var agentSvc service.AgentService

// RegisterAgentService wires in the AgentService implementation.
func RegisterAgentService(svc service.AgentService) {
	agentSvc = svc
}

// ListAgents handles GET /agents?agentids=... or GET /agents
func ListAgents(c *fiber.Ctx) error {
	companyId := c.Locals("companyId").(string)
	idsParam := c.Query("agentids")

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
		return utils.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	if agents == nil {
		agents = make([]*model.Agent, 0)
	}
	return utils.Success(c, agents)
}

// GetAgent handles GET /agents/:agent_id
func GetAgent(c *fiber.Ctx) error {
	companyId := c.Locals("companyId").(string)
	agentId := c.Params("agent_id")

	agent, err := agentSvc.GetByAgentID(c.Context(), companyId, agentId)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	if agent == nil {
		return utils.Error(c, fiber.StatusNotFound, "agent not found")
	}
	return utils.Success(c, agent)
}

// CreateAgent handles POST /agents
func CreateAgent(c *fiber.Ctx) error {
	companyId := c.Locals("companyId").(string)
	var in model.Agent

	if err := c.BodyParser(&in); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, err.Error())
	}

	created, err := agentSvc.Create(c.Context(), companyId, &in)
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(utils.APIResponse{Success: true, Data: created})
}

// UpdateAgentName handles PATCH /agents/:id
func UpdateAgentName(c *fiber.Ctx) error {
	companyId := c.Locals("companyId").(string)
	agentId := c.Params("id")

	var body struct {
		AgentName string `json:"agent_name"`
	}
	if err := c.BodyParser(&body); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, err.Error())
	}

	updated, err := agentSvc.UpdateName(c.Context(), companyId, agentId, body.AgentName)
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, err.Error())
	}
	if updated == nil {
		return utils.Error(c, fiber.StatusNotFound, "agent not found")
	}
	return utils.Success(c, updated)
}

// DeleteAgent handles DELETE /agents/:id
func DeleteAgent(c *fiber.Ctx) error {
	companyId := c.Locals("companyId").(string)
	id := c.Params("id")

	if err := agentSvc.Delete(c.Context(), companyId, id); err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return c.Status(fiber.StatusNoContent).Send(nil)
}
