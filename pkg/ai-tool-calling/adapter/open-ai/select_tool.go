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
func (a *Adapter) SelectTool(ctx context.Context, messages []*domain.ChatMessage) (*domain.ToolSelectionResult, error) {
	// Elusuario debe recibir feeback del recurso a crear eliinar
	// Cunado el usuario use palabras como crealo actualizalo donde
	// require una modificacion del sistema entonces usamos create_post_execute
	// aqui ya no quiero que modifiques los campos por ningun motivo sino llamar al
	// operacion adecuada

	// El detalle es no se si estoy

	resp, err := a.Client.Responses.New(ctx, responses.ResponseNewParams{
		Model: openai.ChatModelGPT5Mini,
		Input: responses.ResponseNewParamsInputUnion{
			OfInputItemList: buildConversation(messages),
		},
		Tools: a.tools,
		ToolChoice: responses.ResponseNewParamsToolChoiceUnion{
			OfToolChoiceMode: param.Opt[responses.ToolChoiceOptions]{
				Value: responses.ToolChoiceOptionsAuto,
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

		// 👇 EXTRAER ACTION
		action := getString(args, "action")

		// 👇 opcional: quitar action de los argumentos
		delete(args, "action")

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
		requiredRaw, _ := toolDef.Schemma["required"].([]string)
		fmt.Println(properties)
		fmt.Println(requiredRaw)
		requiredMap := make(map[string]bool)
		for _, r := range requiredRaw {
			requiredMap[r] = true
		}

		fmt.Println(properties)
		fmt.Println(requiredRaw)
		fmt.Println(ok)
		fmt.Println(requiredMap)
		fmt.Println("**************************************2")

		fields := make([]domain.FieldInfo, 0)
		for name := range properties {
			// Extraemos del prompt si no hay valor
			// Aca podria ser dejarlo vacio
			//fmt.Println(val)
			val := getString(args, name)

			var prompt string

			for i := len(messages) - 1; i >= 0; i-- {
				if messages[i].Role == "user" {
					prompt = messages[i].Content
					break
				}
			}

			val = extractFieldFromText(prompt, name, val)

			//fmt.Println(val)
			fields = append(fields, domain.FieldInfo{
				Name:     name,
				Value:    val,
				Required: requiredMap[name],
			})
		}
		fmt.Println("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
		fmt.Println(fields)
		fmt.Println("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")

		fmt.Println("**************************************3")

		fmt.Println("-----------------------------------------------g")
		fmt.Println(toolCall.Name)
		resp := &domain.ToolSelectionResult{
			Operation:    toolCall.Name,
			Action:       action,
			Arguments:    args,
			Required:     requiredRaw,
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

func getString(m map[string]any, key string) string {
	v, ok := m[key]
	if !ok {
		return ""
	}
	return fmt.Sprintf("%v", v)
}

func buildConversation(messages []*domain.ChatMessage) []responses.ResponseInputItemUnionParam {

	items := []responses.ResponseInputItemUnionParam{
		{
			OfInputMessage: &responses.ResponseInputItemMessageParam{
				Role: "system",
				Content: []responses.ResponseInputContentUnionParam{
					{
						OfInputText: &responses.ResponseInputTextParam{
							Text: `
							You are an admin assistant that manages system resources using tools.

							Rules:
								- If the user wants to delete or update a resource but no ID is provided, call the tool that lists the resources first.
								- Do not repeat the same listing tool if it was already executed in the previous step unless the user explicitly asks again.
								- Use the conversation context to decide the next step.

							Action rules: 
								- When collecting or editing fields, set action="prepare".
								- When the user confirms the operation (for example: yes, create it, go ahead, execute), set action="execute".
								- Do not execute operations unless the user clearly confirms.
								
							Never repeat a previous tool unless the user explicitly asks for it.
							`,
						},
					},
				},
			},
		},
	}

	for _, m := range messages {

		if m.Content == "" {
			continue
		}

		items = append(items,
			responses.ResponseInputItemUnionParam{
				OfInputMessage: &responses.ResponseInputItemMessageParam{
					Role: m.Role,
					Content: []responses.ResponseInputContentUnionParam{
						{
							OfInputText: &responses.ResponseInputTextParam{
								Text: m.Content,
							},
						},
					},
				},
			},
		)
	}

	return items
}
