package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) hitsReset(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)

	header := w.Header()
	bodyContent := fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())
	body := []byte(bodyContent)

	header.Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
