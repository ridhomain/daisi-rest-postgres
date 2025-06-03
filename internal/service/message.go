// internal/service/message.go
package service

import (
	"context"
	"errors"
	"fmt"

	"gitlab.com/timkado/api/daisi-rest-postgres/internal/model"
	"gitlab.com/timkado/api/daisi-rest-postgres/internal/repository"
)

// MessageService defines business operations for reading messages
type MessageService interface {
	// FetchMessagesByChatId returns paginated messages for a chat with sorting
	FetchMessagesByChatId(ctx context.Context, companyId, agentId, chatId string, sort, order string, limit, offset int) (*repository.MessagePage, error)
	// FetchRangeMessagesByChatId returns messages in a specific range for infinite scroll with total count
	FetchRangeMessagesByChatId(ctx context.Context, companyId, agentId, chatId string, sort, order string, start, end int) (*repository.MessagePage, error)
}

// NewMessageService constructs a MessageService backed by the given repository
func NewMessageService(repo repository.MessageRepository) MessageService {
	return &messageService{repo: repo}
}

type messageService struct {
	repo repository.MessageRepository
}

func (s *messageService) FetchMessagesByChatId(
	ctx context.Context,
	companyId, agentId, chatId string,
	sort, order string,
	limit, offset int,
) (*repository.MessagePage, error) {
	// Validate required parameters
	if companyId == "" || agentId == "" || chatId == "" {
		return nil, errors.New("companyId, agentId, and chatId are required")
	}

	// Apply default pagination
	if limit <= 0 {
		limit = 20
	} else if limit > 100 {
		limit = 100 // Cap at 100 messages per request
	}

	if offset < 0 {
		offset = 0
	}

	// Default sort parameters
	if sort == "" {
		sort = "message_timestamp"
	}
	if order == "" {
		order = "DESC"
	}

	// Fetch messages from repository
	page, err := s.repo.FetchMessagesByChatId(ctx, companyId, agentId, chatId, sort, order, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch messages: %w", err)
	}

	// Ensure items is never nil
	if page.Items == nil {
		page.Items = make([]model.Message, 0)
	}

	return page, nil
}

func (s *messageService) FetchRangeMessagesByChatId(
	ctx context.Context,
	companyId, agentId, chatId string,
	sort, order string,
	start, end int,
) (*repository.MessagePage, error) {
	// Validate required parameters
	if companyId == "" || agentId == "" || chatId == "" {
		return nil, errors.New("companyId, agentId, and chatId are required")
	}

	// Validate range parameters
	if start < 0 {
		start = 0
	}

	if end < start {
		end = start
	}

	// Limit range size to prevent excessive data retrieval
	maxRangeSize := 100
	if end-start+1 > maxRangeSize {
		end = start + maxRangeSize - 1
	}

	// Default sort parameters
	if sort == "" {
		sort = "message_timestamp"
	}
	if order == "" {
		order = "DESC"
	}

	// Fetch messages from repository
	page, err := s.repo.FetchRangeMessagesByChatId(ctx, companyId, agentId, chatId, sort, order, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch range messages: %w", err)
	}

	// Ensure items is never nil
	if page.Items == nil {
		page.Items = make([]model.Message, 0)
	}

	return page, nil
}
