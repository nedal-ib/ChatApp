package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func JSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func JSONError(w http.ResponseWriter, message string, statusCode int) {
	type ErrorResponse struct {
		Error string `json:"error"`
	}
	JSONResponse(w, ErrorResponse{Error: message}, statusCode)
}
