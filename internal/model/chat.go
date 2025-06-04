package model

import (
	"time"

	"gorm.io/datatypes"
)

type Chat struct {
	ID                    int64          `json:"-" gorm:"primaryKey;autoIncrement"`
	ChatID                string         `json:"id" gorm:"column:chat_id;uniqueIndex"`
	Jid                   string         `json:"jid" gorm:"column:jid;index"`
	PushName              string         `json:"push_name" gorm:"column:push_name"`
	IsGroup               bool           `json:"is_group" gorm:"column:is_group"`
	GroupName             string         `json:"group_name" gorm:"column:group_name"`
	UnreadCount           int32          `json:"unread_count" gorm:"column:unread_count"`
	LastMessageObj        datatypes.JSON `json:"last_message" gorm:"type:jsonb;column:last_message"`
	ConversationTimestamp int64          `json:"conversation_timestamp" gorm:"column:conversation_timestamp"`
	NotSpam               bool           `json:"not_spam,omitempty" gorm:"column:not_spam"`
	AgentID               string         `json:"agent_id" gorm:"column:agent_id;index"`
	CompanyID             string         `json:"company_id" gorm:"column:company_id"`
	PhoneNumber           string         `json:"phone_number" gorm:"column:phone_number"`
	CreatedAt             time.Time      `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt             time.Time      `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}
