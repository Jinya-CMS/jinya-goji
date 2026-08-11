package api

import (
	"log/slog"
	"net/http"
	"strconv"

	"jinya.de/goji/database"
	"jinya.de/goji/libgoji"
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
		slog.Error("Failed to parse id", "error", err)
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
	preview, err := libgoji.PreviewClassicPage(activeTheme.Id, page)
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
		slog.Error("Failed to parse id", "error", err)
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

	preview, err := libgoji.PreviewModernPage(activeTheme.Id, id, sections)
	if err != nil {
		slog.Error("Failed to render preview", "id", id, "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte(preview))
}

func previewBlogPost(w http.ResponseWriter, r *http.Request) {
	slog.Info("Render preview for blog post")
	err := r.ParseForm()
	if err != nil {
		slog.Error("Failed to parse request form", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		slog.Error("Failed to parse id", "error", err)
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

	preview, err := libgoji.PreviewBlogPost(activeTheme.Id, id, sections)
	if err != nil {
		slog.Error("Failed to render preview", "id", id, "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte(preview))
}
