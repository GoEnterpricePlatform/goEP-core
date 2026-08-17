package renderer

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/web/shared/templates"
)

type Renderer struct {
	Templates *template.Template
}

func NewToolCallingRenderer() *Renderer {
	t := template.New("tool-calling")

	//template.Must(t.ParseGlob("web/chat-tool-calling/templates/layout/*.html"))
	//template.Must(t.ParseGlob("web/chat-tool-calling/templates/*.html"))
	//template.Must(t.ParseGlob("web/chat-tool-calling/templates/components/*.html"))

	return &Renderer{
		Templates: t,
	}
}

func (r *Renderer) Render(w http.ResponseWriter, name string, data any) {
	files := []string{
		"web/standard-web/ai-tool-calling/templates/layout/base.html",
		fmt.Sprintf("web/standard-web/ai-tool-calling/templates/%s.html", name),
	}

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		r.Templates.ExecuteTemplate(w, "error", templates.ErrorData{
			ErrorMsg: "Internal Server Error",
		})
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		// fmt.Println("TEMPLATE ERROR:", err) // for debug
		w.WriteHeader(http.StatusInternalServerError)
		r.Templates.ExecuteTemplate(w, "error", templates.ErrorData{
			ErrorMsg: "Internal Server Error",
		})
		return
	}
}
