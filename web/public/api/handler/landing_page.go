package handler

import (
	"html/template"
	"net/http"
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

	data := struct {
		ActivePage string
	}{
		ActivePage: "landing",
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}