package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) HandleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	hit_num := cfg.fileserverHits.Load()
	resp := fmt.Sprintf(`
	<html>
  		<body>
    		<h1>Welcome, Chirpy Admin</h1>
    		<p>Chirpy has been visited %d times!</p>
  		</body>
	</html>`, hit_num)
	w.Write([]byte(resp))
}

func (cfg *apiConfig) HandleReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
}

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) error {
	response, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(status)
	w.Write(response)
	return nil
}

func respondWithError(w http.ResponseWriter, status int, message string) error {
	return respondWithJSON(w, status, map[string]string{"error": message})
}

func chirpValidationHandler(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(r.Body)
	reqBody := requestBody{}
	err := decoder.Decode(&reqBody)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Invalid request body: %v", err))
		return
	} else if len(reqBody.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	type returnVals struct {
		Valid bool `json:"valid"`
	}

	valBody := returnVals{
		Valid: true,
	}

	err = respondWithJSON(w, 200, valBody)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error marshaling response: %v", err))
		return
	}
}

func main() {

	ServeMux := http.NewServeMux()

	apiCfg := &apiConfig{}

	handler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	ServeMux.Handle("/app/", apiCfg.middleware(handler))
	ServeMux.HandleFunc("GET /api/healthz", handlerReadiness)
	ServeMux.HandleFunc("GET /admin/metrics", apiCfg.HandleMetrics)
	ServeMux.HandleFunc("POST /admin/reset", apiCfg.HandleReset)
	ServeMux.HandleFunc("POST /api/validate_chirp", chirpValidationHandler)

	serverStruct := &http.Server{
		Addr:    ":8080",
		Handler: ServeMux,
	}

	log.Fatal(serverStruct.ListenAndServe())
}
