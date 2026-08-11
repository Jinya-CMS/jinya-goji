package engine

import (
	"context"
	_ "embed"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"jinya.de/goji/analytics"
	"jinya.de/goji/database"
)

//go:embed target/release/engine
var engine []byte

const gojiEngineSocket = "/tmp/goji-frontend-engine.sock"

func StartEngine(ctx context.Context) {
	slog.Info("Preparing frontend engine")

	frontendEnginePath := ""
	if os.Getenv("ENV") != "dev" {
		slog.Debug("Store the embedded binary")
		tmpFile, err := os.CreateTemp("", "engine")
		if err != nil {
			panic(err)
		}

		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		defer tmpFile.Close()
		defer os.Remove(tmpFile.Name())

		_, err = tmpFile.Write(engine)
		if err != nil {
			panic(err)
		}

		err = tmpFile.Chmod(0o777)
		if err != nil {
			panic(err)
		}

		err = tmpFile.Close()
		if err != nil {
			panic(err)
		}

		frontendEnginePath = tmpFile.Name()

		slog.Info("Start frontend engine on socket")
		cmd := exec.CommandContext(ctx, frontendEnginePath, gojiEngineSocket)

		cmd.Env = os.Environ()
		cmd.Dir = wd
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Start()
		if err != nil {
			panic(err)
		}

		go func() {
			err := cmd.Wait()
			if err != nil {
				if ctx.Err() != nil {
					slog.Info("Frontend engine stopped because app is shutting down")
					return
				}

				panic(err)
			}

			slog.Info("Frontend engine stopped")
		}()
	}
}

func getEngineClient() (*http.Client, error) {
	httpc := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", gojiEngineSocket)
			},
		},
	}

	return &httpc, nil
}

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
