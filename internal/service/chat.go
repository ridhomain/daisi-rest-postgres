// internal/service/chat.go
package service

import (
	"context"
	"errors"

	"gitlab.com/timkado/api/daisi-rest-postgres/internal/model"
	"gitlab.com/timkado/api/daisi-rest-postgres/internal/repository"
)

// ChatService defines business operations for chats
type ChatService interface {
	// FetchChats returns a paginated page of chats with total count and joined contact info
	FetchChats(ctx context.Context, companyId string, filter map[string]interface{}, limit, offset int) (*repository.ChatPage, error)
	// FetchRangeChats returns a page of chats with total count and joined contact info
	FetchRangeChats(ctx context.Context, companyId string, filter map[string]interface{}, start, end int) (*repository.ChatPage, error)
	// SearchChats performs a text search across chats and contacts
	SearchChats(ctx context.Context, companyId, query, agentId string) (*repository.ChatPage, error)
}

// NewChatService constructs a ChatService backed by the given repository
func NewChatService(repo repository.ChatRepository) ChatService {
	return &chatService{repo: repo}
}

type chatService struct {
	repo repository.ChatRepository
}

func (s *chatService) FetchChats(
	ctx context.Context,
	companyId string,
	filter map[string]interface{},
	limit, offset int,
) (*repository.ChatPage, error) {
	if companyId == "" {
		return nil, errors.New("companyId is required")
	}

	// Apply default pagination if not specified
	if limit <= 0 {
		limit = 20
	} else if limit > 100 {
		limit = 100 // Cap at 100 to prevent excessive data retrieval
	}

	if offset < 0 {
		offset = 0
	}

	// Validate filter values
	validatedFilter := make(map[string]interface{})
	for key, value := range filter {
		switch key {
		case "agent_id", "assigned_to":
			if strVal, ok := value.(string); ok && strVal != "" {
				validatedFilter[key] = strVal
			}
		case "has_unread", "is_group":
			if boolVal, ok := value.(bool); ok {
				validatedFilter[key] = boolVal
			}
		}
	}

	return s.repo.FetchChats(ctx, companyId, validatedFilter, limit, offset)
}

func (s *chatService) FetchRangeChats(
	ctx context.Context,
	companyId string,
	filter map[string]interface{},
	start, end int,
) (*repository.ChatPage, error) {
	if companyId == "" {
		return nil, errors.New("companyId is required")
	}

	if start < 0 {
		start = 0
	}

	if end < start {
		end = start
	}

	// Limit range to prevent excessive data retrieval
	maxRange := 100
	if end-start+1 > maxRange {
		end = start + maxRange - 1
	}

	// Validate filter values (same as FetchChats)
	validatedFilter := make(map[string]interface{})
	for key, value := range filter {
		switch key {
		case "agent_id", "assigned_to":
			if strVal, ok := value.(string); ok && strVal != "" {
				validatedFilter[key] = strVal
			}
		case "has_unread", "is_group":
			if boolVal, ok := value.(bool); ok {
				validatedFilter[key] = boolVal
			}
		}
	}

	return s.repo.FetchRangeChats(ctx, companyId, validatedFilter, start, end)
}

func (s *chatService) SearchChats(
	ctx context.Context,
	companyId, query, agentId string,
) (*repository.ChatPage, error) {
	if companyId == "" {
		return nil, errors.New("companyId is required")
	}

	// Return empty result if query is empty
	if query == "" {
		return &repository.ChatPage{
			Items: []model.Chat{},
			Total: 0,
		}, nil
	}

	// Limit query length to prevent abuse
	if len(query) > 100 {
		query = query[:100]
	}

	return s.repo.SearchChats(ctx, companyId, query, agentId)
}
