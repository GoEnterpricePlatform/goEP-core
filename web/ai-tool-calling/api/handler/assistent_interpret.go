package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	tcCore "github.com/GoEnterpricePlatform/goEP-core/pkg/ai-tool-calling/core"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/posts/domain"
	"github.com/GoEnterpricePlatform/goEP-core/web/ai-tool-calling/core"
	sharedC "github.com/GoEnterpricePlatform/goEP-core/web/shared/core"
)

var messages []core.ChatMessage

func (h *Handler) AssistantInterpret(w http.ResponseWriter, r *http.Request) {
	fmt.Println("--------------------------------------1")
	req := &tcCore.PromptReq{
		Prompt: r.FormValue("prompt"),
	}
	fmt.Println(req)
	// ! isPrompValid

	fmt.Println("--------------------------------------2")

	result, err := h.toolCallingSrv.SelectTool(context.Background(), req.Prompt)
	if err != nil {
		fmt.Println("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
		fmt.Println(err)
		fmt.Println("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
		h.ToolCallingRenderer.Render(w, "assistant/page", core.AssistantPageData{
			Prompt: req.Prompt,
			Error:  sharedC.UiErrorResp(err),
		})
		return
	}
	fmt.Println("--------------------------------------3")

	if result == nil {
		h.ToolCallingRenderer.Render(w, "assistant/page", core.AssistantPageData{
			Prompt: req.Prompt,
			Error:  "No se detectó ninguna acción.",
		})
		return
	}
	fmt.Println("--------------------------------------4")

	messages = append(messages, core.ChatMessage{
		Role:    "user",
		Content: req.Prompt,
	})

	switch {

	case strings.HasPrefix(result.Operation, "create_"),
		strings.HasPrefix(result.Operation, "update_"),
		strings.HasPrefix(result.Operation, "delete_"):
		fmt.Println("--------------------------------------6")

		table := ArgumentsToTable(
			result.Operation,
			result.Arguments,
			result.Missing,
		)

		messages = append(messages, core.ChatMessage{
			Role:  "assistant",
			Table: table,
		})

		h.ToolCallingRenderer.Render(w, "assistant_chat", core.AssistantPageData{
			Prompt:    "",
			Messages:  messages,
			Arguments: result.Arguments,
			Missing:   result.Missing,
		})

	case result.Operation == "get_all_posts":
		fmt.Println("--------------------------------------121")

		posts, err := h.PostSrv.GetAll(r.Context())
		if err != nil {
			h.ToolCallingRenderer.Render(w, "assistant_chat", core.AssistantPageData{
				Prompt:   req.Prompt,
				Error:    "Error obteniendo posts.",
				Messages: messages, // mostramos el historial
			})
			return
		}
		fmt.Println("--------------------------------------123")

		table := PostsToTable(posts)

		// Lo agregamos al historial como mensaje del asistente
		messages = append(messages, core.ChatMessage{
			Role:  "assistant",
			Table: table,
		})

		// Renderizamos la página del chat con todos los mensajes
		h.ToolCallingRenderer.Render(w, "assistant_chat", core.AssistantPageData{
			Prompt:   "",
			Messages: messages,
		})

	default:
		fmt.Println("--------------------------------------10")

		// Lo agregamos al historial como mensaje del asistente
		messages = append(messages, core.ChatMessage{
			Role:    "assistant",
			Content: "operacion no soportada",
		})

		h.ToolCallingRenderer.Render(w, "assistant_chat", core.AssistantPageData{
			Prompt:   "",
			Messages: messages,
			Error:    "Operación no soportada.",
		})
	}
	fmt.Println("--------------------------------------11")
}

func PostsToTable(posts []*domain.Post) *core.TableView {

	rows := []map[string]any{}

	for _, p := range posts {
		rows = append(rows, map[string]any{
			"id":    p.ID,
			"title": p.Title,
			"desc":  p.Desc,
		})
	}

	return &core.TableView{
		Columns: []string{"id", "title", "desc"},
		Rows:    rows,
	}
}

func ArgumentsToTable(
	title string,
	args map[string]any,
	required []string,
) *core.TableView {

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

	return &core.TableView{
		//Title:   title,
		Columns: []string{"field", "value", "required"},
		Rows:    rows,
	}
}
