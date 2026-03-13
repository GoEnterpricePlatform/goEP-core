package port

import (
	"context"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/ai-tool-calling/domain"
)

type ToolCallingSrv interface {
	SelectTool(ctx context.Context, prompt string) error
	GetAll(ctx context.Context) ([]*domain.ChatMessage,error)
}

type ToolCallingAdt interface {
	SelectTool(ctx context.Context, messages []*domain.ChatMessage) (*domain.ToolSelectionResult, error)
}
