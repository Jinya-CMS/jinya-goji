package main

import (
	"embed"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"jinya.de/goji/api"
	"jinya.de/goji/database"
	"jinya.de/goji/localization"
	"jinya.de/goji/locate"
	"jinya.de/goji/logging"
	"jinya.de/goji/mailing"
	"jinya.de/goji/web"

	_ "github.com/joho/godotenv/autoload"
)

var (
	//go:embed static
	static embed.FS
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

	go func() {
		slog.Info("Populate ip database")
		err := locate.PopulateDatabase()
		if err != nil {
			slog.Error("Failed to populate ip database", "error", err)
		} else {
			slog.Info("Populated ip database")
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
	err := http.ListenAndServe(":8090", router)
	if err != nil {
		panic(err)
	}
}
