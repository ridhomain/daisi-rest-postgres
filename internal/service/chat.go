// internal/service/chat.go
package service

import (
	"context"
	"errors"

	"gitlab.com/timkado/api/daisi-rest-postgres/internal/repository"
)

// ChatService defines business operations for chats.
type ChatService interface {
	// FetchChats returns a paginated page of chats with total count and joined contact info.
	FetchChats(ctx context.Context, companyId string, filter map[string]interface{}, limit, offset int) (*repository.ChatPage, error)
	// FetchRangeChats returns a slice of messages for a given chat ID, with joined contact info.
	FetchRangeChats(ctx context.Context, companyId, chatId string, start, end int) ([]map[string]interface{}, error)
}

// NewChatService constructs a ChatService backed by the given repository.
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
	// you could validate filter keys here if needed
	return s.repo.FetchChats(ctx, companyId, filter, limit, offset)
}

func (s *chatService) FetchRangeChats(
	ctx context.Context,
	companyId, chatId string,
	start, end int,
) ([]map[string]interface{}, error) {
	if companyId == "" {
		return nil, errors.New("companyId is required")
	}
	if chatId == "" {
		return nil, errors.New("chatId is required")
	}
	return s.repo.FetchRangeChats(ctx, companyId, chatId, start, end)
}
