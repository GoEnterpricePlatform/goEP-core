package service

import (
	"github.com/GoEnterpricePlatform/goEP-core/pkg/ai-tool-calling/initializer"
	tcP "github.com/GoEnterpricePlatform/goEP-core/pkg/ai-tool-calling/port"
)

var _ tcP.ToolCallingSrv = &Service{}

type Service struct {
	ToolCallingAdt tcP.ToolCallingAdt
	ToolCallingItz *initializer.Initializer
	SystemPrompt   string
}

func NewToolCallingSrv(toolCallingAdt tcP.ToolCallingAdt, toolCallingItz *initializer.Initializer) *Service {
	return &Service{
		ToolCallingAdt: toolCallingAdt,
		ToolCallingItz: toolCallingItz,
	}
}
