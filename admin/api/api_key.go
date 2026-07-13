package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"jinya.de/goji/api/middleware"
	"jinya.de/goji/database"
)

func GetAllApiKeys(w http.ResponseWriter, r *http.Request) {
	artist := middleware.GetUserFromRequest(r)
	encoder := json.NewEncoder(w)
	apiKeys, err := database.FindApiKeysByArtist(artist.Id)
	if err != nil {
		error("Failed to fetch api keys", http.StatusInternalServerError, encoder, w)
		return
	}

	encoder.Encode(apiKeys)
}

func DeleteApiKey(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	apiKey, err := database.FindApiKeyByHashedKey(key)
	if err != nil {
		error("Api key not found", http.StatusNotFound, json.NewEncoder(w), w)
		return
	}

	err = apiKey.Delete()
	if err != nil {
		error("Api key not found", http.StatusNotFound, json.NewEncoder(w), w)
	}

	w.WriteHeader(http.StatusNoContent)
}
