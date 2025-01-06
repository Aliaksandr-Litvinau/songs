package server

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type ErrorResponse struct {
	Slug       string `json:"slug"`
	Error      string `json:"error,omitempty"`
	httpStatus int
}

func (e ErrorResponse) Render(w http.ResponseWriter) error {
	w.WriteHeader(e.httpStatus)
	return nil
}

func BadRequest(slug string, err error, w http.ResponseWriter) {
	httpRespondWithError(err, slug, w, "Bad Request", http.StatusBadRequest)
}

func NotFound(slug string, err error, w http.ResponseWriter) {
	httpRespondWithError(err, slug, w, "Not Found", http.StatusNotFound)
}

func InternalError(slug string, err error, w http.ResponseWriter) {
	httpRespondWithError(err, slug, w, "Internal Server Error", http.StatusInternalServerError)
}

func RespondWithError(err error, w http.ResponseWriter) {
	log.Printf("Error: %v", err)
	InternalError("internal-server-error", err, w)
}

func httpRespondWithError(err error, slug string, w http.ResponseWriter, msg string, status int) {
	log.Printf("error: %s, slug: %s, msg: %s", err, slug, msg)

	resp := ErrorResponse{
		Slug:       slug,
		httpStatus: status,
	}

	if os.Getenv("DEBUG_ERRORS") != "" && err != nil {
		resp.Error = err.Error()
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(resp)
}
