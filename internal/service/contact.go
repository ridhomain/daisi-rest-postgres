package service

import (
	"context"
	"errors"

	"gitlab.com/timkado/api/daisi-rest-postgres/internal/model"
	"gitlab.com/timkado/api/daisi-rest-postgres/internal/repository"
)

type ContactService interface {
	FetchContacts(ctx context.Context, companyId string, filter model.ContactFilter, sort, order string, limit, offset int) (*model.ContactPage, error)
	GetContactByID(ctx context.Context, companyId, id string) (*model.Contact, error)
	UpdateContact(ctx context.Context, companyId, id string, in model.ContactUpdateInput) (*model.Contact, error)
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
	filter model.ContactFilter,
	sort, order string,
	limit, offset int,
) (*model.ContactPage, error) {
	if companyId == "" {
		return nil, errors.New("companyId is required")
	}
	return s.repo.FetchContacts(ctx, companyId, filter, sort, order, limit, offset)
}

func (s *contactService) GetContactByID(
	ctx context.Context,
	companyId, id string,
) (*model.Contact, error) {
	if companyId == "" || id == "" {
		return nil, errors.New("companyId and id are required")
	}
	return s.repo.GetContactByID(ctx, companyId, id)
}

func (s *contactService) UpdateContact(
	ctx context.Context,
	companyId, id string,
	in model.ContactUpdateInput,
) (*model.Contact, error) {
	if companyId == "" || id == "" {
		return nil, errors.New("companyId and id are required")
	}
	return s.repo.UpdateContact(ctx, companyId, id, in)
}
