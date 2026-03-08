package openai

import (
	"github.com/GoEnterpricePlatform/goEP-core/pkg/ai-tool-calling/initializer"
	ctP "github.com/GoEnterpricePlatform/goEP-core/pkg/ai-tool-calling/port"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

var _ ctP.ToolCallingAdt = &Adapter{}

type Adapter struct {
	Client openai.Client
	tools  []responses.ToolUnionParam
}

func NewToolCallingAdt(client openai.Client) *Adapter {
	return &Adapter{
		Client: client,
	}
}

func (a *Adapter) SetTools(initializer *initializer.Initializer) {
	defs := initializer.GetAllTools()

	tools := mapToOpenAITools(defs)

	a.tools = tools
}
