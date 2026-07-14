package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"jinya.de/goji/api/middleware"
	"jinya.de/goji/database"
)

func getAllModernPages(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	pages, err := database.FindAllModernPages()
	if err != nil {
		sendError("Failed to find all modern pages", http.StatusInternalServerError, encoder, w, err)
		return
	}

	encoder.Encode(pages)
}

func getModernPageById(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Id is invalid", http.StatusBadRequest, encoder, w, err)
		return
	}

	page, err := database.FindModernPageById(id)
	if err != nil {
		sendError("Modern page not found", http.StatusNotFound, encoder, w, err)
		return
	}

	encoder.Encode(page)
}

func createModernPage(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	body, err := BodyFromRequest(r, "name")
	if err != nil {
		sendError("Invalid body", http.StatusBadRequest, encoder, w, err)
		return
	}

	artist := middleware.GetUserFromRequest(r)
	page := database.ModernPage{
		CreatorId:   artist.Id,
		UpdatedById: artist.Id,
	}
	if name, ok := GetFromBody[string](body, "name"); ok {
		page.Name = *name
	}

	err = page.Create()
	if err != nil {
		switch detectErrorType(err) {
		case errorTypeUnique:
			sendError("Name is not unique", http.StatusConflict, encoder, w, err)
			return
		default:
			sendError("Failed to create modern page", http.StatusInternalServerError, encoder, w, err)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	encoder.Encode(&page)
}

func updateModernPage(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Id is invalid", http.StatusBadRequest, encoder, w, err)
		return
	}

	page, err := database.FindModernPageById(id)
	if err != nil {
		sendError("Modern page not found", http.StatusNotFound, encoder, w, err)
		return
	}

	body, err := BodyFromRequest(r)
	if err != nil {
		sendError("Invalid body", http.StatusBadRequest, encoder, w, err)
		return
	}

	if name, ok := GetFromBody[string](body, "name"); ok {
		page.Name = *name
	}

	artist := middleware.GetUserFromRequest(r)
	page.UpdatedById = artist.Id

	err = page.Update()
	if err != nil {
		switch detectErrorType(err) {
		case errorTypeUnique:
			sendError("Name is not unique", http.StatusConflict, encoder, w, err)
			return
		default:
			sendError("Failed to create modern page", http.StatusInternalServerError, encoder, w, err)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteModernPage(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Id is invalid", http.StatusBadRequest, encoder, w, err)
		return
	}

	page, err := database.FindModernPageById(id)
	if err != nil {
		sendError("Modern page not found", http.StatusNotFound, encoder, w, err)
		return
	}

	err = page.Delete()
	if err != nil {
		sendError("Failed to delete modern page", http.StatusInternalServerError, encoder, w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getModernPageSections(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Id is invalid", http.StatusBadRequest, encoder, w, err)
		return
	}

	page, err := database.FindModernPageById(id)
	if err != nil {
		sendError("Modern page not found", http.StatusNotFound, encoder, w, err)
		return
	}

	sections, err := page.GetSections()
	if err != nil {
		sendError("Failed to load sections", http.StatusInternalServerError, encoder, w, err)
		return
	}

	encoder.Encode(sections)
}

func updateModernPageSections(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	decoder := json.NewDecoder(r.Body)
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Id is invalid", http.StatusBadRequest, encoder, w, err)
		return
	}

	page, err := database.FindModernPageById(id)
	if err != nil {
		sendError("Modern page not found", http.StatusNotFound, encoder, w, err)
		return
	}

	var body []struct {
		Gallery *int    `json:"gallery"`
		File    *int    `json:"file"`
		Link    string  `json:"link"`
		Html    *string `json:"html"`
		Type    string  `json:"type"`
	}

	err = decoder.Decode(&body)
	if err != nil {
		sendError("Invalid body", http.StatusBadRequest, encoder, w, err)
		return
	}

	sections := make([]database.ModernPageSection, len(body))

	for i, section := range body {
		if (section.Type == "gallery" && section.Gallery == nil) || (section.Type == "file" && section.File == nil) || (section.Type == "html" && section.Html == nil) {
			sendError("Invalid body", http.StatusBadRequest, encoder, w, err)
			return
		}

		s := database.ModernPageSection{
			Html: section.Html,
		}
		if section.Link != "" {
			s.Target = &section.Link
			s.Action = new("link")
		}

		if section.Gallery != nil {
			if _, err = database.FindGalleryById(*section.Gallery); err != nil {
				sendError(fmt.Sprintf("Gallery %d not found", *section.Gallery), http.StatusConflict, encoder, w, err)
				return
			}

			s.GalleryId = section.Gallery
		}
		if section.File != nil {
			if _, err = database.FindFileById(*section.File); err != nil {
				sendError(fmt.Sprintf("File %d not found", *section.File), http.StatusConflict, encoder, w, err)
				return
			}

			s.FileId = section.File
		}

		sections[i] = s
	}

	err = page.ReplaceSections(sections)
	if err != nil {
		sendError("Invalid body", http.StatusBadRequest, encoder, w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
