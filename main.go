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
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
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
	}

	handler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	ServeMux.Handle("/app/", apiCfg.middleware(handler))
	ServeMux.HandleFunc("GET /api/healthz", handlerReadiness)
	ServeMux.HandleFunc("GET /admin/metrics", apiCfg.HandleMetrics)
	ServeMux.HandleFunc("POST /admin/reset", apiCfg.HandleReset)
	ServeMux.HandleFunc("POST /api/validate_chirp", chirpValidationHandler)
	ServeMux.HandleFunc("POST /api/users", apiCfg.HandleCreateUser)

	serverStruct := &http.Server{
		Addr:    ":8080",
		Handler: ServeMux,
	}

	log.Fatal(serverStruct.ListenAndServe())
}
