// internal/service/contact.go
package service

import (
	"context"
	"errors"

	"gitlab.com/timkado/api/daisi-rest-postgres/internal/model"
	"gitlab.com/timkado/api/daisi-rest-postgres/internal/repository"
)

type ContactService interface {
	FetchContacts(ctx context.Context, companyId string, filter map[string]interface{}, sort, order string, limit, offset int) (*model.ContactPage, error)
	GetContactByID(ctx context.Context, companyId, id string) (*model.Contact, error)
	GetContactByPhoneAndAgent(ctx context.Context, companyId, phoneNumber, agentId string) (*model.Contact, error)
	UpdateContact(ctx context.Context, companyId, id string, in model.ContactUpdateInput) (*model.Contact, error)
	SearchContacts(ctx context.Context, companyId, query, agentId string) (*model.ContactPage, error)
}

func NewContactService(repo repository.ContactRepository) ContactService {
	return &contactService{repo: repo}
}

type contactService struct {
	repo repository.ContactRepository
}

func (s *contactService) FetchContacts(
	ctx context.Context,
	companyId string,
	filter map[string]interface{},
	sort, order string,
	limit, offset int,
) (*model.ContactPage, error) {
	if companyId == "" {
		return nil, errors.New("companyId is required")
	}

	// Apply default pagination
	if limit <= 0 {
		limit = 20
	} else if limit > 100 {
		limit = 100
	}

	if offset < 0 {
		offset = 0
	}

	// Default sort
	if sort == "" {
		sort = "created_at"
	}

	if order == "" {
		order = "DESC"
	}

	// Validate filter values
	validatedFilter := make(map[string]interface{})
	for key, value := range filter {
		switch key {
		case "phone_number", "agent_id", "assigned_to", "tags", "status", "origin":
			if strVal, ok := value.(string); ok && strVal != "" {
				validatedFilter[key] = strVal
			}
		case "has_chat":
			if boolVal, ok := value.(bool); ok {
				validatedFilter[key] = boolVal
			}
		}
	}

	return s.repo.FetchContacts(ctx, companyId, validatedFilter, sort, order, limit, offset)
}

func (s *contactService) GetContactByID(ctx context.Context, companyId, id string) (*model.Contact, error) {
	if companyId == "" || id == "" {
		return nil, errors.New("companyId and id are required")
	}

	return s.repo.GetContactByID(ctx, companyId, id)
}

func (s *contactService) GetContactByPhoneAndAgent(ctx context.Context, companyId, phoneNumber, agentId string) (*model.Contact, error) {
	if companyId == "" || phoneNumber == "" || agentId == "" {
		return nil, errors.New("companyId, phoneNumber, and agentId are required")
	}

	return s.repo.GetContactByPhoneAndAgent(ctx, companyId, phoneNumber, agentId)
}

func (s *contactService) UpdateContact(ctx context.Context, companyId, id string, in model.ContactUpdateInput) (*model.Contact, error) {
	if companyId == "" || id == "" {
		return nil, errors.New("companyId and id are required")
	}

	// Build updates map - only include non-empty values
	updates := make(map[string]interface{})

	// Allow empty string to clear values
	if in.CustomName != nil {
		updates["custom_name"] = *in.CustomName
	}

	if in.AssignedTo != nil {
		updates["assigned_to"] = *in.AssignedTo
	}

	if in.Tags != nil {
		updates["tags"] = *in.Tags
	}

	if in.Avatar != nil {
		updates["avatar"] = *in.Avatar
	}

	if in.Notes != nil {
		updates["notes"] = *in.Notes
	}

	if len(updates) == 0 {
		return nil, errors.New("no fields to update")
	}

	return s.repo.UpdateContact(ctx, companyId, id, updates)
}

func (s *contactService) SearchContacts(ctx context.Context, companyId, query, agentId string) (*model.ContactPage, error) {
	if companyId == "" {
		return nil, errors.New("companyId is required")
	}

	if query == "" {
		return &model.ContactPage{
			Items: []map[string]interface{}{},
			Total: 0,
		}, nil
	}

	// Limit query length
	if len(query) > 100 {
		query = query[:100]
	}

	// Default limit for search
	limit := 50

	return s.repo.SearchContacts(ctx, companyId, query, agentId, limit)
}
