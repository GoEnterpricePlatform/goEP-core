package contract

import "context"

type ToolCategory string

const (
	ToolCategoryRead  ToolCategory = "read"
	ToolCategoryWrite ToolCategory = "write"
)

type ToolHandler func(ctx context.Context, args map[string]any) (map[string]any, error)

type ToolDefinition struct {
	Name        string
	Description string
	Schemma     map[string]interface{}
	Category    ToolCategory
	Handler     ToolHandler
}

type ResourceField struct {
	Name        string
	Description string
	Type        string
}

type ToolProvider interface {
	GetAITools() []ToolDefinition
}
