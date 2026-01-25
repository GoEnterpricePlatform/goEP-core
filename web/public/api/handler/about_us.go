package handler

import (
	"html/template"
	"net/http"
)

func (h Handler) AboutUsPage(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"web/public/templates/layout/base.html",
		"web/public/templates/about_us.html",
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
		ActivePage: "about_us",
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}