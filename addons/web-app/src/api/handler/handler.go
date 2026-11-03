package handler

import (
	"net/http"
	"path"
	"strings"

	"github.com/GoEnterpricePlatform/goEP-core/addons/posts/port"
	"github.com/GoEnterpricePlatform/goEP-core/addons/web-app/resources"
	"github.com/GoEnterpricePlatform/goEP-core/addons/web-app/src/ui/pages"
	moduleRegistry "github.com/GoEnterpricePlatform/goEP-core/internal/modules"
	"github.com/GoEnterpricePlatform/goEP-core/internal/shared/front/api/middlewares"
)

type Handler struct {
	PostSrv          port.PostSrv
	MdwSrvTmpl       *middlewares.MdwSrvTmpl
	AvailableModules []moduleRegistry.Manifest
}

func NewSrcHandler(mux *http.ServeMux, templateV1 *http.ServeMux, postSrv port.PostSrv, mdwSrvtmpl *middlewares.MdwSrvTmpl, availableModules []moduleRegistry.Manifest) *Handler {
	h := &Handler{
		PostSrv:          postSrv,
		MdwSrvTmpl:       mdwSrvtmpl,
		AvailableModules: availableModules,
	}

	// Redirects requests from "/" to the landing page
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/landing", http.StatusFound)
	})

	mux.HandleFunc("/landing", h.LandingPage)
	mux.HandleFunc("/about_us", h.AboutUsPage)
	mux.HandleFunc("/blog", h.BlogPage)
	mux.HandleFunc("/contact", h.ContactPage)
	templateV1.Handle("GET /goep-admin/modules", h.MdwSrvTmpl.Authenticate(http.HandlerFunc(h.ModulesPage)))
	templateV1.Handle("GET /goep-admin/modules/{module}", h.MdwSrvTmpl.Authenticate(http.HandlerFunc(h.ModulePage)))

	templateV1.HandleFunc("/goep-admin/access-denied", h.AccessDeniedPage)

	return h
}

func (h *Handler) ModulesPage(w http.ResponseWriter, r *http.Request) {
	cards := make([]pages.ModuleCard, 0, len(h.AvailableModules))
	for _, manifest := range h.AvailableModules {
		iconURL := ""
		if manifest.Icon != "" {
			iconPath := path.Join("modules", manifest.Name, strings.TrimPrefix(manifest.Icon, "static/"))
			iconURL = resources.StaticPath(iconPath)
		}
		cards = append(cards, pages.ModuleCard{
			Name:    manifest.Name,
			Version: manifest.Version,
			Summary: manifest.Summary,
			IconURL: iconURL,
			Href:    moduleAdminHref(manifest.Name),
		})
	}

	if err := pages.ModulesPage(cards).Render(r.Context(), w); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func moduleAdminHref(name string) string {
	if name == "identity" {
		return "/v1/goep-admin/general"
	}
	return "/v1/goep-admin/modules/" + name
}

func (h *Handler) ModulePage(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("module")
	for _, manifest := range h.AvailableModules {
		if manifest.Name != name {
			continue
		}
		if err := pages.ModulePage(manifest.Name, manifest.Summary).Render(r.Context(), w); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}
	http.NotFound(w, r)
}
