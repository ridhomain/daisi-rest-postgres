package handler

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/timkado/api/daisi-rest-postgres/internal/service"
	"gitlab.com/timkado/api/daisi-rest-postgres/pkg/utils"
)

var messageSvc service.MessageService

// RegisterMessageService wires in the MessageService implementation.
func RegisterMessageService(svc service.MessageService) {
	messageSvc = svc
}

// FetchMessagesByChatId handles:
// GET /messages?agentId=AGENTID&chat_id=CHATID&limit=N
// Returns { total: X, items: [...] }
func FetchMessagesByChatId(c *fiber.Ctx) error {
	companyId := c.Locals("companyId").(string)
	agentId := c.Query("agent_id")
	chatId := c.Query("chat_id")
	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)

	page, err := messageSvc.FetchMessagesByChatId(c.Context(), companyId, agentId, chatId, limit, offset)
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SuccessWithTotal(c, page.Items, page.Total)
}

// FetchRangeMessagesByChatId handles:
// GET /messages/range?agentId=AGENTID&chatId=CHATID&start=0&end=9
// Returns an array of messages (max end-start+1)
func FetchRangeMessagesByChatId(c *fiber.Ctx) error {
	companyId := c.Locals("companyId").(string)
	agentId := c.Query("agent_id")
	chatId := c.Query("chat_id")
	start := c.QueryInt("start", 0)
	end := c.QueryInt("end", start)

	items, err := messageSvc.FetchRangeMessagesByChatId(c.Context(), companyId, agentId, chatId, start, end)
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.Success(c, items)
}
