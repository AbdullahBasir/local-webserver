package main

import (
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
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	hit_num := cfg.fileserverHits.Load()
	resp := fmt.Sprintf("Hits: %d", hit_num)
	w.Write([]byte(resp))
}

func (cfg *apiConfig) HandleReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
}

func main() {

	ServeMux := http.NewServeMux()

	apiCfg := &apiConfig{}

	handler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	ServeMux.Handle("/app/", apiCfg.middleware(handler))
	ServeMux.HandleFunc("/healthz", handlerReadiness)
	ServeMux.HandleFunc("/metrics", apiCfg.HandleMetrics)
	ServeMux.HandleFunc("/reset", apiCfg.HandleReset)

	serverStruct := &http.Server{
		Addr:    ":8080",
		Handler: ServeMux,
	}

	log.Fatal(serverStruct.ListenAndServe())
}
