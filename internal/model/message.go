package model

import (
	"time"

	"gorm.io/datatypes"
)

// Message represents a chat message with fields aligned to proto definitions
type Message struct {
	ID               int64          `json:"-" gorm:"column:id;primaryKey;autoIncrement"`
	MessageID        string         `json:"id" gorm:"column:message_id;index"`
	FromPhone        string         `json:"from_phone,omitempty" gorm:"column:from_phone;index"`
	ToPhone          string         `json:"to_phone,omitempty" gorm:"column:to_phone;index"`
	ChatID           string         `json:"chat_id,omitempty" gorm:"column:chat_id;index"`
	Jid              string         `json:"jid,omitempty" gorm:"column:jid;index"`
	Flow             string         `json:"flow,omitempty" gorm:"column:flow"`
	MessageText      string         `json:"message_text" gorm:"column:message_text"`
	MessageUrl       string         `json:"message_url" gorm:"column:message_url"`
	MessageType      string         `json:"message_type,omitempty" gorm:"column:message_type"`
	AgentID          string         `json:"agent_id,omitempty" gorm:"column:agent_id;index"`
	CompanyID        string         `json:"company_id,omitempty" gorm:"column:company_id"` // CompanyID is implicitly the tenant ID
	MessageObj       datatypes.JSON `json:"message_obj,omitempty" gorm:"type:jsonb;column:message_obj"`
	EditedMessageObj datatypes.JSON `json:"edited_message_obj,omitempty" gorm:"type:jsonb;column:edited_message_obj"`
	Key              datatypes.JSON `json:"key,omitempty" gorm:"type:jsonb;column:key"`
	Status           string         `json:"status,omitempty" gorm:"column:status"`
	IsDeleted        bool           `json:"is_deleted,omitempty" gorm:"column:is_deleted;default:false"`
	MessageTimestamp int64          `json:"message_timestamp,omitempty" gorm:"column:message_timestamp;index"`
	MessageDate      time.Time      `gorm:"column:message_date;type:date;not null"`
	CreatedAt        time.Time      `json:"created_at,omitempty" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt        time.Time      `json:"updated_at,omitempty" gorm:"column:updated_at;autoUpdateTime"`
	// LastMetadata     datatypes.JSON `json:"last_metadata,omitempty" gorm:"type:jsonb;column:last_metadata"`
}
