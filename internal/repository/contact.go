// internal/repository/contact.go
package repository

import (
	"context"
	"errors"
	"fmt"

	"gitlab.com/timkado/api/daisi-rest-postgres/internal/model"
	"gorm.io/gorm"
)

type ContactRepository interface {
	FetchContacts(ctx context.Context, companyId string, filter model.ContactFilter, sort, order string, limit, offset int) (*model.ContactPage, error)
	GetContactByID(ctx context.Context, companyId, id string) (*model.Contact, error)
	UpdateContact(ctx context.Context, companyId, id string, in model.ContactUpdateInput) (*model.Contact, error)
}

type contactRepo struct {
	db *gorm.DB
}

func NewContactRepository(db *gorm.DB) ContactRepository {
	return &contactRepo{db: db}
}

func (r *contactRepo) tableFor(companyId string) *gorm.DB {
	return r.db.Table(fmt.Sprintf("daisi_%s.contacts c", companyId))
}

func (r *contactRepo) FetchContacts(ctx context.Context, companyId string, filter model.ContactFilter, sort, order string, limit, offset int) (*model.ContactPage, error) {
	db := r.tableFor(companyId).
		Select("c.*, ch.push_name, ch.group_name, ch.is_group").
		Joins("LEFT JOIN daisi_" + companyId + ".chats ch ON c.phone_number = ch.phone_number AND c.agent_id = ch.agent_id").
		WithContext(ctx)

	if filter.PhoneNumber != "" {
		db = db.Where("c.phone_number = ?", filter.PhoneNumber)
	}
	if filter.AgentID != "" {
		db = db.Where("c.agent_id = ?", filter.AgentID)
	}
	if filter.Tags != "" {
		db = db.Where("c.tags ILIKE ?", "%"+filter.Tags+"%")
	}
	if filter.AssignedTo != "" {
		db = db.Where("c.assigned_to = ?", filter.AssignedTo)
	}
	if filter.Status != "" {
		db = db.Where("c.status = ?", filter.Status)
	}
	if filter.Origin != "" {
		db = db.Where("c.origin = ?", filter.Origin)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	db = db.Order(fmt.Sprintf("c.%s %s", sort, order)).Limit(limit).Offset(offset)
	if err := db.Find(&results).Error; err != nil {
		return nil, err
	}

	return &model.ContactPage{Total: total, Items: results}, nil
}

func (r *contactRepo) GetContactByID(ctx context.Context, companyId, id string) (*model.Contact, error) {
	var contact model.Contact
	err := r.tableFor(companyId).
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
	db := r.tableFor(companyId).WithContext(ctx)
	if err := db.Where("id = ?", id).First(&contact).Error; err != nil {
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
