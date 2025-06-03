// internal/handler/chat.go
package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/timkado/api/daisi-rest-postgres/internal/service"
	"gitlab.com/timkado/api/daisi-rest-postgres/pkg/utils"
)

var chatSvc service.ChatService

// RegisterChatService wires in the ChatService implementation
func RegisterChatService(svc service.ChatService) {
	chatSvc = svc
}

// FetchChats handles GET /chats?limit=...&offset=...&<filters>
// Common filters: agent_id, assigned_to, has_unread
// Returns a JSON object with "total" and "items"
func FetchChats(c *fiber.Ctx) error {
	companyId := c.Locals("companyId").(string)
	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)

	// Build filter map
	filter := make(map[string]interface{})

	// Most common filter: agent_id
	if agentId := c.Query("agent_id"); agentId != "" {
		filter["agent_id"] = agentId
	}

	// Second most common: assigned_to (from contacts)
	if assignedTo := c.Query("assigned_to"); assignedTo != "" {
		filter["assigned_to"] = assignedTo
	}

	// Unread filter - convert to boolean
	if hasUnreadStr := c.Query("has_unread"); hasUnreadStr != "" {
		if hasUnread, err := strconv.ParseBool(hasUnreadStr); err == nil {
			filter["has_unread"] = hasUnread
		}
	}

	// Group chat filter
	if isGroupStr := c.Query("is_group"); isGroupStr != "" {
		if isGroup, err := strconv.ParseBool(isGroupStr); err == nil {
			filter["is_group"] = isGroup
		}
	}

	page, err := chatSvc.FetchChats(c.Context(), companyId, filter, limit, offset)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessWithTotal(c, page.Items, page.Total)
}

// FetchRangeChats handles GET /chats/range?start=...&end=...&<filters>
// Used for infinite scroll implementation - now returns total count like other endpoints
func FetchRangeChats(c *fiber.Ctx) error {
	companyId := c.Locals("companyId").(string)
	start := c.QueryInt("start", 0)
	end := c.QueryInt("end", start)

	// Build filter map (same as FetchChats)
	filter := make(map[string]interface{})

	if agentId := c.Query("agent_id"); agentId != "" {
		filter["agent_id"] = agentId
	}

	if assignedTo := c.Query("assigned_to"); assignedTo != "" {
		filter["assigned_to"] = assignedTo
	}

	if hasUnreadStr := c.Query("has_unread"); hasUnreadStr != "" {
		if hasUnread, err := strconv.ParseBool(hasUnreadStr); err == nil {
			filter["has_unread"] = hasUnread
		}
	}

	if isGroupStr := c.Query("is_group"); isGroupStr != "" {
		if isGroup, err := strconv.ParseBool(isGroupStr); err == nil {
			filter["is_group"] = isGroup
		}
	}

	page, err := chatSvc.FetchRangeChats(c.Context(), companyId, filter, start, end)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessWithTotal(c, page.Items, page.Total)
}

// SearchChats handles GET /chats/search?q=query
// Searches in: phone_number, push_name, group_name (from chats) and custom_name (from contacts)
func SearchChats(c *fiber.Ctx) error {
	companyId := c.Locals("companyId").(string)
	q := c.Query("q")
	agentId := c.Query("agent_id") // Optional filter by agent

	if q == "" {
		return utils.SuccessWithTotal(c, []any{}, 0)
	}

	page, err := chatSvc.SearchChats(c.Context(), companyId, q, agentId)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessWithTotal(c, page.Items, page.Total)
}
