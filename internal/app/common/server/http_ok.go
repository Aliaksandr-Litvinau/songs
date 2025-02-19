package server

import (
	"encoding/json"
	"net/http"
)

func RespondOK(data interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(data)
}
