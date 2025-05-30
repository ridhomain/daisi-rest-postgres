// internal/handler/contact.go
package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/timkado/api/daisi-rest-postgres/internal/model"
	"gitlab.com/timkado/api/daisi-rest-postgres/internal/service"
	"gitlab.com/timkado/api/daisi-rest-postgres/pkg/utils"
)

var contactSvc service.ContactService

func RegisterContactService(svc service.ContactService) {
	contactSvc = svc
}

// FetchContacts handles GET /contacts?limit=...&offset=...&<filters>
func FetchContacts(c *fiber.Ctx) error {
	companyId := c.Locals("companyId").(string)
	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)
	sort := c.Query("sort", "created_at")
	order := c.Query("order", "desc")

	// Build filter map
	filter := make(map[string]interface{})

	// Common filters
	if phoneNumber := c.Query("phone_number"); phoneNumber != "" {
		filter["phone_number"] = phoneNumber
	}

	if agentId := c.Query("agent_id"); agentId != "" {
		filter["agent_id"] = agentId
	}

	if assignedTo := c.Query("assigned_to"); assignedTo != "" {
		filter["assigned_to"] = assignedTo
	}

	if tags := c.Query("tags"); tags != "" {
		filter["tags"] = tags
	}

	if status := c.Query("status"); status != "" {
		filter["status"] = status
	}

	if origin := c.Query("origin"); origin != "" {
		filter["origin"] = origin
	}

	// Boolean filter for contacts with/without chats
	if hasChatStr := c.Query("has_chat"); hasChatStr != "" {
		if hasChat, err := strconv.ParseBool(hasChatStr); err == nil {
			filter["has_chat"] = hasChat
		}
	}

	page, err := contactSvc.FetchContacts(c.Context(), companyId, filter, sort, order, limit, offset)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessWithTotal(c, page.Items, page.Total)
}

// GetContactByID handles GET /contacts/:id
func GetContactByID(c *fiber.Ctx) error {
	companyId := c.Locals("companyId").(string)
	id := c.Params("id")

	contact, err := contactSvc.GetContactByID(c.Context(), companyId, id)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	if contact == nil {
		return utils.Error(c, fiber.StatusNotFound, "contact not found")
	}

	return utils.Success(c, contact)
}

// GetContactByPhoneAndAgent handles GET /contacts/by-phone?phone_number=...&agent_id=...
func GetContactByPhoneAndAgent(c *fiber.Ctx) error {
	companyId := c.Locals("companyId").(string)
	phoneNumber := c.Query("phone_number")
	agentId := c.Query("agent_id")

	if phoneNumber == "" || agentId == "" {
		return utils.Error(c, fiber.StatusBadRequest, "phone_number and agent_id are required")
	}

	contact, err := contactSvc.GetContactByPhoneAndAgent(c.Context(), companyId, phoneNumber, agentId)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	if contact == nil {
		return utils.Error(c, fiber.StatusNotFound, "contact not found")
	}

	return utils.Success(c, contact)
}

// UpdateContact handles PATCH /contacts/:id
func UpdateContact(c *fiber.Ctx) error {
	companyId := c.Locals("companyId").(string)
	id := c.Params("id")

	var body model.ContactUpdateInput
	if err := c.BodyParser(&body); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, err.Error())
	}

	updated, err := contactSvc.UpdateContact(c.Context(), companyId, id, body)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	if updated == nil {
		return utils.Error(c, fiber.StatusNotFound, "contact not found")
	}

	return utils.Success(c, updated)
}

// SearchContacts handles GET /contacts/search?q=...&agent_id=...
// Searches in: phone_number, custom_name, push_name (from chat)
func SearchContacts(c *fiber.Ctx) error {
	companyId := c.Locals("companyId").(string)
	query := c.Query("q")
	agentId := c.Query("agent_id") // Optional filter

	if query == "" {
		return utils.SuccessWithTotal(c, []any{}, 0)
	}

	page, err := contactSvc.SearchContacts(c.Context(), companyId, query, agentId)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessWithTotal(c, page.Items, page.Total)
}
