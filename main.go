package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/google/uuid"

	"github.com/AbdullahBasir/local-webserver/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	Platform       string
	Secret         string
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

type RefreshToken struct {
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	RevokedAt time.Time `json:"revoked_at"`
}

func main() {

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	dbQueries := database.New(db)

	ServeMux := http.NewServeMux()

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		dbQueries:      dbQueries,
		Platform:       os.Getenv("PLATFORM"),
		Secret:         os.Getenv("JWT_SECRET"),
	}

	handler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	ServeMux.Handle("/app/", apiCfg.middleware(handler))
	ServeMux.HandleFunc("GET /api/healthz", handlerReadiness)
	ServeMux.HandleFunc("GET /admin/metrics", apiCfg.HandleMetrics)
	ServeMux.HandleFunc("POST /admin/reset", apiCfg.HandleReset)
	ServeMux.HandleFunc("POST /api/users", apiCfg.HandleCreateUser)
	ServeMux.HandleFunc("POST /api/chirps", apiCfg.chirpWriter)
	ServeMux.HandleFunc("GET /api/chirps", apiCfg.RetrieveChirps)
	ServeMux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.GetSingleChirp)
	ServeMux.HandleFunc("POST /api/login", apiCfg.HandleUserLogin)

	serverStruct := &http.Server{
		Addr:    ":8080",
		Handler: ServeMux,
	}

	log.Fatal(serverStruct.ListenAndServe())
}
