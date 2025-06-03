// internal/repository/chat.go
package repository

import (
	"context"
	"fmt"
	"strings"

	"gitlab.com/timkado/api/daisi-rest-postgres/internal/database"
	"gitlab.com/timkado/api/daisi-rest-postgres/internal/model"
	"gorm.io/gorm"
)

type ChatPage struct {
	Items []model.Chat `json:"items"`
	Total int64        `json:"total"`
}

type ChatRepository interface {
	FetchChats(ctx context.Context, companyId string, filter map[string]interface{}, limit, offset int) (*ChatPage, error)
	FetchRangeChats(ctx context.Context, companyId string, filter map[string]interface{}, start, end int) (*ChatPage, error)
	SearchChats(ctx context.Context, companyId string, q string, agentId string) (*ChatPage, error)
}

func NewChatRepository() ChatRepository {
	return &chatRepo{db: database.DB}
}

type chatRepo struct {
	db *gorm.DB
}

func (r *chatRepo) chatTable(companyId string) string {
	schemaName := "daisi_" + companyId
	return fmt.Sprintf(`"%s"."%s"`, schemaName, "chats")
}

func (r *chatRepo) contactsTable(companyId string) string {
	schemaName := "daisi_" + companyId
	return fmt.Sprintf(`"%s"."%s"`, schemaName, "contacts")
}

// buildBaseQuery creates the base query with proper LEFT JOIN
// Each chat should ideally have a contact, but not every contact has a chat
func (r *chatRepo) buildBaseQuery(ctx context.Context, companyId string) *gorm.DB {
	chatTbl := r.chatTable(companyId)
	contactsTbl := r.contactsTable(companyId)

	// Join by chat_id
	joinSQL := fmt.Sprintf(
		"LEFT JOIN %s ON %s.chat_id = %s.chat_id",
		contactsTbl, chatTbl, contactsTbl,
	)

	// Select all chat fields plus specific contact fields we need
	selectFields := []string{
		fmt.Sprintf("%s.*", chatTbl),
		// Contact fields that we need from the join
		fmt.Sprintf("%s.custom_name AS contact_custom_name", contactsTbl),
		fmt.Sprintf("%s.assigned_to AS contact_assigned_to", contactsTbl),
		fmt.Sprintf("%s.tags AS contact_tags", contactsTbl),
		fmt.Sprintf("%s.avatar AS contact_avatar", contactsTbl),
		fmt.Sprintf("%s.origin AS contact_origin", contactsTbl),
		// Check if contact exists
		fmt.Sprintf("CASE WHEN %s.id IS NULL THEN FALSE ELSE TRUE END AS has_contact", contactsTbl),
	}

	return r.db.
		Table(chatTbl).
		WithContext(ctx).
		Joins(joinSQL).
		Select(selectFields)
}

// applyFilters handles the most common filters
func (r *chatRepo) applyFilters(query *gorm.DB, filter map[string]interface{}, chatTbl, contactsTbl string) *gorm.DB {
	for key, value := range filter {
		switch key {
		case "agent_id":
			query = query.Where(fmt.Sprintf("%s.agent_id = ?", chatTbl), value)
		case "assigned_to":
			// This filters by contact's assigned_to field
			query = query.Where(fmt.Sprintf("%s.assigned_to = ?", contactsTbl), value)
		case "has_unread":
			// Boolean filter: true for unread_count > 0, false for unread_count = 0
			if hasUnread, ok := value.(bool); ok {
				if hasUnread {
					query = query.Where(fmt.Sprintf("%s.unread_count > ?", chatTbl), 0)
				} else {
					query = query.Where(fmt.Sprintf("%s.unread_count = ?", chatTbl), 0)
				}
			}
		case "is_group":
			query = query.Where(fmt.Sprintf("%s.is_group = ?", chatTbl), value)
		}
	}
	return query
}

// buildCountQuery creates an optimized count query without JOIN
func (r *chatRepo) buildCountQuery(ctx context.Context, companyId string, filter map[string]interface{}) *gorm.DB {
	chatTbl := r.chatTable(companyId)

	countQuery := r.db.
		Table(chatTbl).
		WithContext(ctx)

	// Apply filters for count (only chat table filters for performance)
	for key, value := range filter {
		switch key {
		case "agent_id":
			countQuery = countQuery.Where(fmt.Sprintf("%s.agent_id = ?", chatTbl), value)
		case "has_unread":
			if hasUnread, ok := value.(bool); ok {
				if hasUnread {
					countQuery = countQuery.Where(fmt.Sprintf("%s.unread_count > ?", chatTbl), 0)
				} else {
					countQuery = countQuery.Where(fmt.Sprintf("%s.unread_count = ?", chatTbl), 0)
				}
			}
		case "is_group":
			countQuery = countQuery.Where(fmt.Sprintf("%s.is_group = ?", chatTbl), value)
		}
	}

	return countQuery
}

