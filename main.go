package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/justinjest/chirpy/internal/database"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	database       *database.Queries
	PLATFORM       string
	secret         string
	apiKey         string
}

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	RedStatus    bool      `json:"is_chirpy_red"`
}

type RefreshToken struct {
	ID        string      `json:"id"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	UserID    uuid.UUID   `json:"user_id"`
	Expires   time.Time   `json:"expires_at"`
	Revoked   pq.NullTime `json:"revoked_at"`
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
		PLATFORM:       os.Getenv("PLATFORM"),
		secret:         os.Getenv("SECRET"),
		apiKey:         os.Getenv("API_KEY"),
	}
	mux := http.NewServeMux()
	handler := http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app/", (apiCfg.middlewareMetricsInc(handler)))
	mux.HandleFunc("GET /api/healthz", healthzfunc)
	mux.HandleFunc("GET /admin/metrics", apiCfg.hitsHandler)
	mux.HandleFunc("POST /api/reset", apiCfg.hitsReset)
	mux.HandleFunc("POST /api/users", apiCfg.CreateUser)
	mux.HandleFunc("POST /api/chirps", apiCfg.createChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.getChirps)
	mux.HandleFunc("POST /admin/reset", apiCfg.DropUsers)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.getOneChirp)
	mux.HandleFunc("POST /api/login", apiCfg.userLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.refreshUser)
	mux.HandleFunc("POST /api/revoke", apiCfg.revokeUser)
	mux.HandleFunc("PUT /api/users", apiCfg.updateUser)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.deleteChirp)
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.updateUserRed)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
