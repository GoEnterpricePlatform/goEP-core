package service

import "github.com/GoEnterpricePlatform/goEP-core/pkg/ai-tool-calling/port"

var _ port.ToolCallingSrv = &Service{}

type Service struct {
	ToolCallingAdt port.ToolCallingAdt
}

func NewToolCallingSrv(toolCallingAdt port.ToolCallingAdt) *Service {
	return &Service{
		ToolCallingAdt: toolCallingAdt,
	}
}
