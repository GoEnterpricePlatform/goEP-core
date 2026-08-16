package handler

import (
	"encoding/json"
	"net/http"
)

// It is used to know if the server is running
func Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := struct {
		Msg string `json:"msg"`
	}{
		Msg: "pong",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
