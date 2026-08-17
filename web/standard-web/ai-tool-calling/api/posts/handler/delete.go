package handler

import (
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/web/standard-web/ai-tool-calling/core"
)

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	if id == "" {
		h.ToolCallingRenderer.Render(w, "assistant-page", core.AssistantPageData{})
		return
	}

	err := h.PostSrv.Delete(r.Context(), id)
	if err != nil {
		h.ToolCallingRenderer.Render(w, "assistant-page", core.AssistantPageData{})
		return
	}

	h.ToolCallingRenderer.Render(w, "assistant-page", core.AssistantPageData{})
}
