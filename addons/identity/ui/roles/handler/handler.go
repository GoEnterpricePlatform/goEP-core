package handler

import (
	"net/http"

	roleP "github.com/GoEnterpricePlatform/goEP-core/addons/identity/roles/port"
	"github.com/GoEnterpricePlatform/goEP-core/web/shared/api/middlewares"
)

type Handler struct {
	RoleSrv    roleP.RoleSrv
	ApiBaseUrl string
	MdwSrvTmpl *middlewares.MdwSrvTmpl
}

func NewRolesTmplHandler(
	roleSrv roleP.RoleSrv,
	apiBaseUrl string,
	mdwSrvtmpl *middlewares.MdwSrvTmpl,
) *Handler {
	h := &Handler{
		RoleSrv:    roleSrv,
		ApiBaseUrl: apiBaseUrl,
		MdwSrvTmpl: mdwSrvtmpl,
	}

	return h
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux, templateV1 *http.ServeMux) {

	// Protect routes with middlewares
	// Actions - form submissions
	templateV1.HandleFunc("POST /roles", h.Create)
	//templateV1.HandleFunc("GET /roles/{id}", h.Get)
	//templateV1.HandleFunc("GET /roles", h.GetAll)
	templateV1.HandleFunc("PUT /roles/{id}", h.Update)
	templateV1.HandleFunc("DELETE /roles/{id}", h.Delete)

}
