package contract

type ToolDefinition struct {
	Name        string
	Description string
	Schemma     map[string]interface{}
}

type ToolProvider interface {
	GetAITools() []ToolDefinition
}
