package api

import (
	"net/http"
)

func ClearCache(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
