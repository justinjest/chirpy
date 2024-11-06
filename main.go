package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	"github.com/justinjest/chirpy/internal/database"
	jsonParser "github.com/justinjest/chirpy/internal/json"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	database       *database.Queries
}

func main() {

	const filepathRoot = "."
	const port = "8080"
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Unable to load .env file")
	}
	dbUrl := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("Unable to load postgress file")
	}
	dbQueries := database.New(db)
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		database:       dbQueries,
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
