package service

import (
	"context"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/ai-tool-calling/domain"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/ai-tool-calling/repository/memory"
)

func (s *Service) GetAll(ctx context.Context) ([]*domain.ChatMessage, error) {
	messages := memory.Messages
	/* for _, msg := range messages {
		fmt.Printf("Role: %s\n", msg.Role)
		fmt.Printf("Content: %s\n", msg.Content)

		if msg.Table != nil {
			fmt.Printf("Table: %+v\n", *msg.Table)
		} else {
			fmt.Println("Table: nil")
		}

		fmt.Println("----")
	} */
	return messages, nil
}
