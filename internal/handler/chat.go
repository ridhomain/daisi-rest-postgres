// internal/handler/chat.go
package handler

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/timkado/api/daisi-rest-postgres/internal/service"
	"gitlab.com/timkado/api/daisi-rest-postgres/pkg/utils"
)

var chatSvc service.ChatService

// RegisterChatService wires in the ChatService implementation.
func RegisterChatService(svc service.ChatService) {
	chatSvc = svc
}

// FetchChats handles GET /chats?limit=...&offset=...&<filters>
// Returns a JSON object with "total" and "items".
func FetchChats(c *fiber.Ctx) error {
	companyId := c.Locals("companyId").(string)
	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)

	// build filter map from all other query params
	filter := make(map[string]interface{})
	for k, v := range c.Queries() {
		switch k {
		case "limit", "offset":
			continue
		default:
			filter[k] = v
		}
	}

	page, err := chatSvc.FetchChats(c.Context(), companyId, filter, limit, offset)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	if page.Items == nil {
		page.Items = make([]map[string]interface{}, 0)
	}
	return utils.SuccessWithTotal(c, page.Items, page.Total)
}

// FetchRangeChats handles GET /chats/range?&start=...&end=...&<filters>
func FetchRangeChats(c *fiber.Ctx) error {
	companyId := c.Locals("companyId").(string)
	start := c.QueryInt("start", 0)
	end := c.QueryInt("end", start)

	// build filter map from all other query params
	filter := make(map[string]interface{})
	for k, v := range c.Queries() {
		switch k {
		case "limit", "offset":
			continue
		default:
			filter[k] = v
		}
	}

	items, err := chatSvc.FetchRangeChats(c.Context(), companyId, filter, start, end)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	if items == nil {
		items = make([]map[string]interface{}, 0)
	}
	return utils.Success(c, items)
}

// SearchChats handles GET /chats/search?q=query
func SearchChats(c *fiber.Ctx) error {
	companyId := c.Locals("companyId").(string)
	q := c.Query("q")

	if q == "" {
		return utils.SuccessWithTotal(c, []any{}, 0)
	}

	page, err := chatSvc.SearchChats(c.Context(), companyId, q)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	if page.Items == nil {
		page.Items = make([]map[string]interface{}, 0)
	}
	return utils.SuccessWithTotal(c, page.Items, page.Total)
}
