// internal/service/message.go
package service

import (
	"context"
	"errors"

	"gitlab.com/timkado/api/daisi-rest-postgres/internal/repository"
	"gitlab.com/timkado/api/daisi-rest-postgres/pkg/logger"
	"go.uber.org/zap"
)

// MessageService defines business operations for reading messages.
type MessageService interface {
	// FetchMessagesByChatId returns up to `limit` messages plus the total count.
	FetchMessagesByChatId(ctx context.Context, companyId, agentId, chatId string, limit int) (*repository.MessagePage, error)
	// FetchRangeMessagesByChatId returns messages in [start,end] for a given chat.
	FetchRangeMessagesByChatId(ctx context.Context, companyId, agentId, chatId string, start, end int) ([]map[string]interface{}, error)
}

// NewMessageService constructs a MessageService backed by the given repository.
func NewMessageService(repo repository.MessageRepository) MessageService {
	return &messageService{repo: repo}
}

type messageService struct {
	repo repository.MessageRepository
}

func (s *messageService) FetchMessagesByChatId(
	ctx context.Context,
	companyId, agentId, chatId string,
	limit int,
) (*repository.MessagePage, error) {
	if companyId == "" || agentId == "" || chatId == "" {
		logger.NewLogger().Info("Payload", zap.String("companyId", companyId), zap.String("agentId", agentId), zap.String("chatId", chatId))
		return nil, errors.New("companyId, agentId, and chatId are required")
	}
	return s.repo.FetchMessagesByChatId(ctx, companyId, agentId, chatId, limit)
}

func (s *messageService) FetchRangeMessagesByChatId(
	ctx context.Context,
	companyId, agentId, chatId string,
	start, end int,
) ([]map[string]interface{}, error) {
	if companyId == "" || agentId == "" || chatId == "" {
		return nil, errors.New("companyId, agentId, and chatId are required")
	}
	return s.repo.FetchRangeMessagesByChatId(ctx, companyId, agentId, chatId, start, end)
}
