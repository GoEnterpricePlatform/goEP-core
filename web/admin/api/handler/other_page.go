package handler

import (
	"html/template"
	"net/http"
)

func (h Handler)OtherPage(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"web/admin/templates/layout/base.html",
		"web/admin/templates/components/sidebar.html",
		"web/admin/templates/other.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		ActivePage string
	}{
		ActivePage: "other",
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
