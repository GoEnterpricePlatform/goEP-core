package disabled

import (
	"context"
	"errors"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/ai-tool-calling/domain"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/ai-tool-calling/initializer"
)

type DisabledAdapter struct{}

func NewDisabledAdapter() *DisabledAdapter {
	return &DisabledAdapter{}
}

func (a *DisabledAdapter) SetTools(initializer *initializer.Initializer) {
	// I don't need to add the tools because it is not a provider
}

func (a *DisabledAdapter) SelectTool(ctx context.Context, messages []*domain.ChatMessage, systemPrompt string) (*domain.ToolSelectionResult, error){
	// if we want to show a custom error to the user we can create an ErrToolCallingDisabled,
	// I don't think it's necessary, I just want an internal server error.
	return nil,errors.New("error tool-calling disabled")
}
