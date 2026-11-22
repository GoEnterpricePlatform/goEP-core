package handler

import (
	"net/http"
	"strings"

	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/ui/navigation"
	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/ui/pages"
	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/variations/port"
	"github.com/GoEnterpricePlatform/goEP-core/addons/web-app/src/ui/layouts"
	"github.com/GoEnterpricePlatform/goEP-core/internal/shared/front/api/middlewares"
	"github.com/a-h/templ"
)

type Handler struct {
	VariationSrv port.VariationSrv
	MdwSrvTmpl   *middlewares.MdwSrvTmpl
}

func NewHandler(variationSrv port.VariationSrv, mdwSrvTmpl *middlewares.MdwSrvTmpl) *Handler {
	return &Handler{VariationSrv: variationSrv, MdwSrvTmpl: mdwSrvTmpl}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	if mux == nil {
		return
	}
	view := h.MdwSrvTmpl.Authenticate(h.MdwSrvTmpl.RequirePermission("view.catalog", http.HandlerFunc(h.render)))
	mux.Handle("GET /goep-admin/modules/catalog", view)
	mux.HandleFunc("GET /catalog", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/v1/catalog/", http.StatusPermanentRedirect)
	})
	mux.Handle("GET /catalog/", view)
	protectAction := func(handler http.HandlerFunc) http.Handler {
		return h.MdwSrvTmpl.Authenticate(h.MdwSrvTmpl.RequirePermission("view.catalog", handler))
	}
	mux.Handle("POST /catalog/variations", protectAction(h.CreateVariation))
	mux.Handle("PUT /catalog/variations/{id}", protectAction(h.UpdateVariation))
	mux.Handle("DELETE /catalog/variations/{id}", protectAction(h.DeleteVariation))
	mux.Handle("POST /catalog/variations/{variationId}/options", protectAction(h.CreateOption))
	mux.Handle("PUT /catalog/variations/{variationId}/options/{id}", protectAction(h.UpdateOption))
	mux.Handle("DELETE /catalog/variations/{variationId}/options/{id}", protectAction(h.DeleteOption))
}

func (h *Handler) render(w http.ResponseWriter, r *http.Request) {
	activePage := catalogPage(r.URL.Path)
	var content templ.Component
	if activePage == "variations" {
		variations, err := h.VariationSrv.GetAllVariationsWithOptions(r.Context())
		if err != nil {
			content = pages.CatalogVariationsContent(nil, err.Error())
		} else {
			content = pages.CatalogVariationsContent(variations, "")
		}
	} else if activePage == "paddle-plans" || activePage == "polar-plans" {
		content = pages.CatalogPlaceholderContent(pages.SectionTitle(activePage))
	} else {
		content = pages.CatalogGeneralContent()
	}

	page := layouts.DashboardLayout(activePage, []string{"ROLE_SYSTEM_ADMIN"}, nil, navigation.AdminMenu, content)
	if err := page.Render(r.Context(), w); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func catalogPage(path string) string {
	path = strings.Trim(strings.TrimPrefix(path, "/v1/"), "/")
	if path == "goep-admin/modules/catalog" {
		return "general"
	}
	path = strings.TrimPrefix(path, "catalog")
	path = strings.Trim(path, "/")
	if path == "" {
		return "general"
	}
	return path
}
