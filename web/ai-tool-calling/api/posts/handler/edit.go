package handler

/* func (h *Handler) Edit(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	posts, err := h.PostSrv.GetAll(r.Context())
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := RespData{
		EditingID:  id,
		ActivePage: "post",
		ErrorMsg:   "",
		Posts:      posts,
	}

	h.ToolCallingRenderer.Templates.ExecuteTemplate(w, "base", data)
} */
