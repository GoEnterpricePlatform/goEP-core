package service

import (
	"context"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/ai-tool-calling/domain"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/ai-tool-calling/repository/memory"
)

func (s *Service) GetAll(ctx context.Context) ([]*domain.ChatMessage, error) {
	messages := memory.Messages

	return messages, nil
}
