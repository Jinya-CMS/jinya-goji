package auth

import (
	"log/slog"
	"net/http"
	"slices"
	"strings"
	"time"

	"jinya.de/goji/configuration"
	"jinya.de/goji/database"
)

const (
	AuthenticationCookieName = "JinyaApiKey"
	KnownDeviceCookieName    = "JinyaDeviceCode"
)

type ApiKeyInvalidError struct {
	Message string `json:"message"`
}

func (e *ApiKeyInvalidError) Error() string {
	return e.Message
}

type MissingPermissionsError struct {
	Message string `json:"message"`
	Role    string `json:"role"`
}

func (e *MissingPermissionsError) Error() string {
	return e.Message
}

func GetApiKeyFromRequest(r *http.Request) (*database.ApiKey, error) {
	key := getPlainApiKeyFromRequest(r)
	if key == "" {
		return nil, nil
	}

	return database.FindApiKeyByKey(key)
}

func getPlainApiKeyFromRequest(r *http.Request) string {
	authorizationHeader := r.Header.Get("Authorization")
	apiKey := ""
	if after, ok := strings.CutPrefix(authorizationHeader, "Bearer "); ok {
		apiKey = after
	} else {
		cookie, err := r.Cookie(AuthenticationCookieName)
		if err == nil {
			apiKey = cookie.Value
		}
	}

	if strings.HasPrefix(apiKey, database.ApiKeyPrefix) {
		return apiKey
	}

	return ""
}

func GetApiTokenFromRequest(r *http.Request) (*database.ApiToken, error) {
	token := GetPlainApiTokenFromRequest(r)
	if token == "" {
		return nil, nil
	}

	return database.FindApiTokenByToken(token)
}

func GetPlainApiTokenFromRequest(r *http.Request) string {
	authorizationHeader := r.Header.Get("Authorization")
	if after, ok := strings.CutPrefix(authorizationHeader, "Bearer "); ok {
		return after
	}

	return ""
}

func GetDeviceCodeFromRequest(r *http.Request) (*database.KnownDevice, error) {
	key := GetPlainDeviceCodeFromRequest(r)
	if key == "" {
		return nil, nil
	}

	return database.FindKnownDeviceByKey(key)
}

func GetPlainDeviceCodeFromRequest(r *http.Request) string {
	deviceCode := r.Header.Get(KnownDeviceCookieName)
	if deviceCode != "" {
		return deviceCode
	}

	cookie, err := r.Cookie(KnownDeviceCookieName)
	if err == nil {
		return cookie.Value
	}

	return ""
}

func CheckRequestForUserByApiKey(r *http.Request, role string) (*database.Artist, error) {
	slog.Debug("Check if api key exists in request")
	apiKey, err := GetApiKeyFromRequest(r)
	if err != nil {
		slog.Debug("Failed to get api key from request", "error", err)
		return nil, &ApiKeyInvalidError{Message: "Api key invalid"}
	}

	if apiKey == nil {
		slog.Debug("Api key not found")
		return nil, &ApiKeyInvalidError{Message: "Api key invalid"}
	}

	expireAfterSeconds := configuration.GetInteger("jinya", "api_key_expiry", 86400)

	slog.Debug("Check if api key is valid")
	if apiKey.ValidSince.Add(time.Duration(expireAfterSeconds) * time.Second).Before(time.Now()) {
		slog.Debug("Api key is invalid, delete it")
		_ = apiKey.Delete()
		return nil, &ApiKeyInvalidError{Message: "Api key invalid"}
	}

	apiKey.ValidSince = time.Now()
	artist, err := apiKey.GetArtist()
	if err != nil {
		slog.Debug("Failed to get artist for api key", "error", err)
		return nil, &ApiKeyInvalidError{Message: "Api key invalid"}
	}

	slog.Debug("Check if the api keys artist is enabled")
	if artist == nil || !artist.Enabled {
		slog.Debug("Api key is invalid, delete it")
		_ = apiKey.Delete()
		return nil, &ApiKeyInvalidError{Message: "Api key invalid"}
	}

	slog.Debug("Touch api key and refresh the validation")
	err = apiKey.Update()
	if err != nil {
		slog.Debug("Failed to touch the api key", "error", err)
		return nil, &ApiKeyInvalidError{Message: "Api key invalid"}
	}

	if role != "" {
		cascadedRole := ""
		if role == RoleReader {
			cascadedRole = RoleWriter
		}

		slog.Debug("Check if artist has required role")
		if !slices.Contains(artist.Roles, role) && (cascadedRole == "" || !slices.Contains(artist.Roles, cascadedRole)) {
			slog.Debug("Artist does not have the required role")
			return nil, &MissingPermissionsError{
				Message: "Not enough permissions",
				Role:    role,
			}
		}
	}

	return artist, nil
}

func CheckRequestForUserByApiToken(r *http.Request, role string) (*database.Artist, error) {
	slog.Debug("Check if api token exists in request")
	apiToken, err := GetApiTokenFromRequest(r)
	if err != nil {
		slog.Debug("Failed to get api token from request", "error", err)
		return nil, &ApiKeyInvalidError{Message: "Api token invalid"}
	}

	if apiToken == nil {
		slog.Debug("Api token not found")
		return nil, &ApiKeyInvalidError{Message: "Api token invalid"}
	}

	slog.Debug("Check if api key is valid")
	if apiToken.ExpiryDate != nil && apiToken.ExpiryDate.After(time.Now()) {
		slog.Debug("Api token is invalid, delete it")
		_ = apiToken.Delete()
		return nil, &ApiKeyInvalidError{Message: "Api token invalid"}
	}

	artist, err := apiToken.GetArtist()
	if err != nil {
		slog.Debug("Failed to get artist for api key", "error", err)
		return nil, &ApiKeyInvalidError{Message: "Api key invalid"}
	}

	slog.Debug("Check if the api keys artist is enabled")
	if artist == nil || !artist.Enabled {
		slog.Debug("Api key is invalid, delete it")
		_ = apiToken.Delete()
		return nil, &ApiKeyInvalidError{Message: "Api key invalid"}
	}

	if role != "" {
		cascadedRole := ""
		if role == RoleReader {
			cascadedRole = RoleWriter
		}

		slog.Debug("Check if artist has required role")
		if !slices.Contains(artist.Roles, role) && (cascadedRole == "" || !slices.Contains(artist.Roles, cascadedRole)) {
			slog.Debug("Artist does not have the required role")
			return nil, &MissingPermissionsError{
				Message: "Not enough permissions",
				Role:    role,
			}
		}
	}

	return artist, nil
}
