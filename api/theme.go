package api

import (
	"archive/zip"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"jinya.de/goji/database"
	"jinya.de/goji/theming"
)

func getAllThemes(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	themes, err := database.FindAllThemes()
	if err != nil {
		sendError("Failed to find themes", http.StatusInternalServerError, encoder, w, err)
		return
	}

	encoder.Encode(themes)
}

func getThemeById(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Invalid id format", http.StatusBadRequest, encoder, w, err)
		return
	}

	theme, err := database.FindThemeById(id)
	if err != nil {
		sendError("Theme not found", http.StatusNotFound, encoder, w, err)
		return
	}

	encoder.Encode(theme)
}

func storeTheme(name, path string, r *http.Request) error {
	tmpFile, err := os.CreateTemp("", name)
	if err != nil {
		return err
	}

	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	_, err = io.Copy(tmpFile, r.Body)
	if err != nil {
		return err
	}

	tmpFile.Close()

	fileName := tmpFile.Name()
	archive, err := zip.OpenReader(fileName)
	if err != nil {
		return err
	}

	defer archive.Close()

	for _, f := range archive.File {
		filePath := filepath.Join(path, f.Name)
		if !strings.HasPrefix(filePath, filepath.Clean(path)+string(os.PathSeparator)) {
			return err
		}
		if f.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		fileInArchive, err := f.Open()
		if err != nil {
			dstFile.Close()
			return err
		}

		if _, err = io.Copy(dstFile, fileInArchive); err != nil {
			dstFile.Close()
			fileInArchive.Close()
			return err
		}

		dstFile.Close()
		fileInArchive.Close()
	}

	return nil
}

func createTheme(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	name := uuid.NewString()

	err := storeTheme(name, theming.GetThemePath(name), r)
	if err != nil {
		sendError("Failed to create theme", http.StatusInternalServerError, encoder, w, err)
		return
	}

	err = theming.SyncThemes()
	if err != nil {
		sendError("Failed to create theme", http.StatusInternalServerError, encoder, w, err)
		return
	}

	theme, err := database.FindThemeByName(name)
	if err != nil {
		sendError("Failed to create theme", http.StatusInternalServerError, encoder, w, err)
		return
	}

	slog.Info("Created theme", "id", theme.Id, "name", theme.Name)
	encoder.Encode(theme)
}

func updateTheme(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Invalid id format", http.StatusBadRequest, encoder, w, err)
		return
	}

	theme, err := database.FindThemeById(id)
	if err != nil {
		sendError("Theme not found", http.StatusNotFound, encoder, w, err)
		return
	}

	err = storeTheme(theme.Name, theme.Path, r)
	if err != nil {
		sendError("Failed to update theme", http.StatusInternalServerError, encoder, w, err)
		return
	}

	err = theming.SyncThemes()
	if err != nil {
		sendError("Failed to update theme", http.StatusInternalServerError, encoder, w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteTheme(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Invalid id format", http.StatusBadRequest, encoder, w, err)
		return
	}

	theme, err := database.FindThemeById(id)
	if err != nil {
		sendError("Theme not found", http.StatusNotFound, encoder, w, err)
		return
	}

	err = theme.Delete()
	if err != nil {
		sendError("Failed to delete theme", http.StatusInternalServerError, encoder, w, err)
		return
	}

	slog.Info("Deleted theme", "id", theme.Id, "name", theme.Name)
	_ = theming.DeleteTheme(theme.Name)
}

func activateTheme(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Invalid id format", http.StatusBadRequest, encoder, w, err)
		return
	}

	theme, err := database.FindThemeById(id)
	if err != nil {
		sendError("Theme not found", http.StatusNotFound, encoder, w, err)
		return
	}

	err = theme.MakeActiveTheme()
	if err != nil {
		sendError("Failed to activate theme", http.StatusInternalServerError, encoder, w, err)
		return
	}

	slog.Info("Activated theme", "id", theme.Id, "name", theme.Name)
	w.WriteHeader(http.StatusNoContent)
}

func compileThemeAssets(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Invalid id format", http.StatusBadRequest, encoder, w, err)
		return
	}

	theme, err := database.FindThemeById(id)
	if err != nil {
		sendError("Theme not found", http.StatusNotFound, encoder, w, err)
		return
	}

	t, err := theming.NewTheme(theme)
	if err != nil {
		sendError("Failed to parse theme.yaml", http.StatusInternalServerError, encoder, w, err)
		return
	}

	err = t.CompileAssetCache()
	if err != nil {
		sendError("Failed to compile asset cache", http.StatusInternalServerError, encoder, w, err)
		return
	}

	slog.Info("Compiled theme assets", "id", theme.Id, "name", theme.Name)
	w.WriteHeader(http.StatusNoContent)
}

