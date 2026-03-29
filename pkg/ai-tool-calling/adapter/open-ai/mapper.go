package openai

import (
	"github.com/GoEnterpricePlatform/goEP-core/web/shared/ai/contract"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/responses"
)

func mapToOpenAITools(defs []contract.ToolDefinition) []responses.ToolUnionParam {
	var tools []responses.ToolUnionParam

	for _, def := range defs {
		tools = append(tools, responses.ToolUnionParam{
			OfFunction: &responses.FunctionToolParam{
				Name:        def.Name,
				Description: param.Opt[string]{Value: def.Description},
				Parameters:  def.Schemma,
				
			},
		})
	}

	return tools
}
