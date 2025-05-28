// internal/repository/contact.go
package repository

import (
	"context"
	"errors"
	"fmt"

	"gitlab.com/timkado/api/daisi-rest-postgres/internal/database"
	"gitlab.com/timkado/api/daisi-rest-postgres/internal/model"
	"gorm.io/gorm"
)

type ContactRepository interface {
	FetchContacts(ctx context.Context, companyId string, filter model.ContactFilter, sort, order string, limit, offset int) (*model.ContactPage, error)
	GetContactByID(ctx context.Context, companyId, id string) (*model.Contact, error)
	UpdateContact(ctx context.Context, companyId, id string, in model.ContactUpdateInput) (*model.Contact, error)
}

func NewContactRepository() ContactRepository {
	return &contactRepo{db: database.DB}
}

type contactRepo struct {
	db *gorm.DB
}

func (r *contactRepo) contactTable(companyId string) string {
	schema := "daisi_" + companyId
	return fmt.Sprintf(`"%s"."%s"`, schema, "contacts")
}

func (r *contactRepo) chatTable(companyId string) string {
	schema := "daisi_" + companyId
	return fmt.Sprintf(`"%s"."%s"`, schema, "chats")
}

func (r *contactRepo) FetchContacts(
	ctx context.Context,
	companyId string,
	filter model.ContactFilter,
	sort, order string,
	limit, offset int,
) (*model.ContactPage, error) {
	contactTbl := r.contactTable(companyId)
	chatTbl := r.chatTable(companyId)

	base := r.db.
		Table(contactTbl + " c").
		WithContext(ctx)

	// apply filters
	if filter.PhoneNumber != "" {
		base = base.Where("c.phone_number = ?", filter.PhoneNumber)
	}
	if filter.AgentID != "" {
		base = base.Where("c.agent_id = ?", filter.AgentID)
	}
	if filter.Tags != "" {
		base = base.Where("c.tags ILIKE ?", "%"+filter.Tags+"%")
	}
	if filter.AssignedTo != "" {
		base = base.Where("c.assigned_to = ?", filter.AssignedTo)
	}
	if filter.Status != "" {
		base = base.Where("c.status = ?", filter.Status)
	}
	if filter.Origin != "" {
		base = base.Where("c.origin = ?", filter.Origin)
	}

	// count before join
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, err
	}

	// join and select
	joinSQL := fmt.Sprintf(
		"LEFT JOIN %s ch ON c.phone_number = ch.phone_number AND c.agent_id = ch.agent_id",
		chatTbl,
	)

	selectFields := []string{
		"c.*",
		"ch.push_name",
		"ch.group_name",
		"ch.is_group",
	}

	query := base.
		Joins(joinSQL).
		Select(selectFields).
		Order(fmt.Sprintf("c.%s %s", sort, order)).
		Limit(limit).
		Offset(offset)

	var items []map[string]interface{}
	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}

	return &model.ContactPage{Total: total, Items: items}, nil
}

func (r *contactRepo) GetContactByID(ctx context.Context, companyId, id string) (*model.Contact, error) {
	var contact model.Contact
	err := r.db.
		Table(r.contactTable(companyId)+" c").
		WithContext(ctx).
		Where("c.id = ?", id).
		First(&contact).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &contact, err
}

func (r *contactRepo) UpdateContact(ctx context.Context, companyId, id string, in model.ContactUpdateInput) (*model.Contact, error) {
	var contact model.Contact
	db := r.db.
		Table(r.contactTable(companyId) + " c").
		WithContext(ctx)

	if err := db.Where("c.id = ?", id).First(&contact).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	updates := map[string]interface{}{
		"custom_name": in.CustomName,
		"assigned_to": in.AssignedTo,
		"tags":        in.Tags,
	}

	if err := db.Model(&contact).Updates(updates).Error; err != nil {
		return nil, err
	}
	return &contact, nil
}
