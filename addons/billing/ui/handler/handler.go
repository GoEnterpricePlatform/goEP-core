package handler

import (
	"net/http"
	"strings"

	moduleNavigation "github.com/GoEnterpricePlatform/goEP-core/addons/billing/ui/navigation"
	"github.com/GoEnterpricePlatform/goEP-core/addons/billing/ui/pages"
	"github.com/GoEnterpricePlatform/goEP-core/addons/web-app/src/ui/layouts"
	"github.com/GoEnterpricePlatform/goEP-core/internal/shared/front/api/middlewares"
)

type Handler struct {
	MdwSrvTmpl *middlewares.MdwSrvTmpl
}

func NewHandler(mdwSrvTmpl *middlewares.MdwSrvTmpl) *Handler {
	return &Handler{MdwSrvTmpl: mdwSrvTmpl}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	if mux == nil {
		return
	}
	mux.Handle("GET /goep-admin/modules/billing", h.MdwSrvTmpl.Authenticate(h.MdwSrvTmpl.RequirePermission("view.billing", http.HandlerFunc(h.render))))
	mux.HandleFunc("GET /billing", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/v1/billing/", http.StatusPermanentRedirect)
	})
	mux.Handle("GET /billing/", h.MdwSrvTmpl.Authenticate(h.MdwSrvTmpl.RequirePermission("view.billing", http.HandlerFunc(h.render))))
}

func (h *Handler) render(w http.ResponseWriter, r *http.Request) {
	activePage, pageTitle := modulePage(r.URL.Path)
	if "dashboard" == "dashboard" {
		if "replace" == "topbar" {
			page := layouts.GoEPDashboardLayout(activePage, moduleNavigation.AdminMenu, pages.BillingDashboardContent(activePage))
			if err := page.Render(r.Context(), w); err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		}
		page := layouts.DashboardLayout(activePage, []string{"ROLE_SYSTEM_ADMIN"}, nil, moduleNavigation.AdminMenu, pages.BillingDashboardContent(activePage))
		if err := page.Render(r.Context(), w); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}
	if "replace" == "topbar" {
		page := layouts.GoEPModuleLayout(pages.BillingContent(pageTitle))
		if err := page.Render(r.Context(), w); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}
	if err := pages.BillingPage().Render(r.Context(), w); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func modulePage(path string) (string, string) {
	page := strings.Trim(path, "/")
	page = strings.TrimPrefix(page, "v1/")
	if page == "billing" {
		return "overview", "Billing"
	}
	prefix := "billing" + "/"
	if !strings.HasPrefix(page, prefix) {
		return "overview", "Billing"
	}
	page = strings.TrimPrefix(page, prefix)
	if page == "" {
		return "overview", "Billing"
	}
	title := strings.ReplaceAll(strings.ReplaceAll(page, "-", " "), "/", " / ")
	return page, strings.Title(title)
}
