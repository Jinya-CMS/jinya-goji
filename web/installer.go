package web

import (
	"html/template"
	"net/http"
	"strconv"

	"jinya.de/goji/auth"
	"jinya.de/goji/configuration"
	"jinya.de/goji/database"
)

func checkIfInstallIsAllowed(w http.ResponseWriter, r *http.Request) bool {
	users, err := database.FindAllArtists()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return false
	}

	if len(users) > 0 {
		http.Redirect(w, r, "/designer", http.StatusPermanentRedirect)
		return false
	}

	return true
}

func installer(w http.ResponseWriter, r *http.Request) {
	if !checkIfInstallIsAllowed(w, r) {
		return
	}

	tmpl, err := template.New("layout").ParseFS(tmplFs, "tmpl/installer.gohtml")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func installerMailer(w http.ResponseWriter, r *http.Request) {
	if !checkIfInstallIsAllowed(w, r) {
		return
	}

	err := r.ParseMultipartForm(10 * 1024 * 1024)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to parse form"))
		return
	}

	host := r.FormValue("host")
	port, err := strconv.Atoi(r.FormValue("port"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid port"))
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	encryption := r.FormValue("encryption")
	from := r.FormValue("from")

	err = configuration.SetString("mailer", "host", &host)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to set host"))
		return
	}

	err = configuration.SetInteger("mailer", "port", port)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to set port"))
		return
	}

	err = configuration.SetString("mailer", "username", &username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to set username"))
		return
	}

	err = configuration.SetString("mailer", "password", &password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to set password"))
		return
	}

	err = configuration.SetString("mailer", "encryption", &encryption)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to set encryption"))
		return
	}

	err = configuration.SetString("mailer", "from", &from)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to set from"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func installerFirstAdmin(w http.ResponseWriter, r *http.Request) {
	if !checkIfInstallIsAllowed(w, r) {
		return
	}

	artistName := r.FormValue("artistName")
	email := r.FormValue("email")
	password := r.FormValue("password")

	user := database.Artist{
		Email:      email,
		Enabled:    true,
		Roles:      []string{auth.RoleAdmin, auth.RoleWriter, auth.RoleReader},
		ArtistName: artistName,
	}

	err := user.SetPassword(password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to hash password"))
		return
	}

	err = user.Create()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to create first admin"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func installedCheck(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		users, err := database.FindAllArtists()
		if err == nil && len(users) > 0 {
			next.ServeHTTP(w, r)
			return
		}

		http.Redirect(w, r, "/installer", http.StatusTemporaryRedirect)
	})
}
