package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	ErrorTypeInternalError = "internal-error"
	ErrorTypeNotFound      = "not-found"
	ErrorTypeBadRequest    = "bad-request"
)

const (
	errorTypeUnique = iota
	errorTypeReference
	errorTypeNotNull
	errorTypeOther
)

func detectErrorType(err error) int {
	pgErr, ok := errors.AsType[*pgconn.PgError](err)
	if !ok {
		return errorTypeOther
	}

	if pgErr.Code == pgerrcode.ForeignKeyViolation {
		return errorTypeReference
	}
	if pgErr.Code == pgerrcode.UniqueViolation {
		return errorTypeUnique
	}
	if pgErr.Code == pgerrcode.NotNullViolation {
		return errorTypeNotNull
	}

	return errorTypeOther
}

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
