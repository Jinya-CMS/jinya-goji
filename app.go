package main

import (
	"context"
	"embed"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"jinya.de/goji/api"
	"jinya.de/goji/compiler"
	"jinya.de/goji/database"
	"jinya.de/goji/localization"
	"jinya.de/goji/locate"
	"jinya.de/goji/logging"
	"jinya.de/goji/mailing"
	"jinya.de/goji/native"
	"jinya.de/goji/theming"
	"jinya.de/goji/web"

	_ "github.com/joho/godotenv/autoload"
)

var (
	//go:embed static
	static embed.FS
)

func main() {
	defer native.Cleanup()

	logging.SetupLogging()
	slog.Info("Prepared logging")

	slog.Info("Prepare theme compiler")
	compiler.SetupCompiler()

	slog.Info("Prepare the database")
	database.SetupDatabase()

	slog.Info("Prepare localization")
	localization.SetupLocalization()

	slog.Info("Prepare mailing")
	mailing.SetupMailing()

	defer database.GetDbMap().Db.Close()

	appCtx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

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

	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}

	slog.Info("Server stopped")
}
