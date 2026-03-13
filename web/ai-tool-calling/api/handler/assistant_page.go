package handler

import (
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/web/ai-tool-calling/core"
	//sharedC "github.com/GoEnterpricePlatform/goEP-core/web/shared/core"
)

func (h Handler) AssistantPage(w http.ResponseWriter, r *http.Request) {
	messages, err := h.toolCallingSrv.GetAll(r.Context())
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	h.ToolCallingRenderer.Render(w, "assistant_chat", core.AssistantPageData{
		Prompt:   "",
		Messages: messages,
	})
}
