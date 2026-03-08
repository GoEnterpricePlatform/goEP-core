package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/ai-tool-calling/domain"
	"github.com/GoEnterpricePlatform/goEP-core/web/shared/ai/contract"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/responses"
)

// ! no se si deba poner el required por que se enviará al handler y el lo validará
// ! un usuario se puede leiminar mediante id o email
func (a *Adapter) SelectTool(ctx context.Context, prompt string) (*domain.ToolSelectionResult, error) {
	resp, err := a.Client.Responses.New(ctx, responses.ResponseNewParams{
		Model: openai.ChatModelGPT5Mini,
		Input: responses.ResponseNewParamsInputUnion{
			OfInputItemList: []responses.ResponseInputItemUnionParam{
				{
					OfInputMessage: &responses.ResponseInputItemMessageParam{
						Role: "system",
						Content: []responses.ResponseInputContentUnionParam{
							{
								OfInputText: &responses.ResponseInputTextParam{
									Text: "You are a command parser. Only extract the values explicitly mentioned by the user. Do not invent any values.",
								},
							},
						},
					},
				},
				{
					OfInputMessage: &responses.ResponseInputItemMessageParam{
						Role: "user",
						Content: []responses.ResponseInputContentUnionParam{
							{
								OfInputText: &responses.ResponseInputTextParam{
									Text: prompt,
								},
							},
						},
					},
				},
			},
			
		},
		Tools: a.tools,
		ToolChoice: responses.ResponseNewParamsToolChoiceUnion{
			OfToolChoiceMode: param.Opt[responses.ToolChoiceOptions]{
				Value: responses.ToolChoiceOptionsRequired,
			},
		},
	})
	fmt.Println("-----------------------------------------------b")

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	fmt.Println("-----------------------------------------------c")

	for _, item := range resp.Output {
		fmt.Printf("TYPE: %s\n", item.Type)
		/* fmt.Printf("RAW ITEM: %+v\n", item) */
		fmt.Println("-----------------------------------------------d")

		if item.Type != "function_call" {
			fmt.Println("-----------------------------------------------f")

			continue
		}

		toolCall := item.AsFunctionCall()

		var args map[string]any
		if err := json.Unmarshal([]byte(toolCall.Arguments), &args); err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		fmt.Println(args)

		// Extraemos la definición de la herramienta
		var toolDef contract.ToolDefinition
		for _, t := range a.tools {
			if t.OfFunction != nil && t.OfFunction.Name == toolCall.Name {
				toolDef = contract.ToolDefinition{
					Name:        t.OfFunction.Name,
					Description: t.OfFunction.Description.Value,
					Schemma:     t.OfFunction.Parameters,
				}
				break
			}
		}
		fmt.Println("**************************************")

		properties, ok := toolDef.Schemma["properties"].(map[string]interface{})
		requiredList, _ := toolDef.Schemma["required"].([]string)
		requiredMap := make(map[string]bool)
		for _, r := range requiredList {
			requiredMap[r] = true
		}

		fmt.Println(properties)
		fmt.Println(requiredList)
		fmt.Println(ok)
		fmt.Println(requiredMap)
		fmt.Println("**************************************2")

		fields := make([]domain.FieldInfo, 0)
		for name := range properties {
			// Extraemos del prompt si no hay valor
			// Aca podria ser dejarlo vacio
			//fmt.Println(val)
			val := getString(args, name)

			val = extractFieldFromText(prompt, name, val)

			//fmt.Println(val)
			fields = append(fields, domain.FieldInfo{
				Name:     name,
				Value:    val,
				Required: requiredMap[name],
			})
			fmt.Println(fields)
		}

		// Validamos campos faltantes
		missing := []string{}
		for _, f := range fields {
			if f.Required && f.Value == "" {
				missing = append(missing, f.Name)
			}
		}

		fmt.Println("**************************************3")

		//missing := validateMissing(args, required)
		fmt.Println("-----------------------------------------------g")
		fmt.Println(missing)
		fmt.Println(toolCall.Name)
		resp := &domain.ToolSelectionResult{
			Operation:    toolCall.Name,
			Arguments:    args,
			Missing:      missing,
			NeedsConfirm: true,
		}
		// Reconstruir Arguments si quieres mantenerlo
		resp.Arguments = make(map[string]any)
		for _, f := range fields {
			resp.Arguments[f.Name] = f.Value
		}
		fmt.Printf("Resp: %+v", resp)
		return resp, nil
	}

	return nil, nil
}

// Extrae valor del prompt si no está en args
func extractFieldFromText(prompt, keyword, current string) string {
	if current != "" {
		return current
	}
	re := regexp.MustCompile(`(?i)` + keyword + `\s*["']([^"']+)["']`)
	match := re.FindStringSubmatch(prompt)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

func validateMissing(args map[string]any, required []string) []string {
	var missing []string

	for _, field := range required {
		v, ok := args[field]
		if !ok || v == nil || v == "" {
			missing = append(missing, field)
		}
	}

	return missing
}

func getString(m map[string]any, key string) string {
	v, ok := m[key]
	if !ok {
		return ""
	}
	return fmt.Sprintf("%v", v)
}
