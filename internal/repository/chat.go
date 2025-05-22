package repository

import (
	"context"
	"fmt"

	"gitlab.com/timkado/api/daisi-rest-postgres/internal/database"
	"gorm.io/gorm"
)

type ChatPage struct {
	Items []map[string]interface{}
	Total int64
}

type ChatRepository interface {
	FetchChats(ctx context.Context, companyId string, filter map[string]interface{}, limit, offset int) (*ChatPage, error)
	FetchRangeChats(ctx context.Context, companyId string, start, end int) ([]map[string]interface{}, error)
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

func (r *chatRepo) FetchChats(
	ctx context.Context,
	companyId string,
	filter map[string]interface{},
	limit, offset int,
) (*ChatPage, error) {
	chatTbl := r.chatTable(companyId)
	contactsTbl := r.contactsTable(companyId)

	base := r.db.
		Table(chatTbl).
		WithContext(ctx)

	// apply filters
	for col, val := range filter {
		base = base.Where(fmt.Sprintf("%s = ?", col), val)
	}

	// total count
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, err
	}

	// build join and select fields
	joinSQL := fmt.Sprintf(
		"LEFT JOIN %s ON %s.phone_number = %s.phone_number AND %s.agent_id = %s.agent_id",
		contactsTbl, chatTbl, contactsTbl, chatTbl, contactsTbl,
	)
	joinFields := []string{
		fmt.Sprintf("%s.avatar            AS contact_avatar", contactsTbl),
		fmt.Sprintf("%s.tags              AS contact_tags", contactsTbl),
		fmt.Sprintf("%s.origin            AS contact_origin", contactsTbl),
		fmt.Sprintf("%s.assigned_to       AS contact_assigned_to", contactsTbl),
		fmt.Sprintf("%s.custom_name       AS contact_custom_name", contactsTbl),
		fmt.Sprintf("CASE WHEN %s.id IS NULL THEN FALSE ELSE TRUE END AS has_contact", contactsTbl),
	}
	// select everything from chats plus our join fields
	fields := append([]string{fmt.Sprintf("%s.*", chatTbl)}, joinFields...)

	// fetch page
	rows := base.
		Joins(joinSQL).
		Select(fields).
		Order(fmt.Sprintf("%s.conversation_timestamp DESC", chatTbl)).
		Limit(limit).
		Offset(offset)

	var items []map[string]interface{}
	if err := rows.Find(&items).Error; err != nil {
		return nil, err
	}
	return &ChatPage{Items: items, Total: total}, nil
}

func (r *chatRepo) FetchRangeChats(
	ctx context.Context,
	companyId string,
	start, end int,
) ([]map[string]interface{}, error) {
	chatTbl := r.chatTable(companyId)
	contactsTbl := r.contactsTable(companyId)

	base := r.db.
		Table(chatTbl).
		WithContext(ctx)

	// same join & select fields as above
	joinSQL := fmt.Sprintf(
		"LEFT JOIN %s ON %s.phone_number = %s.phone_number AND %s.agent_id = %s.agent_id",
		contactsTbl, chatTbl, contactsTbl, chatTbl, contactsTbl,
	)
	joinFields := []string{
		fmt.Sprintf("%s.avatar            AS contact_avatar", contactsTbl),
		fmt.Sprintf("%s.tags              AS contact_tags", contactsTbl),
		fmt.Sprintf("%s.origin            AS contact_origin", contactsTbl),
		fmt.Sprintf("%s.assigned_to       AS contact_assigned_to", contactsTbl),
		fmt.Sprintf("%s.custom_name       AS contact_custom_name", contactsTbl),
		fmt.Sprintf("CASE WHEN %s.id IS NULL THEN FALSE ELSE TRUE END AS has_contact", contactsTbl),
	}
	fields := append([]string{fmt.Sprintf("%s.*", chatTbl)}, joinFields...)

	rows := base.
		Joins(joinSQL).
		Select(fields).
		Order(fmt.Sprintf("%s.conversation_timestamp DESC", chatTbl)).
		Offset(start).
		Limit(end - start + 1)

	var items []map[string]interface{}
	if err := rows.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
