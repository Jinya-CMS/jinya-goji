package main

import (
	"embed"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"jinya.de/goji/api"
	"jinya.de/goji/database"
	"jinya.de/goji/logging"
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

	defer database.GetDbMap().Db.Close()

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
