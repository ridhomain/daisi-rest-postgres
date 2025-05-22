// internal/handler/chat.go
package handler

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/timkado/api/daisi-rest-postgres/internal/service"
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

	// parse pagination
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if page.Items == nil {
		page.Items = make([]map[string]interface{}, 0)
	}
	return c.JSON(fiber.Map{
		"total": page.Total,
		"items": page.Items,
	})
}

// FetchRangeChats handles GET /chats/range?&start=...&end=...
func FetchRangeChats(c *fiber.Ctx) error {
	companyId := c.Locals("companyId").(string)
	start := c.QueryInt("start", 0)
	end := c.QueryInt("end", start)

	items, err := chatSvc.FetchRangeChats(c.Context(), companyId, start, end)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if items == nil {
		items = make([]map[string]interface{}, 0)
	}
	return c.JSON(items)
}
