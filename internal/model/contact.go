package model

import (
	"time"
)

type Contact struct {
	ID                    string     `json:"id" gorm:"primaryKey;type:text"`
	PhoneNumber           string     `json:"phone_number" gorm:"type:text;uniqueIndex:uniq_agent_phone" validate:"required"`
	AgentID               string     `json:"agent_id,omitempty" gorm:"type:text;uniqueIndex:uniq_agent_phone"`
	Type                  string     `json:"type,omitempty" gorm:"type:text"`                  // e.g., PERSONAL, AGENT, OTHER
	CustomName            string     `json:"custom_name,omitempty" gorm:"type:text"`           // Alias or custom label
	Notes                 string     `json:"notes,omitempty" gorm:"type:text"`                 // Freeform notes
	Tags                  string     `json:"tags,omitempty" gorm:"type:text"`                  // Text of tags (comma-separated)
	CompanyID             string     `json:"company_id,omitempty" gorm:"column:company_id"`    // Company / company ID
	Avatar                string     `json:"avatar,omitempty" gorm:"type:text"`                // URL or reference to profile picture
	AssignedTo            string     `json:"assigned_to,omitempty" gorm:"index;type:text"`     // Assigned agent ID (optional)
	Pob                   string     `json:"pob,omitempty" gorm:"type:text"`                   // Place of birth
	Dob                   *time.Time `json:"dob,omitempty" gorm:"type:date"`                   // Date of birth (pointer for nullability)
	Gender                string     `json:"gender,omitempty" gorm:"type:text;default:MALE"`   // MALE or FEMALE (default MALE)
	Origin                string     `json:"origin,omitempty" gorm:"type:text"`                // Origin source (manual, import, etc.)
	PushName              string     `json:"push_name,omitempty" gorm:"type:text"`             // Push name from WA metadata
	Status                string     `json:"status,omitempty" gorm:"type:text;default:ACTIVE"` // ACTIVE or DISABLED (default ACTIVE)
	FirstMessageID        string     `json:"first_message_id,omitempty" gorm:"type:text"`      // Message ID of first message received (nullable)
	FirstMessageTimestamp int64      `json:"first_message_timestamp,omitempty" gorm:"column:first_message_timestamp"`
	CreatedAt             time.Time  `json:"created_at,omitempty" gorm:"autoCreateTime"`
	UpdatedAt             time.Time  `json:"updated_at,omitempty" gorm:"autoUpdateTime"`
}

// ContactFilter - keeping for compatibility but not used in improved implementation
type ContactFilter struct {
	PhoneNumber string
	AgentID     string
	Tags        string
	AssignedTo  string
	Status      string
	Origin      string
}

// ContactUpdateInput with pointer fields to allow partial updates
type ContactUpdateInput struct {
	CustomName *string `json:"custom_name,omitempty"`
	AssignedTo *string `json:"assigned_to,omitempty"`
	Tags       *string `json:"tags,omitempty"`
	Avatar     *string `json:"avatar,omitempty"`
	Notes      *string `json:"notes,omitempty"`
}

type ContactPage struct {
	Total int64                    `json:"total"`
	Items []map[string]interface{} `json:"items"`
}
