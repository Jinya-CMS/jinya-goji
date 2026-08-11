package engine

import (
	_ "embed"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"jinya.de/goji/analytics"
	"jinya.de/goji/database"
)

func ForwardBlogRequest(w http.ResponseWriter, r *http.Request) {
	statusCode := forwardRequest(w, r)
	if statusCode > 199 && 300 > statusCode {
		go analytics.TrackBlogPostRequest(r)
	}
}

func ForwardStandardRequest(w http.ResponseWriter, r *http.Request) {
	statusCode := forwardRequest(w, r)
	if statusCode > 199 && 300 > statusCode {
		go analytics.TrackStandardRequest(r)
	}
}

func forwardRequest(w http.ResponseWriter, r *http.Request) int {
	callEngine := func(uriPath, uriQuery string, isError bool, errorStatusCode int) ([]byte, int) {
		var (
			statusCode int
			body       []byte
		)

		uri := r.URL
		uri.Host = r.Host
		if r.TLS != nil {
			uri.Scheme = "https"
		} else {
			uri.Scheme = "http"
		}

		if isError {
			body, statusCode = handleError(strings.HasPrefix(r.Header.Get("Accept"), "application/json"), uri.String(), errorStatusCode)
		} else {
			body, statusCode = handleRequest(strings.HasPrefix(r.Header.Get("Accept"), "application/json"), uri.String())
		}

		return body, statusCode
	}

	callEngineError := func(statusCode int) {
		activeTheme, err := database.GetActiveTheme()
		if err != nil {
			w.WriteHeader(statusCode)
			return
		}

		if activeTheme.ErrorBehavior == "homepage" {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		} else {
			callEngine("/internal/error-handler", "error="+strconv.Itoa(statusCode), true, statusCode)
			w.WriteHeader(statusCode)
		}
	}

	body, statusCode := callEngine(r.URL.Path, r.URL.RawQuery, false, 0)
	contentType := http.DetectContentType(body)

	if statusCode > 300 {
		callEngineError(http.StatusInternalServerError)
		return statusCode
	}

	bodyAsString := string(body)
	if strings.HasPrefix(contentType, "text/plain") {
		if strings.HasPrefix(bodyAsString, "<") {
			contentType = "text/html"
		} else if (strings.HasPrefix(bodyAsString, "[") && strings.HasSuffix(bodyAsString, "]")) || strings.HasPrefix(bodyAsString, "{") && strings.HasSuffix(bodyAsString, "}") {
			contentType = "application/json"
		}
	}

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)

	if _, err := w.Write(body); err != nil {
		slog.Error("Failed to write response from frontend engine", "error", err)
		callEngineError(http.StatusInternalServerError)
		return http.StatusInternalServerError
	}

	return statusCode
}
