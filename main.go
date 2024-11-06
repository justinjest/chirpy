package main

import (
	"log"
	"net/http"
	"sync/atomic"

	jsonParser "github.com/justinjest/chirpy/internal/json"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}
	mux := http.NewServeMux()
	handler := http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app/", (apiCfg.middlewareMetricsInc(handler)))
	mux.HandleFunc("GET /api/healthz", healthzfunc)
	mux.HandleFunc("GET /admin/metrics", apiCfg.hitsHandler)
	mux.HandleFunc("POST /api/reset", apiCfg.hitsReset)
	mux.HandleFunc("POST /api/validate_chirp", jsonParser.Handler)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
