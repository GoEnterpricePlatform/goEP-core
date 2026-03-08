package initializer

import (
	"github.com/GoEnterpricePlatform/goEP-core/web/shared/ai/contract"
)

type Initializer struct {
	providers []contract.ToolProvider
}

func NewAIItz() *Initializer {
	return &Initializer{}
}

func (i *Initializer) RegisterTool(provider contract.ToolProvider) {
	i.providers = append(i.providers, provider)
}

func (i *Initializer) GetAllTools() []contract.ToolDefinition {
	var tools []contract.ToolDefinition

	for _, p := range i.providers {
		tools = append(tools, p.GetAITools()...)
	}
	return tools
}
