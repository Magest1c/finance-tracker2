package main

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, data any) error {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)

}

func writeError(w http.ResponseWriter, status int, message string) error {

	response := errorResponse{
		Error: message,
	}
	return writeJSON(w, status, response)

}
