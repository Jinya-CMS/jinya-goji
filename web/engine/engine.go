package engine

import (
	"context"
	_ "embed"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"jinya.de/goji/analytics"
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
	client, err := getEngineClient()
	if err != nil {
		slog.Error("Failed to get engine client", "error", err)
		http.Error(w, "Failed to handle your request", http.StatusInternalServerError)
		return 0
	}

	var body []byte
	var statusCode int

	if os.Getenv("ENV") == "dev" {
		engineReq := r.Clone(r.Context())
		engineReq.URL.Scheme = "http"
		engineReq.URL.Host = r.Host
		engineReq.RequestURI = ""

		res, err := client.Do(engineReq)
		if err != nil {
			slog.Error("Failed to send request to frontend engine", "error", err)
			http.Error(w, "Failed to handle your request", http.StatusInternalServerError)
			return http.StatusInternalServerError
		}
		defer res.Body.Close()

		for key, values := range res.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}

		body, err = io.ReadAll(res.Body)
		if err != nil {
			slog.Error("Failed to read body from frontend engine", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return http.StatusInternalServerError
		}

		statusCode = res.StatusCode
	} else {
		uri := r.URL
		uri.Host = r.Host
		if r.TLS != nil {
			uri.Scheme = "https"
		} else {
			uri.Scheme = "http"
		}
		body, statusCode = handleRequest(strings.HasPrefix(r.Header.Get("Accept"), "application/json"), uri.String())
	}

	contentType := http.DetectContentType(body)

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)

	if _, err = w.Write(body); err != nil {
		slog.Error("Failed to write response from frontend engine", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return http.StatusInternalServerError
	}

	return statusCode
}