func updateThemeConfiguration(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	decoder := json.NewDecoder(r.Body)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Invalid id format", http.StatusBadRequest, encoder, w, err)
		return
	}

	theme, err := database.FindThemeById(id)
	if err != nil {
		sendError("Theme not found", http.StatusNotFound, encoder, w, err)
		return
	}

	var body struct {
		Configuration map[string]map[string]any `json:"configuration"`
	}
	err = decoder.Decode(&body)
	if err != nil {
		sendError("Invalid body", http.StatusBadRequest, encoder, w, err)
		return
	}

	for key, group := range body.Configuration {
		for itemKey, item := range group {
			if item == nil {
				delete(body.Configuration[key], itemKey)
			}
		}
	}

	theme.Configuration = body.Configuration
	err = theme.Update()
	if err != nil {
		sendError("Failed to update theme configuration", http.StatusInternalServerError, encoder, w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getThemeDefaultConfiguration(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Invalid id format", http.StatusBadRequest, encoder, w, err)
		return
	}

	theme, err := database.FindThemeById(id)
	if err != nil {
		sendError("Theme not found", http.StatusNotFound, encoder, w, err)
		return
	}

	t, err := theming.NewTheme(theme)
	if err != nil {
		sendError("Failed to parse theme.yaml", http.StatusInternalServerError, encoder, w, err)
		return
	}

	config := t.GetConfigurationValues()
	encoder.Encode(config)
}

func getThemeConfigurationStructure(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Invalid id format", http.StatusBadRequest, encoder, w, err)
		return
	}

	theme, err := database.FindThemeById(id)
	if err != nil {
		sendError("Theme not found", http.StatusNotFound, encoder, w, err)
		return
	}

	t, err := theming.NewTheme(theme)
	if err != nil {
		sendError("Failed to parse theme.yaml", http.StatusInternalServerError, encoder, w, err)
		return
	}

	config := t.GetConfigurationStructure()
	encoder.Encode(config)
}

func getThemePreview(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Invalid id format", http.StatusBadRequest, encoder, w, err)
		return
	}

	theme, err := database.FindThemeById(id)
	if err != nil {
		sendError("Theme not found", http.StatusNotFound, encoder, w, err)
		return
	}

	t, err := theming.NewTheme(theme)
	if err != nil {
		sendError("Failed to parse theme.yaml", http.StatusInternalServerError, encoder, w, err)
		return
	}

	http.ServeFile(w, r, t.GetPreviewImagePath())
}

func getThemeStyling(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Invalid id format", http.StatusBadRequest, encoder, w, err)
		return
	}

	theme, err := database.FindThemeById(id)
	if err != nil {
		sendError("Theme not found", http.StatusNotFound, encoder, w, err)
		return
	}

	t, err := theming.NewTheme(theme)
	if err != nil {
		sendError("Failed to parse theme.yaml", http.StatusInternalServerError, encoder, w, err)
		return
	}

	styling, err := t.GetStyleVariables()
	if err != nil {
		sendError("Failed to variables sass", http.StatusInternalServerError, encoder, w, err)
		return
	}

	encoder.Encode(styling)
}

