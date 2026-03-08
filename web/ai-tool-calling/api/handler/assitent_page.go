package handler

import (
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/web/ai-tool-calling/core"
)

func (h Handler) AssistantPage(w http.ResponseWriter, r *http.Request) {
	h.ToolCallingRenderer.Render(w, "assistant_chat", core.AssistantPageData{
		Prompt: "",
	})
}
