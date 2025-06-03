// internal/handler/message.go
package handler

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/timkado/api/daisi-rest-postgres/internal/service"
	"gitlab.com/timkado/api/daisi-rest-postgres/pkg/utils"
)

var messageSvc service.MessageService

// RegisterMessageService wires in the MessageService implementation
func RegisterMessageService(svc service.MessageService) {
	messageSvc = svc
}

// FetchMessagesByChatId handles GET /messages?agent_id=...&chat_id=...&limit=...&offset=...&sort=...&order=...
// Returns paginated messages for a specific chat
func FetchMessagesByChatId(c *fiber.Ctx) error {
	companyId := c.Locals("companyId").(string)
	agentId := c.Query("agent_id")
	chatId := c.Query("chat_id")
	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)
	sort := c.Query("sort", "message_timestamp")
	order := c.Query("order", "desc")

	// Validate required parameters
	if agentId == "" || chatId == "" {
		return utils.Error(c, fiber.StatusBadRequest, "agent_id and chat_id are required")
	}

	page, err := messageSvc.FetchMessagesByChatId(c.Context(), companyId, agentId, chatId, sort, order, limit, offset)
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SuccessWithTotal(c, page.Items, page.Total)
}

// FetchRangeMessagesByChatId handles GET /messages/range?agent_id=...&chat_id=...&start=...&end=...&sort=...&order=...
// Returns messages within a specific range for infinite scroll with total count
func FetchRangeMessagesByChatId(c *fiber.Ctx) error {
	companyId := c.Locals("companyId").(string)
	agentId := c.Query("agent_id")
	chatId := c.Query("chat_id")
	start := c.QueryInt("start", 0)
	end := c.QueryInt("end", start)
	sort := c.Query("sort", "message_timestamp")
	order := c.Query("order", "desc")

	// Validate required parameters
	if agentId == "" || chatId == "" {
		return utils.Error(c, fiber.StatusBadRequest, "agent_id and chat_id are required")
	}

	page, err := messageSvc.FetchRangeMessagesByChatId(c.Context(), companyId, agentId, chatId, sort, order, start, end)
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SuccessWithTotal(c, page.Items, page.Total)
}
