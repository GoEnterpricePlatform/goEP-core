package openai

import (
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type OpenaiClient struct {
	Client openai.Client
}

func NewOpenAIClient(apiKey string) *OpenaiClient {
	return &OpenaiClient{
		Client: openai.NewClient(
			option.WithAPIKey(apiKey),
		),
	}
}
