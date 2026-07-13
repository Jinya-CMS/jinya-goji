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

func GetAllArtists(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	artists, err := database.FindAllArtists()
	if err != nil {
		error("Artists not found", http.StatusInternalServerError, encoder, w)
		return
	}

	encoder.Encode(artists)
}

func GetArtistById(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		error("Invalid id", http.StatusBadRequest, encoder, w)
		return
	}

	artist, err := database.FindArtistById(id)
	if err != nil {
		error("Artist not found", http.StatusNotFound, encoder, w)
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

func CreateArtist(w http.ResponseWriter, r *http.Request) {
	var body artistBody
	encoder := json.NewEncoder(w)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		error("Invalid body", http.StatusBadRequest, encoder, w)
		return
	}

	artist := database.Artist{
		Email:      body.Email,
		Enabled:    body.Enabled,
		Roles:      body.Roles,
		ArtistName: body.ArtistName,
	}
	err = artist.SetPassword(body.Password)
	if err != nil {
		error("Failed to hash password", http.StatusInternalServerError, encoder, w)
		return
	}

	err = artist.Create()
	if err != nil {
		error("Failed to save artist", http.StatusInternalServerError, encoder, w)
		return
	}

	w.WriteHeader(http.StatusCreated)
	encoder.Encode(artist)
}

func UpdateArtist(w http.ResponseWriter, r *http.Request) {
	var body artistBody
	encoder := json.NewEncoder(w)
	decoder := json.NewDecoder(r.Body)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		error("Invalid id", http.StatusBadRequest, encoder, w)
		return
	}

	artist, err := database.FindArtistById(id)
	if err != nil {
		error("Artist not found", http.StatusNotFound, encoder, w)
		return
	}

	err = decoder.Decode(&body)
	if err != nil {
		error("Invalid body", http.StatusBadRequest, encoder, w)
		return
	}

	artist.Email = body.Email
	artist.Roles = body.Roles
	artist.ArtistName = body.ArtistName
	if body.Password != "" {
		err = artist.SetPassword(body.Password)
		if err != nil {
			error("Failed to hash password", http.StatusInternalServerError, encoder, w)
			return
		}
	}

	err = artist.Update()
	if err != nil {
		error("Failed to save artist", http.StatusInternalServerError, encoder, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteArtist(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		error("Invalid id", http.StatusBadRequest, encoder, w)
		return
	}

	artist, err := database.FindArtistById(id)
	if err != nil {
		error("Artist not found", http.StatusNotFound, encoder, w)
		return
	}

	err = artist.Delete()
	if err != nil {
		error("Failed to delete artist", http.StatusInternalServerError, encoder, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func ActivateUser(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		error("Invalid id", http.StatusBadRequest, encoder, w)
		return
	}

	artist, err := database.FindArtistById(id)
	if err != nil {
		error("Artist not found", http.StatusNotFound, encoder, w)
		return
	}

	artist.Enabled = true
	err = artist.Update()
	if err != nil {
		error("Failed to activate artist", http.StatusInternalServerError, encoder, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeactivateUser(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		error("Invalid id", http.StatusBadRequest, encoder, w)
		return
	}

	artist, err := database.FindArtistById(id)
	if err != nil {
		error("Artist not found", http.StatusNotFound, encoder, w)
		return
	}

	artist.Enabled = false
	err = artist.Update()
	if err != nil {
		error("Failed to deactivate artist", http.StatusInternalServerError, encoder, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func UpdateUserProfilePicture(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		error("Invalid id", http.StatusBadRequest, encoder, w)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		error("Invalid profile picture data", http.StatusBadRequest, encoder, w)
		return
	}

	err = storage.SaveProfilePicture(id, data)
	if err != nil {
		error("Failed to delete profile picture", http.StatusInternalServerError, encoder, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteUserProfilePicture(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		error("Invalid id", http.StatusBadRequest, encoder, w)
		return
	}

	err = storage.DeleteProfilePicture(id)
	if err != nil {
		error("Failed to delete profile picture", http.StatusInternalServerError, encoder, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func ResetTotp(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		error("Invalid id", http.StatusBadRequest, encoder, w)
		return
	}

	artist, err := database.FindArtistById(id)
	if err != nil {
		error("Artists not found", http.StatusNotFound, encoder, w)
		return
	}

	artist.TotpSecret = nil
	artist.TotpMode = database.TotpModeEmail
	err = artist.Update()
	if err != nil {
		error("Failed to reset totp for artist", http.StatusInternalServerError, encoder, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
