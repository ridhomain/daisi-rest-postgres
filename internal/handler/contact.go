// internal/handler/contact.go
package handler

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/timkado/api/daisi-rest-postgres/internal/model"
	"gitlab.com/timkado/api/daisi-rest-postgres/internal/service"
	"gitlab.com/timkado/api/daisi-rest-postgres/pkg/utils"
)

var contactSvc service.ContactService

func RegisterContactService(svc service.ContactService) {
	contactSvc = svc
}

func FetchContacts(c *fiber.Ctx) error {
	companyId := c.Locals("companyId").(string)
	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)
	sort := c.Query("sort", "created_at")
	order := c.Query("order", "desc")

	filter := model.ContactFilter{
		PhoneNumber: c.Query("phone_number"),
		AgentID:     c.Query("agent_id"),
		Tags:        c.Query("tags"),
		AssignedTo:  c.Query("assigned_to"),
		Status:      c.Query("status"),
		Origin:      c.Query("origin"),
	}

	page, err := contactSvc.FetchContacts(c.Context(), companyId, filter, sort, order, limit, offset)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	if page.Items == nil {
		page.Items = []map[string]interface{}{}
	}
	return utils.SuccessWithTotal(c, page.Items, page.Total)
}

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
