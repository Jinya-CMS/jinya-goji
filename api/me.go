package api

import (
	"encoding/json"
	"net/http"
	"slices"

	"jinya.de/goji/api/middleware"
)

func getMe(w http.ResponseWriter, r *http.Request) {
	artist := middleware.GetUserFromRequest(r)
	encoder := json.NewEncoder(w)
	encoder.Encode(artist)
}

func updateMe(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func updateMeColorScheme(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	decoder := json.NewDecoder(r.Body)

	var body struct {
		ColorScheme string `json:"colorScheme"`
	}
	err := decoder.Decode(&body)
	if err != nil || !slices.Contains([]string{"auto", "light", "dark"}, body.ColorScheme) {
		sendError("Invalid body", http.StatusBadRequest, encoder, w, err)
		return
	}

	artist := middleware.GetUserFromRequest(r)
	switch body.ColorScheme {
	case "auto":
		artist.PrefersColorScheme = nil
		break
	case "light":
		artist.PrefersColorScheme = new(false)
		break
	case "dark":
		artist.PrefersColorScheme = new(true)
		break
	}

	err = artist.Update()
	if err != nil {
		sendError("Failed to update color scheme", http.StatusInternalServerError, encoder, w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func updateMePreferences(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func updateMeProfilePicture(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func verifyOtpApp(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func setupOtpApp(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
