package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"jinya.de/goji/database"
	"jinya.de/goji/storage"
)

func getAllArtists(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	artists, err := database.FindAllArtists()
	if err != nil {
		sendError("Artists not found", http.StatusInternalServerError, encoder, w)
		return
	}

	encoder.Encode(artists)
}

func getArtistById(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Invalid id", http.StatusBadRequest, encoder, w)
		return
	}

	artist, err := database.FindArtistById(id)
	if err != nil {
		sendError("Artist not found", http.StatusNotFound, encoder, w)
		return
	}

	encoder.Encode(artist)
}

type artistBody struct {
	ArtistName string   `json:"artistName"`
	Email      string   `json:"email"`
	Password   string   `json:"password"`
	Enabled    bool     `json:"enabled"`
	Roles      []string `json:"roles"`
}

func createArtist(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	body, err := ParseBody[artistBody](r)
	if err != nil {
		sendError("Invalid body", http.StatusBadRequest, encoder, w)
		return
	}

	artist := database.Artist{
		Email:                body.Email,
		Enabled:              body.Enabled,
		Roles:                body.Roles,
		ArtistName:           body.ArtistName,
		LoginMailEnabled:     true,
		NewDeviceMailEnabled: true,
	}
	err = artist.SetPassword(body.Password)
	if err != nil {
		sendError("Failed to hash password", http.StatusInternalServerError, encoder, w)
		return
	}

	err = artist.Create()
	if err != nil {
		sendError("Failed to save artist", http.StatusInternalServerError, encoder, w)
		return
	}

	w.WriteHeader(http.StatusCreated)
	encoder.Encode(artist)
}

func updateArtist(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Invalid id", http.StatusBadRequest, encoder, w)
		return
	}

	artist, err := database.FindArtistById(id)
	if err != nil {
		sendError("Artist not found", http.StatusNotFound, encoder, w)
		return
	}

	body, err := ParseBody[artistBody](r)
	if err != nil {
		sendError("Invalid body", http.StatusBadRequest, encoder, w)
		return
	}

	artist.Email = body.Email
	artist.Roles = body.Roles
	artist.ArtistName = body.ArtistName
	if body.Password != "" {
		err = artist.SetPassword(body.Password)
		if err != nil {
			sendError("Failed to hash password", http.StatusInternalServerError, encoder, w)
			return
		}
	}

	err = artist.Update()
	if err != nil {
		sendError("Failed to save artist", http.StatusInternalServerError, encoder, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteArtist(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Invalid id", http.StatusBadRequest, encoder, w)
		return
	}

	artist, err := database.FindArtistById(id)
	if err != nil {
		sendError("Artist not found", http.StatusNotFound, encoder, w)
		return
	}

	err = artist.Delete()
	if err != nil {
		sendError("Failed to delete artist", http.StatusInternalServerError, encoder, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func activateUser(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Invalid id", http.StatusBadRequest, encoder, w)
		return
	}

	artist, err := database.FindArtistById(id)
	if err != nil {
		sendError("Artist not found", http.StatusNotFound, encoder, w)
		return
	}

	artist.Enabled = true
	err = artist.Update()
	if err != nil {
		sendError("Failed to activate artist", http.StatusInternalServerError, encoder, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deactivateUser(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Invalid id", http.StatusBadRequest, encoder, w)
		return
	}

	artist, err := database.FindArtistById(id)
	if err != nil {
		sendError("Artist not found", http.StatusNotFound, encoder, w)
		return
	}

	artist.Enabled = false
	err = artist.Update()
	if err != nil {
		sendError("Failed to deactivate artist", http.StatusInternalServerError, encoder, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func updateUserProfilePicture(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Invalid id", http.StatusBadRequest, encoder, w)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		sendError("Invalid profile picture data", http.StatusBadRequest, encoder, w)
		return
	}

	err = storage.SaveProfilePicture(id, data)
	if err != nil {
		sendError("Failed to delete profile picture", http.StatusInternalServerError, encoder, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteUserProfilePicture(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Invalid id", http.StatusBadRequest, encoder, w)
		return
	}

	err = storage.DeleteProfilePicture(id)
	if err != nil {
		sendError("Failed to delete profile picture", http.StatusInternalServerError, encoder, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func resetTotp(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Invalid id", http.StatusBadRequest, encoder, w)
		return
	}

	artist, err := database.FindArtistById(id)
	if err != nil {
		sendError("Artists not found", http.StatusNotFound, encoder, w)
		return
	}

	artist.TotpSecret = nil
	artist.TotpMode = database.TotpModeEmail
	err = artist.Update()
	if err != nil {
		sendError("Failed to reset totp for artist", http.StatusInternalServerError, encoder, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
