package handler

import (
	"html/template"
	"net/http"
)

func (h Handler) SignInPage(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"web/admin/templates/sign_in.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Error string
	}{
		Error: "",
	}

	err = ts.ExecuteTemplate(w, "sign-in", data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

}
