package api

import (
	"encoding/json"
	"net/http"
)

func ParseBody[T any](r *http.Request) (*T, error) {
	decoder := json.NewDecoder(r.Body)
	var body T
	err := decoder.Decode(&body)

	return &body, err
}
