package openai

import (
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func NewOpenAIClient(apiKey string) openai.Client {
	return openai.NewClient(
		option.WithAPIKey(apiKey),
	)
}
