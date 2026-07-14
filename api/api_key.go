package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"jinya.de/goji/api/middleware"
	"jinya.de/goji/database"
)

func getAllApiKeys(w http.ResponseWriter, r *http.Request) {
	artist := middleware.GetUserFromRequest(r)
	encoder := json.NewEncoder(w)
	apiKeys, err := database.FindApiKeysByArtist(artist.Id)
	if err != nil {
		sendError("Failed to fetch api keys", http.StatusInternalServerError, encoder, w, err)
		return
	}

	encoder.Encode(apiKeys)
}

func deleteApiKey(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	encoder := json.NewEncoder(w)
	apiKey, err := database.FindApiKeyByHashedKey(key)
	if err != nil {
		sendError("Api key not found", http.StatusNotFound, encoder, w, err)
		return
	}

	err = apiKey.Delete()
	if err != nil {
		sendError("Api key not found", http.StatusNotFound, encoder, w, err)
	}

	w.WriteHeader(http.StatusNoContent)
}
