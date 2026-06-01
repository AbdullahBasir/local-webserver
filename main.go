package main

import (
	"log"
	"net/http"
)

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
