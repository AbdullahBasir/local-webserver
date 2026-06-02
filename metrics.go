package main

import (
	"fmt"
	"net/http"
)

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
	cfg.fileserverHits.Store(0)
	if cfg.Platform != "dev" {
		respondWithError(w, 403, "Forbidden")
		return
	}
	err := cfg.dbQueries.ResetData(r.Context())
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Failed to reset data: %v", err))
		return
	}
	respondWithJSON(w, 200, struct{}{})
}
