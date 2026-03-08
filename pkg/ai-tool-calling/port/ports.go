package port

import (
	"context"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/ai-tool-calling/domain"
)

type ToolCallingSrv interface {
	SelectTool(ctx context.Context, prompt string) (*domain.ToolSelectionResult, error)
}

type ToolCallingAdt interface {
	SelectTool(ctx context.Context, prompt string) (*domain.ToolSelectionResult, error)
}
