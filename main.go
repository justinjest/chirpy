package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func healthzfunc(w http.ResponseWriter, req *http.Request) {
	header := w.Header()
	body := []byte("OK")
	header.Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) hitsReset(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)

	header := w.Header()
	bodyContent := fmt.Sprintf("Hits: %v", cfg.fileserverHits)
	body := []byte(bodyContent)

	header.Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (cfg *apiConfig) hitsHandler(w http.ResponseWriter, req *http.Request) {
	header := w.Header()
	bodyContent := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileserverHits.Load())
	body := []byte(bodyContent)

	header.Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	var apiCfg apiConfig
	mux := http.NewServeMux()
	handler := http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app/", (apiCfg.middlewareMetricsInc(handler)))
	mux.HandleFunc("GET /api/healthz", healthzfunc)
	mux.HandleFunc("GET /admin/metrics", apiCfg.hitsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.hitsReset)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
