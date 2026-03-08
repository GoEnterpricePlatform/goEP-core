package service

import (
	"context"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/ai-tool-calling/domain"
)

func (s *Service) SelectTool(ctx context.Context, prompt string) (*domain.ToolSelectionResult, error) {
	resp, err := s.ToolCallingAdt.SelectTool(ctx, prompt)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
