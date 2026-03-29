package handler

import (
	"net/http"

	toolCallingP "github.com/GoEnterpricePlatform/goEP-core/pkg/ai-tool-calling/port"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/posts/port"
	"github.com/GoEnterpricePlatform/goEP-core/web/ai-tool-calling/renderer"
	"github.com/GoEnterpricePlatform/goEP-core/web/shared/api/middlewares"
)

type Handler struct {
	toolCallingSrv      toolCallingP.ToolCallingSrv
	ApiBaseUrl          string
	ToolCallingRenderer *renderer.Renderer
	MdwSrvTmpl          *middlewares.MdwSrvTmpl
	PostSrv             port.PostSrv
}

func NewChatToolCallingHandler(
	toolCallingSrv toolCallingP.ToolCallingSrv,
	apiBaseUrl string,
	toolCallingRenderer *renderer.Renderer,
	mdwSrvtmpl *middlewares.MdwSrvTmpl,
	postSrv port.PostSrv,
) *Handler {
	h := &Handler{
		toolCallingSrv:      toolCallingSrv,
		ApiBaseUrl:          apiBaseUrl,
		ToolCallingRenderer: toolCallingRenderer,
		MdwSrvTmpl:          mdwSrvtmpl,
		PostSrv:             postSrv,
	}

	return h
}

func (h Handler) RegisterRoutes(mux *http.ServeMux, muxV1 *http.ServeMux) {

	// Pages - render html
	muxV1.Handle("/admin/assistant", h.MdwSrvTmpl.Authenticate(http.HandlerFunc(h.AssistantPage)))

	// Actions - form submissions,
	// ! proteger con algo como el middleware que se usa para web
	muxV1.HandleFunc("POST /admin/assistant", h.AssistantInterpret)
	muxV1.HandleFunc("POST /admin/execute", h.AssistantInterpret)
}
