package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	// TODO: Register payment handlers
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}
	log.Printf("Payment service running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
