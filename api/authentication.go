package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"jinya.de/goji/api/middleware"
	"jinya.de/goji/auth"
	"jinya.de/goji/database"
	"jinya.de/goji/mailing"
)

func login(w http.ResponseWriter, r *http.Request) {
	type loginBody struct {
		Username      string `json:"username"`
		Password      string `json:"password"`
		TwoFactorCode string `json:"twoFactorCode"`
	}

	encoder := json.NewEncoder(w)
	body, err := ParseBody[loginBody](r)
	if err != nil {
		sendError("Failed to parse body", http.StatusBadRequest, encoder, w, err)
		return
	}

	deviceCode, _ := auth.GetDeviceCodeFromRequest(r)
	if deviceCode.DeviceKey == "" && body.TwoFactorCode == "" {
		sendError("Bad Credentials", http.StatusForbidden, encoder, w, err)
		return
	}

	artist, err := database.FindArtistByEmail(body.Username)
	if err != nil {
		sendError("Bad Credentials", http.StatusForbidden, encoder, w, err)
		return
	}

	if !artist.ValidatePassword(body.Password) {
		_ = artist.RegisterFailedLogin()
		sendError("Bad Credentials", http.StatusForbidden, encoder, w, err)
		return
	}

	sendApiKey := func() {
		key, apiKey, err := database.NewApiKey(artist, r)
		if err != nil {
			sendError("Bad Credentials", http.StatusForbidden, encoder, w, err)
			return
		}

		err = apiKey.Create()
		if err != nil {
			sendError("Bad Credentials", http.StatusForbidden, encoder, w, err)
			return
		}

		var code string
		var knownDevice *database.KnownDevice
		if deviceCode.DeviceKey == "" {
			code, knownDevice, err = database.NewKnownDevice(artist, r)
			if err != nil {
				sendError("Bad Credentials", http.StatusForbidden, encoder, w, err)
				return
			}

			err = knownDevice.Create()
			if err != nil {
				sendError("Bad Credentials", http.StatusForbidden, encoder, w, err)
				return
			}

			if artist.NewDeviceMailEnabled {
				go mailing.SendNewDeviceMail(r, artist.ArtistName, artist.Email)
			}
		} else {
			code = auth.GetPlainDeviceCodeFromRequest(r)
		}

		if artist.LoginMailEnabled {
			go mailing.SendNewLoginMail(r, artist.ArtistName, artist.Email)
		}

		http.SetCookie(w, &http.Cookie{
			Name:     auth.AuthenticationCookieName,
			Value:    key,
			Expires:  time.Now().Add(10 * 365 * 24 * time.Hour),
			Secure:   r.TLS != nil,
			Path:     "/",
			HttpOnly: true,
		})
		http.SetCookie(w, &http.Cookie{
			Name:     auth.KnownDeviceCookieName,
			Value:    code,
			Expires:  time.Now().Add(10 * 365 * 24 * time.Hour),
			Secure:   r.TLS != nil,
			Path:     "/",
			HttpOnly: false,
		})

		encoder.Encode(map[string]any{
			"apiKey":     key,
			"deviceCode": code,
			"roles":      artist.Roles,
		})
	}

	if deviceCode != nil && artist.ValidateDevice(deviceCode.DeviceKey) {
		sendApiKey()
		return
	}

	if artist.VerifyTotpCode(body.TwoFactorCode) {
		sendApiKey()
		return
	}

	sendError("Bad Credentials", http.StatusForbidden, encoder, w, err)
}

func request2fa(w http.ResponseWriter, r *http.Request) {
	type loginBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	encoder := json.NewEncoder(w)

	body, err := ParseBody[loginBody](r)
	if err != nil {
		sendError("Bad Credentials", http.StatusForbidden, encoder, w, err)
		return
	}

	artist, err := database.FindArtistByEmail(body.Username)
	if err != nil {
		sendError("Bad Credentials", http.StatusForbidden, encoder, w, err)
		return
	}

	if !artist.ValidatePassword(body.Password) {
		_ = artist.RegisterFailedLogin()
		sendError("Bad Credentials", http.StatusForbidden, encoder, w, err)
		return
	}

	if artist.Enabled && artist.TotpMode == database.TotpModeApp && artist.TotpSecret != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if artist.Enabled {
		err = artist.SetTwoFactorCode()
		if err != nil {
			sendError("Bad Credentials", http.StatusForbidden, encoder, w, err)
			return
		}

		err = artist.Update()
		if err != nil {
			sendError("Bad Credentials", http.StatusForbidden, encoder, w, err)
			return
		}

		go mailing.SendTwoFactorCodeMail(r, artist.ArtistName, artist.Email, *artist.TwoFactorToken)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	sendError("Bad Credentials", http.StatusForbidden, encoder, w, err)
}

func logout(w http.ResponseWriter, r *http.Request) {
	fully := r.URL.Query().Get("fully") == "true"
	apiKey, err := auth.GetApiKeyFromRequest(r)
	if err != nil || apiKey == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = apiKey.Delete()
	if err != nil {
		slog.Error("Failed to delete the api key", "error", err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     auth.AuthenticationCookieName,
		Value:    "",
		MaxAge:   -1,
		Secure:   r.TLS != nil,
		HttpOnly: true,
	})
	if fully {
		deviceCode, err := auth.GetDeviceCodeFromRequest(r)
		if err != nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		err = deviceCode.Delete()
		if err != nil {
			slog.Error("Failed to delete the known device", "error", err)
		}

		http.SetCookie(w, &http.Cookie{
			Name:     auth.KnownDeviceCookieName,
			Value:    "",
			MaxAge:   -1,
			Secure:   r.TLS != nil,
			HttpOnly: true,
		})
	}
}

func checkApiKey(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func updatePassword(w http.ResponseWriter, r *http.Request) {
	artist := middleware.GetUserFromRequest(r)
	type passwordBody struct {
		Password    string `json:"password"`
		OldPassword string `json:"oldPassword"`
	}

	encoder := json.NewEncoder(w)

	body, err := ParseBody[passwordBody](r)
	if err != nil {
		sendError("Failed to parse body", http.StatusBadRequest, encoder, w, err)
		return
	}

	err = artist.ChangePassword(body.OldPassword, body.Password)
	if err != nil {
		sendError("Failed to change password", http.StatusInternalServerError, encoder, w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
