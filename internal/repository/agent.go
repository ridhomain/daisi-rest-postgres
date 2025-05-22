package repository

import (
	"context"
	"errors"
	"fmt"

	"gitlab.com/timkado/api/daisi-rest-postgres/internal/database"
	"gitlab.com/timkado/api/daisi-rest-postgres/internal/model"
	"gorm.io/gorm"
)

// AgentRepository defines CRUD operations on agents within a tenant schema.
type AgentRepository interface {
	GetByAgentID(ctx context.Context, companyId, agentId string) (*model.Agent, error)
	ListByCompanyID(ctx context.Context, companyId string) ([]*model.Agent, error)
	ListByAgentIDs(ctx context.Context, companyId string, agentIds []string) ([]*model.Agent, error)
	Create(ctx context.Context, companyId string, a *model.Agent) (*model.Agent, error)
	UpdateName(ctx context.Context, companyId, agentId, newName string) (*model.Agent, error)
	Delete(ctx context.Context, companyId, id string) error
}

// NewAgentRepository returns the GORM-backed implementation.
func NewAgentRepository() AgentRepository {
	return &agentRepo{db: database.DB}
}

type agentRepo struct {
	db *gorm.DB
}

// tableFor scopes all queries to the tenant’s schema.
func (r *agentRepo) tableFor(companyId string) *gorm.DB {
	schemaTable := fmt.Sprintf("daisi_%s.agents", companyId)
	// Use a fresh session to avoid cross-request state
	return r.db.Session(&gorm.Session{}).Table(schemaTable)
}

func (r *agentRepo) GetByAgentID(ctx context.Context, companyId, agentId string) (*model.Agent, error) {
	var a model.Agent
	err := r.
		tableFor(companyId).
		WithContext(ctx).
		Where("agent_id = ?", agentId).
		First(&a).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &a, err
}

func (r *agentRepo) ListByCompanyID(ctx context.Context, companyId string) ([]*model.Agent, error) {
	var agents []*model.Agent
	if err := r.
		tableFor(companyId).
		WithContext(ctx).
		Find(&agents).
		Error; err != nil {
		return nil, err
	}
	return agents, nil
}

func (r *agentRepo) ListByAgentIDs(ctx context.Context, companyId string, agentIds []string) ([]*model.Agent, error) {
	var agents []*model.Agent
	db := r.tableFor(companyId).WithContext(ctx)

	if len(agentIds) > 0 {
		db = db.Where("agent_id IN ?", agentIds)
	}
	if err := db.Find(&agents).Error; err != nil {
		return nil, err
	}
	return agents, nil
}

func (r *agentRepo) Create(ctx context.Context, companyId string, a *model.Agent) (*model.Agent, error) {
	if err := r.
		tableFor(companyId).
		WithContext(ctx).
		Create(a).Error; err != nil {
		return nil, err
	}
	return a, nil
}

func (r *agentRepo) UpdateName(ctx context.Context, companyId, agentId, newName string) (*model.Agent, error) {
	// Fetch existing
	var a model.Agent
	err := r.
		tableFor(companyId).
		WithContext(ctx).
		Where("agent_id = ?", agentId).
		First(&a).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Update name
	a.AgentName = newName
	if err := r.
		tableFor(companyId).
		WithContext(ctx).
		Model(&a).
		Update("agent_name", newName).
		Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *agentRepo) Delete(ctx context.Context, companyId, id string) error {
	return r.
		tableFor(companyId).
		WithContext(ctx).
		Delete(&model.Agent{}, "id = ?", id).
		Error
}
