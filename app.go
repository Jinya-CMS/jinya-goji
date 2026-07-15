package main

import (
	"context"
	"embed"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"jinya.de/goji/api"
	"jinya.de/goji/database"
	"jinya.de/goji/localization"
	"jinya.de/goji/locate"
	"jinya.de/goji/logging"
	"jinya.de/goji/mailing"
	"jinya.de/goji/storage"
	"jinya.de/goji/theming"
	"jinya.de/goji/web"

	_ "github.com/joho/godotenv/autoload"
)

var (
	//go:embed static
	static embed.FS
	//go:embed compiler/target/release/compiler
	themeCompiler []byte
)

func main() {
	logging.SetupLogging()
	slog.Info("Prepared logging")

	slog.Info("Preparing the database")
	database.SetupDatabase()

	localization.SetupLocalization()

	slog.Info("Preparing mailing")
	mailing.SetupMailing()

	defer database.GetDbMap().Db.Close()

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	appCtx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	themeCompilerPath := ""
	if os.Getenv("ENV") != "dev" {
		slog.Info("Create theme compiler binary")
		tmpFile, err := os.CreateTemp("", "theme-compiler")
		if err != nil {
			panic(err)
		}

		defer tmpFile.Close()
		defer os.Remove(tmpFile.Name())

		_, err = tmpFile.Write(themeCompiler)
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

		themeCompilerPath = tmpFile.Name()

		slog.Info("Start theme compiler on socket")
		cmd := exec.CommandContext(appCtx, themeCompilerPath, "/var/run/goji-compiler.sock")

		cachePath := os.Getenv("CACHE_PATH")
		if cachePath == "" {
			path.Join(wd, "cache")
		}

		var rustEnvVars []string
		rustEnvVars = append(rustEnvVars, "RUST_LOG="+os.Getenv("LOG_LEVEL"))
		rustEnvVars = append(rustEnvVars, "STORAGE_PATH="+storage.GetStoragePath())
		rustEnvVars = append(rustEnvVars, "CACHE_PATH="+cachePath)
		rustEnvVars = append(rustEnvVars, "THEME_PATH="+theming.GetThemeBasePath())
		rustEnvVars = append(rustEnvVars, "DATABASE_URL="+os.Getenv("DATABASE_URL"))

		cmd.Env = rustEnvVars
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
				if appCtx.Err() != nil {
					slog.Info("Theme compiler stopped because app is shutting down")
					return
				}

				panic(err)
			}

			slog.Info("Theme compiler stopped")
		}()
	}

	go func() {
		slog.Info("Populate ip database")
		err := locate.PopulateDatabase()
		if err != nil {
			slog.Error("Failed to populate ip database", "error", err)
		} else {
			slog.Info("Populated ip database")
		}
	}()

	go func() {
		slog.Info("Trigger theme sync")
		err := theming.SyncThemes()
		if err != nil {
			slog.Error("Failed to sync themes", "error", err)
		} else {
			slog.Info("Theme sync successful")
		}
	}()

	slog.Info("Starting server")
	router := mux.NewRouter()

	api.SetupAdminApiRouter(router)
	web.SetupWebRouter(router)

	if os.Getenv("ENV") == "dev" {
		router.PathPrefix("/static/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Service-Worker-Allowed", "/")
			http.FileServerFS(os.DirFS(".")).ServeHTTP(w, r)
		})
	} else {
		router.PathPrefix("/static/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Service-Worker-Allowed", "/")
			http.FileServerFS(static).ServeHTTP(w, r)
		})
	}

	slog.Info("Serving at localhost:8090...")
	server := &http.Server{
		Addr:    ":8090",
		Handler: router,
	}

	go func() {
		<-appCtx.Done()

		slog.Info("Shutting down server")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			slog.Error("Failed to gracefully shut down server", "error", err)
		}
	}()

	err = server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}

	slog.Info("Server stopped")
}
