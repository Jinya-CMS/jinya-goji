package api

import (
	"encoding/json"
	"net/http"
)

const (
	ErrorTypeInternalError = "internal-error"
	ErrorTypeNotFound      = "not-found"
	ErrorTypeBadRequest    = "bad-request"
)

func sendError(message string, statusCode int, encoder *json.Encoder, w http.ResponseWriter) {
	w.WriteHeader(statusCode)

	errorType := "unknown-error"
	switch statusCode {
	case http.StatusInternalServerError:
		errorType = "internal-error"
		break
	case http.StatusBadRequest:
		errorType = "bad-request"
		break
	case http.StatusNotFound:
		errorType = "not-found"
		break
	case http.StatusUnauthorized:
		errorType = "unauthorized"
		break
	case http.StatusForbidden:
		errorType = "forbidden"
		break
	}

	encoder.Encode(
		map[string]any{
			"success": false,
			"error": map[string]any{
				"message": message,
				"type":    errorType,
			},
		},
	)
}
