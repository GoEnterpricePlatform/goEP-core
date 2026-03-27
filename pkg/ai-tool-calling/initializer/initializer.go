package initializer

import (
	"fmt"
	"os"

	"github.com/GoEnterpricePlatform/goEP-core/web/shared/ai/contract"
)

type Initializer struct {
	providers []contract.ToolProvider
	tools     map[string]contract.ToolDefinition
}

func NewAIItz() *Initializer {
	return &Initializer{
		tools: map[string]contract.ToolDefinition{},
	}
}

func (i *Initializer) RegisterTool(provider contract.ToolProvider) {
	i.providers = append(i.providers, provider)
	for _, t := range provider.GetAITools() {
		i.tools[t.Name] = t
	}
}

func (i *Initializer) GetTool(name string) (contract.ToolDefinition, bool) {
	t, ok := i.tools[name]
	return t, ok
}

func (i *Initializer) GetAllTools() []contract.ToolDefinition {
	var tools []contract.ToolDefinition

	for _, t := range i.tools {
		tools = append(tools, t)
	}
	return tools
}

func (i *Initializer) GetSystemPrompt() (string, error) {
	name := "pkg/ai-tool-calling/initializer/system_prompt/admin_tool_selection/admin_tool_selection.md"
	body, err := os.ReadFile(name)
	if err != nil {
		return "", fmt.Errorf("Failed to read admin_tool_selection.md: %w", err)
	}
	return string(body), nil
}
