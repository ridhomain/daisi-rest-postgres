// internal/repository/message.go
package repository

import (
	"context"
	"fmt"

	"gitlab.com/timkado/api/daisi-rest-postgres/internal/database"
	"gorm.io/gorm"
)

// MessagePage holds a page of messages plus the exact total count.
type MessagePage struct {
	Items []map[string]interface{} `json:"items"`
	Total int64                    `json:"total"`
}

// MessageRepository defines read operations on a tenant's partitioned messages table.
type MessageRepository interface {
	// FetchMessagesByChatId returns up to `limit` messages plus total count.
	FetchMessagesByChatId(ctx context.Context, companyId, agentId, chatId string, limit int) (*MessagePage, error)
	// FetchRangeMessagesByChatId returns messages in [start,end] for a given chat.
	FetchRangeMessagesByChatId(ctx context.Context, companyId, agentId, chatId string, start, end int) ([]map[string]interface{}, error)
}

func NewMessageRepository() MessageRepository {
	return &messageRepo{db: database.DB}
}

type messageRepo struct {
	db *gorm.DB
}

// messageTable returns the fully‐qualified, quoted parent table name:
//
//	"daisi_<companyId>"."messages"
func (r *messageRepo) messageTable(companyId string) string {
	schema := "daisi_" + companyId
	return fmt.Sprintf(`"%s"."%s"`, schema, "messages")
}

func (r *messageRepo) FetchMessagesByChatId(
	ctx context.Context,
	companyId, agentId, chatId string,
	limit int,
) (*MessagePage, error) {
	tbl := r.messageTable(companyId)

	base := r.db.
		Table(tbl).
		WithContext(ctx).
		Where("agent_id = ?", agentId).
		Where("chat_id   = ?", chatId).
		Where("key IS NOT NULL")

	// Get total count
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, err
	}

	// Fetch up to `limit` rows
	rows := base.
		Order("message_timestamp DESC").
		Limit(limit)

	var items []map[string]interface{}
	if err := rows.Find(&items).Error; err != nil {
		return nil, err
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
	tbl := r.messageTable(companyId)

	db := r.db.
		Table(tbl).
		WithContext(ctx).
		Where("agent_id = ?", agentId).
		Where("chat_id   = ?", chatId).
		Where("key IS NOT NULL").
		Order("message_timestamp DESC").
		Offset(start).
		Limit(end - start + 1)

	var items []map[string]interface{}
	if err := db.Find(&items).Error; err != nil {
		return nil, err
	}
	if items == nil {
		items = make([]map[string]interface{}, 0)
	}
	return items, nil
}
