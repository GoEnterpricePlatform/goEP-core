package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/ai-tool-calling/domain"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/ai-tool-calling/repository/memory"
)

func (s *Service) SelectTool(ctx context.Context, prompt string) error {
	fmt.Println("############################################")
	memory.Messages = append(memory.Messages, &domain.ChatMessage{
		Role:    "user",
		Content: prompt,
	})
	messages := memory.Messages
	result, err := s.ToolCallingAdt.SelectTool(ctx, messages)
	if err != nil {
		return err
	}
	fmt.Println("############################################2")

	if result == nil {
		memory.Messages = append(messages, &domain.ChatMessage{
			Role:    "system",
			Content: "no hay result",
		})
		return nil
	}
	fmt.Println("############################################3")

	if result.Action == "prepare" {
		fmt.Println("############################################4")

		table := ArgumentsToTable(
			result.Operation,
			result.Arguments,
			result.Required,
		)

		msgContent := fmt.Sprintf(
			"Preparing %s operation. Please review the fields.",
			result.Operation,
		)

		memory.Messages = append(messages, &domain.ChatMessage{
			Role:    "system",
			Content: msgContent,
			Table:   table,
		})

		return nil
	}
	fmt.Println("############################################5")

	tool, ok := s.ToolCallingItz.GetTool(result.Operation)
	if !ok {
		return errors.New("tool not found")
	}
	fmt.Println("############################################6")

	data, err := tool.Handler(ctx, result.Arguments)
	if err != nil {
		return err
	}

	table := MapToTable(data)
	fmt.Println("############################################7")

	msgContent := fmt.Sprintf(
		"Tool %s executed successfully. The following table shows the result.",
		result.Operation,
	)

	memory.Messages = append(messages, &domain.ChatMessage{
		Role:    "system",
		Content: msgContent,
		Table:   table,
	})

	for _, msg := range messages {
		fmt.Printf("Role: %s\n", msg.Role)
		fmt.Printf("Content: %s\n", msg.Content)

		if msg.Table != nil {
			fmt.Printf("Table: %+v\n", *msg.Table)
		} else {
			fmt.Println("Table: nil")
		}

		fmt.Println("----")
	}

	return nil
}

func MapToTable(data map[string]any) *domain.TableView {
	cols, _ := data["columns"].([]string)
	rows, _ := data["rows"].([]map[string]any)

	return &domain.TableView{
		Columns: cols,
		Rows:    rows,
	}
}

func ArgumentsToTable(
	title string,
	args map[string]any,
	required []string,
) *domain.TableView {

	requiredMap := map[string]bool{}

	for _, r := range required {
		requiredMap[r] = true
	}

	rows := []map[string]any{}

	for k, v := range args {

		rows = append(rows, map[string]any{
			"field":    k,
			"value":    v,
			"required": requiredMap[k],
		})
	}

	return &domain.TableView{
		//Title:   title,
		Columns: []string{"field", "value", "required"},
		Rows:    rows,
	}
}
