package handler

import (
	"net/http"

	postP "github.com/GoEnterpricePlatform/goEP-core/pkg/posts/port"
	"github.com/GoEnterpricePlatform/goEP-core/web/ai-tool-calling/renderer"
	"github.com/GoEnterpricePlatform/goEP-core/web/shared/api/middlewares"
)

type Handler struct {
	PostSrv             postP.PostSrv
	ApiBaseUrl          string
	ToolCallingRenderer *renderer.Renderer
	MdwSrvTmpl          *middlewares.MdwSrvTmpl
}

func NewPostWebAiHandler(
	postSrv postP.PostSrv,
	apiBaseUrl string,
	toolCallingRenderer *renderer.Renderer,
	mdwSrvtmpl *middlewares.MdwSrvTmpl,
) *Handler {
	h := &Handler{
		ApiBaseUrl:          apiBaseUrl,
		PostSrv:             postSrv,
		ToolCallingRenderer: toolCallingRenderer,
		MdwSrvTmpl:          mdwSrvtmpl,
	}

	return h
}

func (h Handler) RegisterRoutes(muxV1 *http.ServeMux) {

	// Actions - form submissions
	muxV1.HandleFunc("POST /admin/ai/posts", h.Create)
	muxV1.HandleFunc("POST /admin/ai/posts/update", h.Update)
	muxV1.HandleFunc("POST /admin/ai/posts/delete", h.Delete)
	//muxV1.HandleFunc("GET /admin/ai/posts/edit", h.Edit)
}
