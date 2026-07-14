package api

import (
	"encoding/json"
	"io"
	"net/http"
	"slices"

	"github.com/pquerna/otp/totp"
	"jinya.de/goji/api/middleware"
	"jinya.de/goji/database"
	"jinya.de/goji/storage"
)

func getMe(w http.ResponseWriter, r *http.Request) {
	artist := middleware.GetUserFromRequest(r)
	encoder := json.NewEncoder(w)
	encoder.Encode(artist)
}

func updateMe(w http.ResponseWriter, r *http.Request) {
	artist := middleware.GetUserFromRequest(r)
	encoder := json.NewEncoder(w)

	body, err := BodyFromRequest(r)
	if err != nil {
		sendError("Invalid body", http.StatusBadRequest, encoder, w, err)
		return
	}

	if email, ok := GetFromBody[string](body, "email"); ok {
		artist.Email = *email
	}
	if artistName, ok := GetFromBody[string](body, "artistName"); ok {
		artist.ArtistName = *artistName
	}
	if aboutMe, ok := GetFromBody[string](body, "aboutMe"); ok {
		artist.AboutMe = aboutMe
	}

	err = artist.Update()
	if err != nil {
		sendError("Failed to update profile", http.StatusInternalServerError, encoder, w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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
	artist := middleware.GetUserFromRequest(r)
	encoder := json.NewEncoder(w)

	body, err := BodyFromRequest(r)
	if err != nil {
		sendError("Invalid body", http.StatusBadRequest, encoder, w, err)
		return
	}

	if loginMailEnabled, ok := GetFromBody[bool](body, "loginMailEnabled"); ok {
		artist.LoginMailEnabled = *loginMailEnabled
	}
	if newDeviceMailEnabled, ok := GetFromBody[bool](body, "newDeviceMailEnabled"); ok {
		artist.NewDeviceMailEnabled = *newDeviceMailEnabled
	}

	err = artist.Update()
	if err != nil {
		sendError("Failed to update preferences", http.StatusInternalServerError, encoder, w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func updateMeProfilePicture(w http.ResponseWriter, r *http.Request) {
	artist := middleware.GetUserFromRequest(r)
	encoder := json.NewEncoder(w)

	data, err := io.ReadAll(r.Body)
	if err != nil {
		sendError("Failed to read body", http.StatusBadRequest, encoder, w, err)
		return
	}

	err = storage.SaveProfilePicture(artist.Id, data)
	if err != nil {
		sendError("Failed to save profile picture", http.StatusInternalServerError, encoder, w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func verifyOtpApp(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	artist := middleware.GetUserFromRequest(r)

	body, err := BodyFromRequest(r, "code")
	if err != nil {
		sendError("Invalid body", http.StatusBadRequest, encoder, w, err)
		return
	}

	code, _ := GetFromBody[string](body, "code")
	valid := totp.Validate(*code, *artist.TotpSecret)
	if valid {
		artist.TotpMode = database.TotpModeApp
		err = artist.Update()
		if err != nil {
			sendError("Failed to save totp mode", http.StatusInternalServerError, encoder, w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	} else {
		sendError("Code is invalid", http.StatusForbidden, encoder, w, nil)
	}
}

func setupOtpApp(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	artist := middleware.GetUserFromRequest(r)
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Jinya Goji",
		AccountName: artist.ArtistName + " (" + artist.Email + ")",
	})
	if err != nil {
		sendError("Failed to generate totp key", http.StatusInternalServerError, encoder, w, err)
		return
	}

	artist.TotpSecret = new(key.Secret())
	err = artist.Update()
	if err != nil {
		sendError("Failed to save artist", http.StatusInternalServerError, encoder, w, err)
		return
	}

	encoder.Encode(map[string]string{
		"url": key.URL(),
	})
}
