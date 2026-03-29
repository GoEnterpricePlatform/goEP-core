package core

import "github.com/GoEnterpricePlatform/goEP-core/pkg/ai-tool-calling/domain"

type AssistantPageData struct {
	Messages []*domain.ChatMessage
	Prompt   string
	Error    string
	//Arguments map[string]any
	/* 	ErrorSubTmpl string
	Success      string */
}
