package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/ai-tool-calling/domain"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/ai-tool-calling/repository/memory"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

func (s *Service) SelectTool(ctx context.Context, prompt string) error {
	memory.Messages = append(memory.Messages, &domain.ChatMessage{
		Role:    "user",
		Content: prompt,
	})
	messages := memory.Messages

	// We have attached the table to provide context to the llm
	chatMessageWithTable := BuildLLMMessages(memory.Messages)

	result, err := s.ToolCallingAdt.SelectTool(ctx, chatMessageWithTable, s.SystemPrompt)
	if err != nil {
		return sharedD.ManageError(err, "")
	}

	//fmt.Printf("Result from service: %+v\n", result)

	if result == nil {
		memory.Messages = append(messages, &domain.ChatMessage{
			Role:    "system",
			Content: "no hay result",
		})
		return nil
	}

	tool, ok := s.ToolCallingItz.GetTool(result.Operation)
	if !ok {
		return errors.New("tool not found")
	}

	// fmt.Printf("Arguments: %v\n", result.Arguments)

	if result.Action == "prepare" {
		table := ArgumentsToTable(
			result.Operation,
			result.Arguments,
			result.Required,
		)

		msgContent := fmt.Sprintf(
			"Preparing %s operation. %s",
			result.Operation, result.Message,
		)

		memory.Messages = append(messages, &domain.ChatMessage{
			Role:    "system",
			Content: msgContent,
			Table:   table,
		})

		return nil
	}

	data, err := tool.Handler(ctx, result.Arguments)
	if err != nil {
		return err
	}

	table := MapToTable(data)

	msgContent := fmt.Sprintf(
		"Tool %s executed successfully. %s ",
		result.Operation, result.Message,
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

func TableToString(table *domain.TableView) string {
	if table == nil {
		return ""
	}

	var sb strings.Builder

	// Header row
	sb.WriteString("| ")
	for _, col := range table.Columns {
		sb.WriteString(col + " | ")
	}
	sb.WriteString("\n")

	// Separator
	sb.WriteString("| ")
	for range table.Columns {
		sb.WriteString("--- | ")
	}
	sb.WriteString("\n")

	// Rows
	for _, row := range table.Rows {
		sb.WriteString("| ")
		for _, col := range table.Columns {
			val := fmt.Sprintf("%v", row[col])
			sb.WriteString(val + " | ")
		}
		sb.WriteString("\n")
	}

	return sb.String()
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

func BuildLLMMessages(messages []*domain.ChatMessage) []*domain.ChatMessage {
	var result []*domain.ChatMessage

	for _, m := range messages {
		content := m.Content

		// Here you add the table ONLY for the LLM
		if m.Table != nil {
			content += "\n\nTable:\n" + TableToString(m.Table)
		}
		result = append(result, &domain.ChatMessage{
			Role:    m.Role,
			Content: content,
		})
	}

	return result
}
