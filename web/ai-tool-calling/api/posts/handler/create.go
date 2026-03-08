package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/posts/domain"
	"github.com/GoEnterpricePlatform/goEP-core/web/ai-tool-calling/core"
	sharedC "github.com/GoEnterpricePlatform/goEP-core/web/shared/core"
)

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("-------------------------------01")
	desc := r.FormValue("desc")
	p := domain.Post{
		Title: r.FormValue("title"),
		Desc:  &desc,
	}
	fmt.Println("-------------------------------02")

	err := h.PostSrv.Create(context.Background(), &p)
	if err != nil {
		fmt.Println("-------------------------------03")

		h.ToolCallingRenderer.Render(w, "assistant_home", core.AssistantPageData{
			Error: "Error creando post: " + sharedC.UiErrorResp(err),
		})
		return
	}
	fmt.Println("-------------------------------04")

	h.ToolCallingRenderer.Render(w, "assistant_home", core.AssistantPageData{
	})
}
