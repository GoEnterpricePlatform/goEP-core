package handler

import (
	"html/template"
	"net/http"

	"github.com/amorindev/go-cms-tmpl/pkg/posts/domain"
)

func (h Handler) LandingPage(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"web/public/templates/layout/base.html",
		"web/public/templates/landing.html",
		"web/public/templates/components/header.html",
		"web/public/templates/components/footer.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	posts, err := h.PostSrv.GetAll(r.Context())
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		ActivePage string
		Posts      []*domain.Post
	}{
		ActivePage: "landing",
		Posts:      posts,
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
