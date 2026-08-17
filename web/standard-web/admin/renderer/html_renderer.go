package renderer

import (
	"html/template"
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/web/shared/templates"
)

type Renderer struct {
	Templates *template.Template
}

func NewAdminRenderer() *Renderer {
	t := template.New("admin")

	template.Must(t.ParseGlob("web/standard-web/admin/templates/*.html"))
	template.Must(t.ParseGlob("web/standard-web/admin/templates/components/*.html"))
	template.Must(t.ParseGlob("web/standard-web/admin/templates/layout/*.html"))
	template.Must(t.ParseGlob("web/standard-web/admin/templates/posts/*.html"))
	template.Must(t.ParseGlob("web/shared/templates/*.html"))

	return &Renderer{
		Templates: t,
	}
}

func (r *Renderer) Render(w http.ResponseWriter, name string, data any) {
	err := r.Templates.ExecuteTemplate(w, name, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		r.Templates.ExecuteTemplate(w, "error", templates.ErrorData{
			ErrorMsg: "Internal Server Error",
		})
	}
}
