package handler

import (
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/posts/domain"
	"github.com/GoEnterpricePlatform/goEP-core/web/standard-web/ai-tool-calling/core"
	sharedC "github.com/GoEnterpricePlatform/goEP-core/web/shared/core"
)

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	desc := r.FormValue("desc")

	post := domain.Post{
		ID:    id,
		Title: r.FormValue("title"),
		Desc:  &desc,
	}

	err := h.PostSrv.Update(r.Context(), id, &post)
	if err != nil {
		h.ToolCallingRenderer.Render(w, "assistant-page", core.AssistantPageData{
			Error: "Error actualizando post: " + sharedC.UiErrorResp(err),
		})
		return
	}

	h.ToolCallingRenderer.Render(w, "assistant-page", core.AssistantPageData{
		/* Success: "Post actualizado correctamente", */
	})
}
