package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"jinya.de/goji/auth"
	"jinya.de/goji/database"
)

const RequestArtistKey = "request_artist"

func GetUserFromRequest(r *http.Request) *database.Artist {
	return r.Context().Value(RequestArtistKey).(*database.Artist)
}

func CheckAuthentication(next http.Handler) http.Handler {
	return CheckAuthenticationAndRole("")(next)
}

func CheckAuthenticationAndRole(role string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := auth.CheckRequestForUserByApiKey(r, role)
			token := auth.GetPlainApiTokenFromRequest(r)
			if err != nil && strings.HasPrefix(token, database.ApiTokenPrefix) {
				user, err = auth.CheckRequestForUserByApiToken(r, role)
			}
			if err != nil {
				if _, ok := errors.AsType[*auth.MissingPermissionsError](err); ok {
					w.WriteHeader(http.StatusForbidden)
				} else {
					w.WriteHeader(http.StatusUnauthorized)
				}
				json.NewEncoder(w).Encode(
					map[string]any{
						"success": false,
						"error": map[string]any{
							"message": err.Error(),
							"type":    "invalid-api-key",
						},
					},
				)
				return
			}

			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), RequestArtistKey, user)))
		})
	}
}
