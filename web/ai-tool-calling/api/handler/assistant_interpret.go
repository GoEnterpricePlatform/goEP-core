package handler

import (
	"context"
	"fmt"
	"net/http"

	tcCore "github.com/GoEnterpricePlatform/goEP-core/pkg/ai-tool-calling/core"
	"github.com/GoEnterpricePlatform/goEP-core/web/ai-tool-calling/core"
	sharedC "github.com/GoEnterpricePlatform/goEP-core/web/shared/core"
)

func (h *Handler) AssistantInterpret(w http.ResponseWriter, r *http.Request) {
	fmt.Println("--------------------------------------1")
	req := &tcCore.PromptReq{
		Prompt: r.FormValue("prompt"),
	}
	fmt.Println(req)
	// ! isPrompValid

	fmt.Println("--------------------------------------2")

	err := h.toolCallingSrv.SelectTool(context.Background(), req.Prompt)
	if err != nil {
		fmt.Println("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
		fmt.Println(err)
		fmt.Println("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
		messages, err := h.toolCallingSrv.GetAll(r.Context())
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		h.ToolCallingRenderer.Render(w, "assistant_chat", core.AssistantPageData{
			Prompt:   req.Prompt,
			Error:    sharedC.UiErrorResp(err),
			Messages: messages,
		})
		return
	}
	fmt.Println("--------------------------------------4")
	http.Redirect(w, r, "/v1/admin/assistant", http.StatusSeeOther)
}
