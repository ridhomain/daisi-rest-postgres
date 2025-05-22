package model

import (
	"time"
)

// Agent represents the agents table structure, containing information about connected WhatsApp agents.
type Agent struct {
	// ID is the internal database primary key.
	ID int64 `json:"-" gorm:"primaryKey;autoIncrement"`
	// AgentID is the unique identifier for the agent (e.g., from the WA client).
	AgentID string `json:"agent_id" gorm:"column:agent_id;uniqueIndex" validate:"required"`
	// QRCode is the QR code content used for linking/pairing the agent, if applicable.
	QRCode string `json:"qr_code" gorm:"column:qr_code"`
	// Status indicates the current connection status of the agent (e.g., 'connected', 'disconnected').
	Status string `json:"status" gorm:"column:status"`
	// AgentName is a user-defined custom label or name for the agent.
	AgentName string `json:"agent_name" gorm:"column:agent_name"`
	// HostName is the name of the device or host machine running the agent instance.
	HostName string `json:"host_name" gorm:"column:host_name"`
	// Version stores the version information of the agent software.
	Version string `json:"version" gorm:"column:version"`
	// CompanyID identifies the company/tenant this agent belongs to.
	CompanyID string `json:"company_id" gorm:"column:company_id"` // CompanyID is implicitly the tenant ID
	// CreatedAt is the timestamp when the agent record was first created.
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	// UpdatedAt is the timestamp when the agent record was last updated.
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}
