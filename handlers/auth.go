package handlers

import (
	"log"
	"net/http"
)

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("AuthHandler called")

	w.Write([]byte("Auth functionality not implemented"))
}
