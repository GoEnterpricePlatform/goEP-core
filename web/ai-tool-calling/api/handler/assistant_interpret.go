package handler

import (
	"context"
	"net/http"

	tcCore "github.com/GoEnterpricePlatform/goEP-core/pkg/ai-tool-calling/core"
	"github.com/GoEnterpricePlatform/goEP-core/web/ai-tool-calling/core"
	sharedC "github.com/GoEnterpricePlatform/goEP-core/web/shared/core"
)

func (h *Handler) AssistantInterpret(w http.ResponseWriter, r *http.Request) {
	req := &tcCore.PromptReq{
		Prompt: r.FormValue("prompt"),
	}

	if err := req.Validate(); err != nil {
		messages, getAllErr := h.toolCallingSrv.GetAll(r.Context())
		if getAllErr != nil {
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

	err := h.toolCallingSrv.SelectTool(context.Background(), req.Prompt)
	if err != nil {
		messages, getAllErr := h.toolCallingSrv.GetAll(r.Context())
		if getAllErr != nil {
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
	http.Redirect(w, r, "/v1/admin/assistant", http.StatusSeeOther)
}
