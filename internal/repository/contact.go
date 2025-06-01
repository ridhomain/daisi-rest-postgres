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
	FetchContacts(ctx context.Context, companyId string, filter map[string]interface{}, sort, order string, limit, offset int) (*model.ContactPage, error)
	GetContactByID(ctx context.Context, companyId, id string) (*model.Contact, error)
	GetContactByPhoneAndAgent(ctx context.Context, companyId, phoneNumber, agentId string) (*model.Contact, error)
	UpdateContact(ctx context.Context, companyId, id string, updates map[string]interface{}) (*model.Contact, error)
	SearchContacts(ctx context.Context, companyId string, query, agentId string, limit int) (*model.ContactPage, error)
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

// buildBaseQuery creates the base query with optional chat join
func (r *contactRepo) buildBaseQuery(ctx context.Context, companyId string, includeChat bool) *gorm.DB {
	contactTbl := r.contactTable(companyId)

	query := r.db.
		Table(contactTbl + " c").
		WithContext(ctx)

	if includeChat {
		chatTbl := r.chatTable(companyId)

		// LEFT JOIN using chat_id for direct relationship
		joinSQL := fmt.Sprintf(
			"LEFT JOIN %s ch ON c.chat_id = ch.chat_id",
			chatTbl,
		)

		// Select contact fields plus chat fields
		selectFields := []string{
			"c.*",
			"ch.push_name AS chat_push_name",
			"ch.group_name AS chat_group_name",
			"ch.is_group AS chat_is_group",
			"ch.conversation_timestamp AS last_conversation_timestamp",
			"CASE WHEN ch.chat_id IS NULL THEN FALSE ELSE TRUE END AS has_chat",
		}

		query = query.Joins(joinSQL).Select(selectFields)
	}

	return query
}

// applyFilters applies common filters for contacts
func (r *contactRepo) applyFilters(query *gorm.DB, filter map[string]interface{}) *gorm.DB {
	for key, value := range filter {
		switch key {
		case "phone_number":
			query = query.Where("c.phone_number = ?", value)
		case "agent_id":
			query = query.Where("c.agent_id = ?", value)
		case "assigned_to":
			query = query.Where("c.assigned_to = ?", value)
		case "tags":
			// Exact tag match - handle comma-separated tags properly
			if tagStr, ok := value.(string); ok && tagStr != "" {
				// Use regex to match exact tag with word boundaries
				// This ensures TAG1 doesn't match TAG11
				query = query.Where("c.tags ~ ?", `(^|,)`+tagStr+`(,|$)`)
			}
		case "status":
			query = query.Where("c.status = ?", value)
		case "origin":
			query = query.Where("c.origin = ?", value)
		case "has_chat":
			// Filter contacts that have/don't have associated chats
			if hasChatBool, ok := value.(bool); ok {
				if hasChatBool {
					query = query.Where("ch.chat_id IS NOT NULL")
				} else {
					query = query.Where("ch.chat_id IS NULL")
				}
			}
		}
	}
	return query
}

func (r *contactRepo) FetchContacts(
	ctx context.Context,
	companyId string,
	filter map[string]interface{},
	sort, order string,
	limit, offset int,
) (*model.ContactPage, error) {
	// Build query with chat join
	query := r.buildBaseQuery(ctx, companyId, true)

	// Apply filters
	query = r.applyFilters(query, filter)

	// Count total before pagination
	var total int64
	countQuery := *query
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count contacts: %w", err)
	}

	// Validate and apply sorting
	allowedSortFields := map[string]bool{
		"created_at":                  true,
		"updated_at":                  true,
		"custom_name":                 true,
		"phone_number":                true,
		"last_conversation_timestamp": true,
	}

	if !allowedSortFields[sort] {
		sort = "created_at"
	}

	if order != "ASC" && order != "asc" {
		order = "DESC"
	}

	// Handle special sort fields that come from join
	sortField := sort
	if sort == "last_conversation_timestamp" {
		sortField = "last_conversation_timestamp"
	} else {
		sortField = "c." + sort
	}

	query = query.Order(fmt.Sprintf("%s %s", sortField, order))

	// Apply pagination
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	// Fetch data
	var items []model.Contact
	if err := query.Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch contacts: %w", err)
	}

	if items == nil {
		items = make([]model.Contact, 0)
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

func (r *contactRepo) GetContactByPhoneAndAgent(ctx context.Context, companyId, phoneNumber, agentId string) (*model.Contact, error) {
	var contact model.Contact
	err := r.db.
		Table(r.contactTable(companyId)+" c").
		WithContext(ctx).
		Where("c.phone_number = ? AND c.agent_id = ?", phoneNumber, agentId).
		First(&contact).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &contact, err
}

func (r *contactRepo) UpdateContact(ctx context.Context, companyId, id string, updates map[string]interface{}) (*model.Contact, error) {
	var contact model.Contact
	db := r.db.
		Table(r.contactTable(companyId) + " c").
		WithContext(ctx)

	// First, fetch the existing contact
	if err := db.Where("c.id = ?", id).First(&contact).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	// Apply updates
	if err := db.Model(&contact).Updates(updates).Error; err != nil {
		return nil, err
	}

	// Fetch updated contact
	if err := db.Where("c.id = ?", id).First(&contact).Error; err != nil {
		return nil, err
	}

	return &contact, nil
}

func (r *contactRepo) SearchContacts(ctx context.Context, companyId string, query, agentId string, limit int) (*model.ContactPage, error) {
	if query == "" {
		return &model.ContactPage{Items: []model.Contact{}, Total: 0}, nil
	}

	contactTbl := r.contactTable(companyId)
	chatTbl := r.chatTable(companyId)

	// Build search query with chat_id join
	db := r.db.
		Table(contactTbl + " c").
		WithContext(ctx).
		Joins(fmt.Sprintf(
			"LEFT JOIN %s ch ON c.chat_id = ch.chat_id",
			chatTbl,
		))

	// Search pattern
	searchPattern := "%" + query + "%"

	// Search in multiple fields
	searchConditions := `
		c.phone_number ILIKE ? OR 
		c.custom_name ILIKE ? OR 
		ch.push_name ILIKE ?
	`

	db = db.Where(searchConditions, searchPattern, searchPattern, searchPattern)

	// Apply agent filter if provided
	if agentId != "" {
		db = db.Where("c.agent_id = ?", agentId)
	}

	// Select fields
	selectFields := []string{
		"c.*",
		"ch.push_name AS chat_push_name",
		"ch.conversation_timestamp AS last_conversation_timestamp",
		"CASE WHEN ch.chat_id IS NULL THEN FALSE ELSE TRUE END AS has_chat",
	}

	db = db.Select(selectFields).
		Order("c.created_at DESC").
		Limit(limit)

	// Fetch data
	var items []model.Contact
	if err := db.Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to search contacts: %w", err)
	}

	if items == nil {
		items = make([]model.Contact, 0)
	}

	return &model.ContactPage{Items: items, Total: int64(len(items))}, nil
}