func updateThemeStyling(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	decoder := json.NewDecoder(r.Body)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Invalid id format", http.StatusBadRequest, encoder, w, err)
		return
	}

	theme, err := database.FindThemeById(id)
	if err != nil {
		sendError("Theme not found", http.StatusNotFound, encoder, w, err)
		return
	}

	var body struct {
		Variables map[string]string `json:"variables"`
	}

	err = decoder.Decode(&body)
	if err != nil {
		sendError("Invalid body", http.StatusBadRequest, encoder, w, err)
		return
	}

	theme.ScssVariables = body.Variables
	err = theme.Update()
	if err != nil {
		sendError("Failed to save style variables", http.StatusInternalServerError, encoder, w, err)
		return
	}

	t, err := theming.NewTheme(theme)
	if err != nil {
		sendError("Failed to compile theme", http.StatusInternalServerError, encoder, w, err)
		return
	}

	err = t.CompileAssetCache()
	if err != nil {
		sendError("Failed to compile theme", http.StatusInternalServerError, encoder, w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getThemeBlogCategories(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Invalid id format", http.StatusBadRequest, encoder, w, err)
		return
	}

	theme, err := database.FindThemeById(id)
	if err != nil {
		sendError("Theme not found", http.StatusNotFound, encoder, w, err)
		return
	}

	categories, err := theme.GetCategories()
	if err != nil {
		sendError("Failed to get blog categories", http.StatusInternalServerError, encoder, w, err)
		return
	}

	result := map[string]*database.BlogCategory{}
	for _, cat := range categories {
		result[cat.Name], err = cat.GetBlogCategory()
		if err != nil {
			sendError("Failed to get blog categories", http.StatusInternalServerError, encoder, w, err)
			return
		}
	}

	encoder.Encode(result)
}

func updateThemeBlogCategory(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Invalid id format", http.StatusBadRequest, encoder, w, err)
		return
	}

	theme, err := database.FindThemeById(id)
	if err != nil {
		sendError("Theme not found", http.StatusNotFound, encoder, w, err)
		return
	}

	name := vars["name"]
	body, err := BodyFromRequest(r, "blogCategory")
	if err != nil {
		sendError("Invalid body", http.StatusBadRequest, encoder, w, err)
		return
	}

	blogCategory, ok := GetIntFromBody(body, "blogCategory")
	if !ok {
		sendError("Invalid body", http.StatusBadRequest, encoder, w, err)
		return
	}

	err = theme.SetBlogCategory(name, blogCategory)
	if err != nil {
		switch detectErrorType(err) {
		case errorTypeReference:
			sendError("Blog category does not exist", http.StatusConflict, encoder, w, err)
			return
		default:
			sendError("Failed to save theme link", http.StatusInternalServerError, encoder, w, err)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func getThemeClassicPages(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Invalid id format", http.StatusBadRequest, encoder, w, err)
		return
	}

	theme, err := database.FindThemeById(id)
	if err != nil {
		sendError("Theme not found", http.StatusNotFound, encoder, w, err)
		return
	}

	pages, err := theme.GetClassicPages()
	if err != nil {
		sendError("Failed to get classic pages", http.StatusInternalServerError, encoder, w, err)
		return
	}

	result := map[string]*database.ClassicPage{}
	for _, page := range pages {
		result[page.Name], err = page.GetClassicPage()
		if err != nil {
			sendError("Failed to get classic pages", http.StatusInternalServerError, encoder, w, err)
			return
		}
	}

	encoder.Encode(result)
}

func updateThemeClassicPage(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Invalid id format", http.StatusBadRequest, encoder, w, err)
		return
	}

	theme, err := database.FindThemeById(id)
	if err != nil {
		sendError("Theme not found", http.StatusNotFound, encoder, w, err)
		return
	}

	name := vars["name"]
	body, err := BodyFromRequest(r, "classicPage")
	if err != nil {
		sendError("Invalid body", http.StatusBadRequest, encoder, w, err)
		return
	}

	classicPage, ok := GetIntFromBody(body, "classicPage")
	if !ok {
		sendError("Invalid body", http.StatusBadRequest, encoder, w, err)
		return
	}

	err = theme.SetClassicPage(name, classicPage)
	if err != nil {
		switch detectErrorType(err) {
		case errorTypeReference:
			sendError("Classic page does not exist", http.StatusConflict, encoder, w, err)
			return
		default:
			sendError("Failed to save theme link", http.StatusInternalServerError, encoder, w, err)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func getThemeFiles(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Invalid id format", http.StatusBadRequest, encoder, w, err)
		return
	}

	theme, err := database.FindThemeById(id)
	if err != nil {
		sendError("Theme not found", http.StatusNotFound, encoder, w, err)
		return
	}

	files, err := theme.GetFiles()
	if err != nil {
		sendError("Failed to get files", http.StatusInternalServerError, encoder, w, err)
		return
	}

	result := map[string]*database.File{}
	for _, file := range files {
		result[file.Name], err = file.GetFile()
		if err != nil {
			sendError("Failed to get files", http.StatusInternalServerError, encoder, w, err)
			return
		}
	}

	encoder.Encode(result)
}

func updateThemeFile(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Invalid id format", http.StatusBadRequest, encoder, w, err)
		return
	}

	theme, err := database.FindThemeById(id)
	if err != nil {
		sendError("Theme not found", http.StatusNotFound, encoder, w, err)
		return
	}

	name := vars["name"]
	body, err := BodyFromRequest(r, "file")
	if err != nil {
		sendError("Invalid body", http.StatusBadRequest, encoder, w, err)
		return
	}

	file, ok := GetIntFromBody(body, "file")
	if !ok {
		sendError("Invalid body", http.StatusBadRequest, encoder, w, err)
		return
	}

	err = theme.SetFile(name, file)
	if err != nil {
		switch detectErrorType(err) {
		case errorTypeReference:
			sendError("File does not exist", http.StatusConflict, encoder, w, err)
			return
		default:
			sendError("Failed to save theme link", http.StatusInternalServerError, encoder, w, err)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func getThemeGalleries(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Invalid id format", http.StatusBadRequest, encoder, w, err)
		return
	}

	theme, err := database.FindThemeById(id)
	if err != nil {
		sendError("Theme not found", http.StatusNotFound, encoder, w, err)
		return
	}

	galleries, err := theme.GetGalleries()
	if err != nil {
		sendError("Failed to get galleries", http.StatusInternalServerError, encoder, w, err)
		return
	}

	result := map[string]*database.Gallery{}
	for _, gallery := range galleries {
		result[gallery.Name], err = gallery.GetGallery()
		if err != nil {
			sendError("Failed to get galleries", http.StatusInternalServerError, encoder, w, err)
			return
		}
	}

	encoder.Encode(result)
}

func updateThemeGallery(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Invalid id format", http.StatusBadRequest, encoder, w, err)
		return
	}

	theme, err := database.FindThemeById(id)
	if err != nil {
		sendError("Theme not found", http.StatusNotFound, encoder, w, err)
		return
	}

	name := vars["name"]
	body, err := BodyFromRequest(r, "gallery")
	if err != nil {
		sendError("Invalid body", http.StatusBadRequest, encoder, w, err)
		return
	}

	gallery, ok := GetIntFromBody(body, "gallery")
	if !ok {
		sendError("Invalid body", http.StatusBadRequest, encoder, w, err)
		return
	}

	err = theme.SetGallery(name, gallery)
	if err != nil {
		switch detectErrorType(err) {
		case errorTypeReference:
			sendError("Gallery does not exist", http.StatusConflict, encoder, w, err)
			return
		default:
			sendError("Failed to save theme link", http.StatusInternalServerError, encoder, w, err)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func getThemeMenus(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Invalid id format", http.StatusBadRequest, encoder, w, err)
		return
	}

	theme, err := database.FindThemeById(id)
	if err != nil {
		sendError("Theme not found", http.StatusNotFound, encoder, w, err)
		return
	}

	menus, err := theme.GetMenus()
	if err != nil {
		sendError("Failed to get menus", http.StatusInternalServerError, encoder, w, err)
		return
	}

	result := map[string]*database.Menu{}
	for _, menu := range menus {
		result[menu.Name], err = menu.GetMenu()
		if err != nil {
			sendError("Failed to get menus", http.StatusInternalServerError, encoder, w, err)
			return
		}
	}

	encoder.Encode(result)
}

func updateThemeMenu(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Invalid id format", http.StatusBadRequest, encoder, w, err)
		return
	}

	theme, err := database.FindThemeById(id)
	if err != nil {
		sendError("Theme not found", http.StatusNotFound, encoder, w, err)
		return
	}

	name := vars["name"]
	body, err := BodyFromRequest(r, "menu")
	if err != nil {
		sendError("Invalid body", http.StatusBadRequest, encoder, w, err)
		return
	}

	menu, ok := GetIntFromBody(body, "menu")
	if !ok {
		sendError("Invalid body", http.StatusBadRequest, encoder, w, err)
		return
	}

	err = theme.SetMenu(name, menu)
	if err != nil {
		switch detectErrorType(err) {
		case errorTypeReference:
			sendError("Menu does not exist", http.StatusConflict, encoder, w, err)
			return
		default:
			sendError("Failed to save theme link", http.StatusInternalServerError, encoder, w, err)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func getThemeModernPages(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Invalid id format", http.StatusBadRequest, encoder, w, err)
		return
	}

	theme, err := database.FindThemeById(id)
	if err != nil {
		sendError("Theme not found", http.StatusNotFound, encoder, w, err)
		return
	}

	pages, err := theme.GetModernPages()
	if err != nil {
		sendError("Failed to get modern pages", http.StatusInternalServerError, encoder, w, err)
		return
	}

	result := map[string]*database.ModernPage{}
	for _, page := range pages {
		result[page.Name], err = page.GetModernPage()
		if err != nil {
			sendError("Failed to get modern pages", http.StatusInternalServerError, encoder, w, err)
			return
		}
	}

	encoder.Encode(result)
}

func updateThemeModernPage(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Invalid id format", http.StatusBadRequest, encoder, w, err)
		return
	}

	theme, err := database.FindThemeById(id)
	if err != nil {
		sendError("Theme not found", http.StatusNotFound, encoder, w, err)
		return
	}

	name := vars["name"]
	body, err := BodyFromRequest(r, "modernPage")
	if err != nil {
		sendError("Invalid body", http.StatusBadRequest, encoder, w, err)
		return
	}

	modernPage, ok := GetIntFromBody(body, "modernPage")
	if !ok {
		sendError("Invalid body", http.StatusBadRequest, encoder, w, err)
		return
	}

	err = theme.SetModernPage(name, modernPage)
	if err != nil {
		switch detectErrorType(err) {
		case errorTypeReference:
			sendError("Modern page does not exist", http.StatusConflict, encoder, w, err)
			return
		default:
			sendError("Failed to save theme link", http.StatusInternalServerError, encoder, w, err)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
