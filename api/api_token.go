package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"jinya.de/goji/api/middleware"
	"jinya.de/goji/database"
)

func getAllApiTokens(w http.ResponseWriter, r *http.Request) {
	artist := middleware.GetUserFromRequest(r)
	encoder := json.NewEncoder(w)
	apiTokens, err := database.FindApiTokensByArtist(artist.Id)
	if err != nil {
		sendError("Failed to fetch api tokens", http.StatusInternalServerError, encoder, w, err)
		return
	}

	encoder.Encode(apiTokens)
}

func createApiToken(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	encoder := json.NewEncoder(w)

	var body struct {
		Name       string     `json:"name"`
		ExpiryDate *time.Time `json:"expiryDate"`
	}
	err := decoder.Decode(&body)
	if err != nil {
		sendError("Failed to decode body", http.StatusBadRequest, encoder, w, err)
		return
	}

	persistentToken, apiToken, err := database.NewApiToken(body.Name, body.ExpiryDate, middleware.GetUserFromRequest(r))
	if err != nil {
		sendError("Failed to create api token", http.StatusInternalServerError, encoder, w, err)
		return
	}

	err = apiToken.Create()
	if err != nil {
		sendError("Failed to create api token", http.StatusInternalServerError, encoder, w, err)
		return
	}

	apiToken.Token = persistentToken

	w.WriteHeader(http.StatusCreated)
	encoder.Encode(apiToken)
}

func deleteApiToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]

	encoder := json.NewEncoder(w)
	apiToken, err := database.FindApiTokenByHashedToken(token)
	if err != nil {
		sendError("Api token not found", http.StatusNotFound, encoder, w, err)
		return
	}

	err = apiToken.Delete()
	if err != nil {
		sendError("Api token not found", http.StatusNotFound, encoder, w, err)
	}

	w.WriteHeader(http.StatusNoContent)
}
