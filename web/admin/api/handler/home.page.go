package handler

import (
	"net/http"
	"text/template"
)

func (h Handler) HomePage(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"web/admin/api/web/templates/base.html",
		"web/admin/api/web/templates/home.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		ApiBaseUrl string
		ActivePage string
	}{
		ApiBaseUrl: h.ApiBaseUrl,
		ActivePage: "home",
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
