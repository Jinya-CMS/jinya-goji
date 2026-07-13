package api

import (
	"encoding/json"
	"net/http"

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
	w.WriteHeader(http.StatusNotImplemented)
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
