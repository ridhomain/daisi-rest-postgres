// internal/repository/message.go
package repository

import (
	"context"
	"fmt"

	"gitlab.com/timkado/api/daisi-rest-postgres/internal/database"
	"gorm.io/gorm"
)

// MessagePage holds a page of messages plus the exact total count
type MessagePage struct {
	Items []map[string]interface{} `json:"items"`
	Total int64                    `json:"total"`
}

// MessageRepository defines read operations on a tenant's partitioned messages table
type MessageRepository interface {
	// FetchMessagesByChatId returns messages for a specific chat with pagination
	FetchMessagesByChatId(ctx context.Context, companyId, agentId, chatId string, limit, offset int) (*MessagePage, error)
	// FetchRangeMessagesByChatId returns messages in [start,end] range for infinite scroll
	FetchRangeMessagesByChatId(ctx context.Context, companyId, agentId, chatId string, start, end int) ([]map[string]interface{}, error)
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
	limit, offset int,
) (*MessagePage, error) {
	// Build base query
	baseQuery := r.buildBaseQuery(ctx, companyId, agentId, chatId)

	// Get total count
	var total int64
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count messages: %w", err)
	}

	// Fetch messages with pagination
	query := r.buildBaseQuery(ctx, companyId, agentId, chatId).
		Order("message_timestamp DESC").
		Limit(limit).
		Offset(offset)

	var items []map[string]interface{}
	if err := query.Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch messages: %w", err)
	}

	if items == nil {
		items = make([]map[string]interface{}, 0)
	}

	return &MessagePage{Items: items, Total: total}, nil
}

func (r *messageRepo) FetchRangeMessagesByChatId(
	ctx context.Context,
	companyId, agentId, chatId string,
	start, end int,
) ([]map[string]interface{}, error) {
	// Calculate limit from range
	limit := end - start + 1
	if limit <= 0 {
		return []map[string]interface{}{}, nil
	}

	// Build query
	query := r.buildBaseQuery(ctx, companyId, agentId, chatId).
		Order("message_timestamp DESC").
		Offset(start).
		Limit(limit)

	var items []map[string]interface{}
	if err := query.Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch range messages: %w", err)
	}

	if items == nil {
		items = make([]map[string]interface{}, 0)
	}

	return items, nil
}
