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
		sendError("Artists not found", http.StatusInternalServerError, encoder, w, err)
		return
	}

	encoder.Encode(artists)
}

func getArtistById(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Invalid id", http.StatusBadRequest, encoder, w, err)
		return
	}

	artist, err := database.FindArtistById(id)
	if err != nil {
		sendError("Artist not found", http.StatusNotFound, encoder, w, err)
		return
	}

	encoder.Encode(&artist)
}

func createArtist(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	body, err := BodyFromRequest(r, "email", "artistName", "password")
	if err != nil {
		sendError("Invalid body", http.StatusBadRequest, encoder, w, err)
		return
	}

	artist := database.Artist{
		LoginMailEnabled:     true,
		NewDeviceMailEnabled: true,
	}

	if email, ok := GetFromBody[string](body, "email"); ok {
		artist.Email = *email
	}
	if name, ok := GetFromBody[string](body, "artistName"); ok {
		artist.ArtistName = *name
	}
	if password, ok := GetFromBody[string](body, "password"); ok {
		err = artist.SetPassword(*password)
		if err != nil {
			sendError("Failed to hash password", http.StatusInternalServerError, encoder, w, err)
			return
		}
	}
	if enabled, ok := GetFromBody[bool](body, "enabled"); ok {
		artist.Enabled = *enabled
	} else {
		artist.Enabled = false
	}
	if roles, ok := GetArrayFromBody[string](body, "roles"); ok {
		artist.Roles = roles
	}

	err = artist.Create()
	if err != nil {
		sendError("Failed to save artist", http.StatusInternalServerError, encoder, w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	encoder.Encode(&artist)
}

func updateArtist(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Invalid id", http.StatusBadRequest, encoder, w, err)
		return
	}

	artist, err := database.FindArtistById(id)
	if err != nil {
		sendError("Artist not found", http.StatusNotFound, encoder, w, err)
		return
	}

	body, err := BodyFromRequest(r)
	if err != nil {
		sendError("Invalid body", http.StatusBadRequest, encoder, w, err)
		return
	}

	if email, ok := GetFromBody[string](body, "email"); ok {
		artist.Email = *email
	}
	if name, ok := GetFromBody[string](body, "artistName"); ok {
		artist.ArtistName = *name
	}
	if password, ok := GetFromBody[string](body, "password"); ok {
		err = artist.SetPassword(*password)
		if err != nil {
			sendError("Failed to hash password", http.StatusInternalServerError, encoder, w, err)
			return
		}
	}
	if enabled, ok := GetFromBody[bool](body, "enabled"); ok {
		artist.Enabled = *enabled
	} else {
		artist.Enabled = false
	}
	if roles, ok := GetArrayFromBody[string](body, "roles"); ok {
		artist.Roles = roles
	}

	err = artist.Update()
	if err != nil {
		sendError("Failed to save artist", http.StatusInternalServerError, encoder, w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteArtist(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Invalid id", http.StatusBadRequest, encoder, w, err)
		return
	}

	artist, err := database.FindArtistById(id)
	if err != nil {
		sendError("Artist not found", http.StatusNotFound, encoder, w, err)
		return
	}

	err = artist.Delete()
	if err != nil {
		sendError("Failed to delete artist", http.StatusInternalServerError, encoder, w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func activateUser(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Invalid id", http.StatusBadRequest, encoder, w, err)
		return
	}

	artist, err := database.FindArtistById(id)
	if err != nil {
		sendError("Artist not found", http.StatusNotFound, encoder, w, err)
		return
	}

	artist.Enabled = true
	err = artist.Update()
	if err != nil {
		sendError("Failed to activate artist", http.StatusInternalServerError, encoder, w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deactivateUser(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Invalid id", http.StatusBadRequest, encoder, w, err)
		return
	}

	artist, err := database.FindArtistById(id)
	if err != nil {
		sendError("Artist not found", http.StatusNotFound, encoder, w, err)
		return
	}

	artist.Enabled = false
	err = artist.Update()
	if err != nil {
		sendError("Failed to deactivate artist", http.StatusInternalServerError, encoder, w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func updateUserProfilePicture(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Invalid id", http.StatusBadRequest, encoder, w, err)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		sendError("Invalid profile picture data", http.StatusBadRequest, encoder, w, err)
		return
	}

	err = storage.SaveProfilePicture(id, data)
	if err != nil {
		sendError("Failed to delete profile picture", http.StatusInternalServerError, encoder, w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteUserProfilePicture(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Invalid id", http.StatusBadRequest, encoder, w, err)
		return
	}

	err = storage.DeleteProfilePicture(id)
	if err != nil {
		sendError("Failed to delete profile picture", http.StatusInternalServerError, encoder, w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func resetTotp(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Invalid id", http.StatusBadRequest, encoder, w, err)
		return
	}

	artist, err := database.FindArtistById(id)
	if err != nil {
		sendError("Artists not found", http.StatusNotFound, encoder, w, err)
		return
	}

	artist.TotpSecret = nil
	artist.TotpMode = database.TotpModeEmail
	err = artist.Update()
	if err != nil {
		sendError("Failed to reset totp for artist", http.StatusInternalServerError, encoder, w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
