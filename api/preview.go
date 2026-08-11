package api

import (
	"log/slog"
	"net/http"
	"strconv"

	"jinya.de/goji/database"
	"jinya.de/goji/web/engine"
)

func previewClassicPage(w http.ResponseWriter, r *http.Request) {
	slog.Info("Render preview for classic page")
	err := r.ParseForm()
	if err != nil {
		slog.Error("Failed to parse request form", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		slog.Error("Failed to parse integer", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	content := r.FormValue("content")

	page, err := database.FindClassicPageById(id)
	if err != nil {
		slog.Error("Failed to find classic page by id", "id", id, "error", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	activeTheme, err := database.GetActiveTheme()
	if err != nil {
		slog.Error("Failed to get active theme", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	page.Content = content
	preview, err := engine.PreviewClassicPage(activeTheme.Id, page)
	if err != nil {
		slog.Error("Failed to render preview", "id", id, "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte(preview))
}

func previewModernPage(w http.ResponseWriter, r *http.Request) {
	slog.Info("Render preview for modern page")
	err := r.ParseForm()
	if err != nil {
		slog.Error("Failed to parse request form", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		slog.Error("Failed to parse integer", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sections := r.FormValue("sections")

	activeTheme, err := database.GetActiveTheme()
	if err != nil {
		slog.Error("Failed to get active theme", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	preview, err := engine.PreviewModernPage(activeTheme.Id, id, sections)
	if err != nil {
		slog.Error("Failed to render preview", "id", id, "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte(preview))
}
