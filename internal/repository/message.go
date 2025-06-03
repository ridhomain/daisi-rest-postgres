// internal/repository/message.go
package repository

import (
	"context"
	"fmt"

	"gitlab.com/timkado/api/daisi-rest-postgres/internal/database"
	"gitlab.com/timkado/api/daisi-rest-postgres/internal/model"
	"gorm.io/gorm"
)

// MessagePage holds a page of messages plus the exact total count
type MessagePage struct {
	Items []model.Message `json:"items"`
	Total int64           `json:"total"`
}

// MessageRepository defines read operations on a tenant's partitioned messages table
type MessageRepository interface {
	// FetchMessagesByChatId returns messages for a specific chat with pagination
	FetchMessagesByChatId(ctx context.Context, companyId, agentId, chatId string, sort, order string, limit, offset int) (*MessagePage, error)
	// FetchRangeMessagesByChatId returns messages in [start,end] range for infinite scroll
	FetchRangeMessagesByChatId(ctx context.Context, companyId, agentId, chatId string, sort, order string, start, end int) (*MessagePage, error)
}

func NewMessageRepository() MessageRepository {
	return &messageRepo{db: database.DB}
}

type messageRepo struct {
	db *gorm.DB
}

// messageTable returns the fully-qualified, quoted parent table name
func (r *messageRepo) messageTable(companyId string) string {
	schema := "daisi_" + companyId
	return fmt.Sprintf(`"%s"."%s"`, schema, "messages")
}

func (r *messageRepo) validateSort(sort, order string) (string, string) {
	// Allowed sort fields
	allowedSortFields := map[string]bool{
		"message_timestamp": true,
		// "created_at":        true,
		// "updated_at":        true,
		// "from_phone":        true,
		// "to_phone":          true,
		// "message_type":      true,
		// "flow":              true,
	}

	// Default sort
	if sort == "" || !allowedSortFields[sort] {
		sort = "message_timestamp"
	}

	// Validate order
	if order != "ASC" && order != "asc" {
		order = "DESC"
	}

	return sort, order
}

// buildBaseQuery creates the base query for messages
func (r *messageRepo) buildBaseQuery(ctx context.Context, companyId, agentId, chatId string) *gorm.DB {
	tbl := r.messageTable(companyId)

	return r.db.
		Table(tbl).
		WithContext(ctx).
		Where("agent_id = ?", agentId).
		Where("chat_id = ?", chatId).
		Where("key IS NOT NULL") // Only get valid messages
}

func (r *messageRepo) FetchMessagesByChatId(
	ctx context.Context,
	companyId, agentId, chatId string,
	sort, order string,
	limit, offset int,
) (*MessagePage, error) {
	// Validate sort parameters
	sort, order = r.validateSort(sort, order)

	// Build base query
	baseQuery := r.buildBaseQuery(ctx, companyId, agentId, chatId)

	// Get total count
	var total int64
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count messages: %w", err)
	}

	// Fetch messages with pagination
	query := r.buildBaseQuery(ctx, companyId, agentId, chatId).
		Order(fmt.Sprintf("%s %s", sort, order)).
		Limit(limit).
		Offset(offset)

	var items []model.Message
	if err := query.Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch messages: %w", err)
	}

	if items == nil {
		items = make([]model.Message, 0)
	}

	return &MessagePage{Items: items, Total: total}, nil
}

func (r *messageRepo) FetchRangeMessagesByChatId(
	ctx context.Context,
	companyId, agentId, chatId string,
	sort, order string,
	start, end int,
) (*MessagePage, error) {
	// Validate sort parameters
	sort, order = r.validateSort(sort, order)

	// Calculate limit from range
	limit := end - start + 1
	if limit <= 0 {
		return &MessagePage{Items: []model.Message{}, Total: 0}, nil
	}

	// Build query
	baseQuery := r.buildBaseQuery(ctx, companyId, agentId, chatId)

	// Get total count
	var total int64
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count messages: %w", err)
	}

	// Build query for range
	query := r.buildBaseQuery(ctx, companyId, agentId, chatId).
		Order(fmt.Sprintf("%s %s", sort, order)).
		Offset(start).
		Limit(limit)

	var items []model.Message
	if err := query.Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch range messages: %w", err)
	}

	if items == nil {
		items = make([]model.Message, 0)
	}

	return &MessagePage{Items: items, Total: total}, nil
}
