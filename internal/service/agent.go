package service

import (
	"context"
	"errors"

	"gitlab.com/timkado/api/daisi-rest-postgres/internal/model"
	"gitlab.com/timkado/api/daisi-rest-postgres/internal/repository"
)

// AgentService defines business operations on agents in a tenant schema.
type AgentService interface {
	GetByAgentID(ctx context.Context, companyId, agentId string) (*model.Agent, error)
	ListByCompanyID(ctx context.Context, companyId string) ([]*model.Agent, error)
	ListByAgentIDs(ctx context.Context, companyId string, agentIds []string) ([]*model.Agent, error)
	Create(ctx context.Context, companyId string, in *model.Agent) (*model.Agent, error)
	UpdateName(ctx context.Context, companyId, agentId, newName string) (*model.Agent, error)
	Delete(ctx context.Context, companyId, id string) error
}

// NewAgentService wires the repository into the service.
func NewAgentService(repo repository.AgentRepository) AgentService {
	return &agentService{repo: repo}
}

type agentService struct {
	repo repository.AgentRepository
}

func (s *agentService) GetByAgentID(ctx context.Context, companyId, agentId string) (*model.Agent, error) {
	if companyId == "" || agentId == "" {
		return nil, errors.New("companyId and agentId are required")
	}
	return s.repo.GetByAgentID(ctx, companyId, agentId)
}

func (s *agentService) ListByCompanyID(ctx context.Context, companyId string) ([]*model.Agent, error) {
	if companyId == "" {
		return nil, errors.New("companyId is required")
	}
	return s.repo.ListByCompanyID(ctx, companyId)
}

func (s *agentService) ListByAgentIDs(ctx context.Context, companyId string, agentIds []string) ([]*model.Agent, error) {
	if companyId == "" {
		return nil, errors.New("companyId is required")
	}
	return s.repo.ListByAgentIDs(ctx, companyId, agentIds)
}

func (s *agentService) Create(ctx context.Context, companyId string, a *model.Agent) (*model.Agent, error) {
	if companyId == "" || a.AgentID == "" {
		return nil, errors.New("companyId, id, and agentId are required")
	}
	// enforce tenant
	a.CompanyID = companyId
	return s.repo.Create(ctx, companyId, a)
}

func (s *agentService) UpdateName(ctx context.Context, companyId, agentId, newName string) (*model.Agent, error) {
	if companyId == "" || agentId == "" {
		return nil, errors.New("companyId and agentId are required")
	}
	if newName == "" {
		return nil, errors.New("newName is required")
	}
	return s.repo.UpdateName(ctx, companyId, agentId, newName)
}

func (s *agentService) Delete(ctx context.Context, companyId, id string) error {
	if companyId == "" || id == "" {
		return errors.New("companyId and id are required")
	}
	return s.repo.Delete(ctx, companyId, id)
}
