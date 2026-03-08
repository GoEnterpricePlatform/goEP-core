package handler

import (
	"fmt"
	"net/http"
)

func (h *Handler) AssistantExecute(w http.ResponseWriter, r *http.Request) {
	operation := r.FormValue("operation")
	fmt.Println(operation)
}
