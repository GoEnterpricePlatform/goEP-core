package openai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/ai-tool-calling/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
	"github.com/GoEnterpricePlatform/goEP-core/web/shared/ai/contract"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/responses"
)

// "action" and "delete" are internal fields; we remove them so they are not displayed in the table.
func (a *Adapter) SelectTool(ctx context.Context, messages []*domain.ChatMessage, systemPrompt string) (*domain.ToolSelectionResult, error) {
	resp, err := a.Client.Responses.New(ctx, responses.ResponseNewParams{
		Model: openai.ChatModelGPT5Mini,
		Input: responses.ResponseNewParamsInputUnion{
			OfInputItemList: buildConversation(messages, systemPrompt),
		},
		Tools: a.tools,
		ToolChoice: responses.ResponseNewParamsToolChoiceUnion{
			OfToolChoiceMode: param.Opt[responses.ToolChoiceOptions]{
				Value: responses.ToolChoiceOptionsRequired,
			},
		},
	})

	if err != nil {
		var apiErr *openai.Error
		if errors.As(err, &apiErr) {
			if apiErr.Code == "invalid_api_key" {
				return nil, fmt.Errorf("%w: invalid openai api key: %w", sharedD.ErrInvalidApiKey, err)
			}
		}
		return nil, fmt.Errorf("openai request failed: %w", err)
	}

	for _, item := range resp.Output {

		// fmt.Printf("TYPE: %s\n", item.Type)

		if item.Type != "function_call" {
			continue
		}

		toolCall := item.AsFunctionCall()

		var args map[string]any
		if err := json.Unmarshal([]byte(toolCall.Arguments), &args); err != nil {
			return nil, err
		}

		// fmt.Println(args)

		// EXTRACT ACTION
		action := getString(args, "action")

		// remove action from arguments
		delete(args, "action")

		// We extract the definition of the tool
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

		properties, _ := toolDef.Schemma["properties"].(map[string]interface{})
		requiredRaw, _ := toolDef.Schemma["required"].([]string)
		requiredMap := make(map[string]bool)
		for _, r := range requiredRaw {
			requiredMap[r] = true
		}

		// fmt.Println(properties)
		//fmt.Println(requiredRaw)
		//fmt.Println(ok)
		//fmt.Println(requiredMap)

		fields := make([]domain.FieldInfo, 0)
		for name := range properties {
			// We extract the value from the prompt if there is no value
			val := getString(args, name)

			var prompt string

			for i := len(messages) - 1; i >= 0; i-- {
				if messages[i].Role == "user" {
					prompt = messages[i].Content
					break
				}
			}

			val = extractFieldFromText(prompt, name, val)

			fields = append(fields, domain.FieldInfo{
				Name:     name,
				Value:    val,
				Required: requiredMap[name],
			})
		}

		// fmt.Println(toolCall.Name)

		message := getString(args, "message")
		delete(args, "message")

		resp := &domain.ToolSelectionResult{
			Operation:    toolCall.Name,
			Action:       action,
			Arguments:    args,
			Message:      message,
			Required:     requiredRaw,
			NeedsConfirm: true,
		}
		return resp, nil
	}

	return nil, nil
}

// Extract value from prompt if not in args
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

func buildConversation(messages []*domain.ChatMessage, systemPrompt string) []responses.ResponseInputItemUnionParam {

	items := []responses.ResponseInputItemUnionParam{
		{
			OfInputMessage: &responses.ResponseInputItemMessageParam{
				Role: "system",
				Content: []responses.ResponseInputContentUnionParam{
					{
						OfInputText: &responses.ResponseInputTextParam{
							Text: systemPrompt,
						},
					},
				},
			},
		},
	}

	// fmt.Println("------------------------- Build Conversation")
	for _, m := range messages {

		if m.Content == "" {
			continue
		}

		// fmt.Printf("Role: %s\n", m.Role)
		//fmt.Printf("Content: %s\n", m.Content)

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
	// fmt.Println("------------------------- Build Conversation end")

	return items
}
