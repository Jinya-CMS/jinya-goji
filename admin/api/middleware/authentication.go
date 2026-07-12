package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"jinya.de/goji/auth"
	"jinya.de/goji/database"
)

const RequestArtistKey = "request_artist"

func GetUserFromRequest(r *http.Request) *database.Artist {
	return r.Context().Value(RequestArtistKey).(*database.Artist)
}

func CheckAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.CheckRequestForUser(r, "")
		if err != nil {
			if _, ok := errors.AsType[*auth.MissingPermissionsError](err); ok {
				w.WriteHeader(http.StatusForbidden)
			}
			json.NewEncoder(w).Encode(err)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), RequestArtistKey, user)))
	})
}

func CheckAuthenticationAndRole(role string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := auth.CheckRequestForUser(r, role)
			if err != nil {
				if _, ok := errors.AsType[*auth.MissingPermissionsError](err); ok {
					w.WriteHeader(http.StatusForbidden)
				}
				json.NewEncoder(w).Encode(err)
				return
			}

			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), RequestArtistKey, user)))
		})
	}
}
