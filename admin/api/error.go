package api

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	Message string `json:"message"`
}

func error(message string, statusCode int, encoder *json.Encoder, w http.ResponseWriter) {
	w.WriteHeader(statusCode)
	encoder.Encode(errorResponse{Message: message})
}
