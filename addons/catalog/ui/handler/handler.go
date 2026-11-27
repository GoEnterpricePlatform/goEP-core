package handler

import (
	"net/http"
	"strings"

	paddlePort "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/port"
	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/ui/navigation"
	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/ui/pages"
	variationPort "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/variations/port"
	"github.com/GoEnterpricePlatform/goEP-core/addons/web-app/src/ui/layouts"
	"github.com/GoEnterpricePlatform/goEP-core/internal/shared/front/api/middlewares"
	"github.com/a-h/templ"
)

type Handler struct {
	VariationSrv  variationPort.VariationSrv
	PaddlePlanSrv paddlePort.PaddlePlanSrv
	MdwSrvTmpl    *middlewares.MdwSrvTmpl
}

func NewHandler(variationSrv variationPort.VariationSrv, paddlePlanSrv paddlePort.PaddlePlanSrv, mdwSrvTmpl *middlewares.MdwSrvTmpl) *Handler {
	return &Handler{VariationSrv: variationSrv, PaddlePlanSrv: paddlePlanSrv, MdwSrvTmpl: mdwSrvTmpl}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	if mux == nil {
		return
	}
	view := h.MdwSrvTmpl.Authenticate(h.MdwSrvTmpl.RequirePermission("view.catalog", http.HandlerFunc(h.render)))
	mux.Handle("GET /goep-admin/modules/catalog", view)
	mux.Handle("GET /catalog/paddle-plans/new", view)
	mux.Handle("GET /catalog/paddle-plans/{id}/edit", view)
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
	mux.Handle("POST /catalog/paddle-plans", protectAction(h.CreatePaddlePlan))
	mux.Handle("POST /catalog/paddle-plans/items", protectAction(h.GeneratePaddlePlanItems))
	mux.Handle("POST /catalog/paddle-plans/items/{itemIndex}/features", protectAction(h.AddPaddlePlanFeature))
	mux.Handle("DELETE /catalog/paddle-plans/items/{itemIndex}/features/{featureIndex}", protectAction(h.RemovePaddlePlanFeature))
	mux.Handle("PUT /catalog/paddle-plans/{id}", protectAction(h.UpdatePaddlePlan))
	mux.Handle("DELETE /catalog/paddle-plans/{id}", protectAction(h.DeletePaddlePlan))
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
	} else if activePage == "paddle-plans/new" {
		variations, err := h.VariationSrv.GetAllVariationsWithOptions(r.Context())
		loadError := ""
		if err != nil {
			loadError = err.Error()
		}
		content = pages.CatalogPaddlePlanFormPage(nil, variations, loadError)
	} else if strings.HasPrefix(activePage, "paddle-plans/") && strings.HasSuffix(activePage, "/edit") {
		id := strings.TrimSuffix(strings.TrimPrefix(activePage, "paddle-plans/"), "/edit")
		plan, planErr := h.PaddlePlanSrv.Get(r.Context(), id)
		if planErr != nil {
			content = pages.CatalogPaddlePlansPage(nil, planErr.Error())
		} else {
			variations, variationErr := h.VariationSrv.GetAllVariationsWithOptions(r.Context())
			loadError := ""
			if variationErr != nil {
				loadError = variationErr.Error()
			}
			content = pages.CatalogPaddlePlanFormPage(plan, variations, loadError)
		}
	} else if activePage == "paddle-plans" {
		plans, _, _, err := h.PaddlePlanSrv.GetAll(r.Context(), 100, 1)
		if err != nil {
			content = pages.CatalogPaddlePlansPage(nil, err.Error())
		} else {
			content = pages.CatalogPaddlePlansPage(plans, "")
		}
	} else if activePage == "polar-plans" {
		content = pages.CatalogPlaceholderContent(pages.SectionTitle(activePage))
	} else {
		content = pages.CatalogGeneralContent()
	}

	layoutPage := activePage
	if strings.HasPrefix(activePage, "paddle-plans/") {
		layoutPage = "paddle-plans"
	}
	page := layouts.DashboardLayout(layoutPage, []string{"ROLE_SYSTEM_ADMIN"}, nil, navigation.AdminMenu, content)
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
