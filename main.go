package main

import (
	"log"
	"net/http"
)

func main() {

	ServeMux := http.NewServeMux()

	serverStruct := &http.Server{
		Addr:    ":8080",
		Handler: ServeMux,
	}
	err := serverStruct.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
