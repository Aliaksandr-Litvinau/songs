package main

import (
	"encoding/json"
	"log"
	"net/http"
	"songs/internal/app/infrastructure/musicapi"
)

func main() {
	http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		group := r.URL.Query().Get("group")
		song := r.URL.Query().Get("song")

		if group == "" || song == "" {
			http.Error(w, "Missing group or song parameter", http.StatusBadRequest)
			return
		}

		log.Printf("Received request: %s %s", r.Method, r.URL.String())

		response := musicapi.SongDetailResponse{
			ReleaseDate: "2006-07-16",
			Text:        "Test lyrics text!!!!!!!!!",
			Link:        "https://example.com/songs/1",
		}

		w.Header().Set("Content-Type", "application/json")
		data, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error marshaling response: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		log.Printf("Sending response: %s", string(data))
		_, err = w.Write(data)
		if err != nil {
			log.Printf("Error writing response: %v", err)
			return
		}
	})

	log.Printf("Starting mock server on :8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