func (r *chatRepo) FetchChats(
	ctx context.Context,
	companyId string,
	filter map[string]interface{},
	limit, offset int,
) (*ChatPage, error) {
	chatTbl := r.chatTable(companyId)
	contactsTbl := r.contactsTable(companyId)

	// Get total count with optimized query
	var total int64
	countQuery := r.buildCountQuery(ctx, companyId, filter)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count chats: %w", err)
	}

	// Build data query with JOIN
	dataQuery := r.buildBaseQuery(ctx, companyId)
	dataQuery = r.applyFilters(dataQuery, filter, chatTbl, contactsTbl)

	// Always sort by conversation_timestamp DESC (newest first)
	dataQuery = dataQuery.Order(fmt.Sprintf("%s.conversation_timestamp DESC", chatTbl))

	// Apply pagination
	if limit > 0 {
		dataQuery = dataQuery.Limit(limit)
	}
	if offset > 0 {
		dataQuery = dataQuery.Offset(offset)
	}

	// Fetch data
	var items []model.Chat
	if err := dataQuery.Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch chats: %w", err)
	}

	if items == nil {
		items = make([]model.Chat, 0)
	}

	return &ChatPage{Items: items, Total: total}, nil
}

func (r *chatRepo) FetchRangeChats(
	ctx context.Context,
	companyId string,
	filter map[string]interface{},
	start, end int,
) (*ChatPage, error) {
	chatTbl := r.chatTable(companyId)
	contactsTbl := r.contactsTable(companyId)

	// Get total count with optimized query
	var total int64
	countQuery := r.buildCountQuery(ctx, companyId, filter)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count chats: %w", err)
	}

	// Build query with JOIN
	query := r.buildBaseQuery(ctx, companyId)
	query = r.applyFilters(query, filter, chatTbl, contactsTbl)

	// Always sort by conversation_timestamp DESC for range queries
	query = query.Order(fmt.Sprintf("%s.conversation_timestamp DESC", chatTbl))

	// Apply range
	if start >= 0 {
		query = query.Offset(start)
	}

	limit := end - start + 1
	if limit > 0 {
		query = query.Limit(limit)
	}

	// Fetch data
	var items []model.Chat
	if err := query.Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch range chats: %w", err)
	}

	if items == nil {
		items = make([]model.Chat, 0)
	}

	return &ChatPage{Items: items, Total: total}, nil
}

func (r *chatRepo) SearchChats(
	ctx context.Context,
	companyId string,
	query string,
	agentId string,
) (*ChatPage, error) {
	if query == "" {
		return &ChatPage{Items: []model.Chat{}, Total: 0}, nil
	}

	chatTbl := r.chatTable(companyId)
	contactTbl := r.contactsTable(companyId)

	// Build search query with proper escaping
	searchPattern := "%" + strings.ReplaceAll(query, "%", "\\%") + "%"

	db := r.buildBaseQuery(ctx, companyId)

	// Search in chat fields (phone_number, push_name, group_name) and contact field (custom_name)
	searchConditions := fmt.Sprintf(`
		%s.phone_number ILIKE ? OR 
		%s.push_name ILIKE ? OR 
		%s.group_name ILIKE ? OR
		%s.custom_name ILIKE ?
	`, chatTbl, chatTbl, chatTbl, contactTbl)

	db = db.Where(searchConditions, searchPattern, searchPattern, searchPattern, searchPattern)

	// Apply agent_id filter if provided
	if agentId != "" {
		db = db.Where(fmt.Sprintf("%s.agent_id = ?", chatTbl), agentId)
	}

	// Always sort by conversation_timestamp DESC
	db = db.Order(fmt.Sprintf("%s.conversation_timestamp DESC", chatTbl))

	// Limit search results to prevent excessive data
	db = db.Limit(100)

	// Fetch data
	var items []model.Chat
	if err := db.Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to search chats: %w", err)
	}

	if items == nil {
		items = make([]model.Chat, 0)
	}

	return &ChatPage{Items: items, Total: int64(len(items))}, nil
}
