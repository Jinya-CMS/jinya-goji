package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"jinya.de/goji/api/middleware"
	"jinya.de/goji/database"
)

func getAllClassicPages(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	pages, err := database.FindAllClassicPages()
	if err != nil {
		sendError("Failed to find all classic pages", http.StatusInternalServerError, encoder, w, err)
		return
	}

	encoder.Encode(pages)
}

func getClassicPageById(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Id is invalid", http.StatusBadRequest, encoder, w, err)
		return
	}

	page, err := database.FindClassicPageById(id)
	if err != nil {
		sendError("Classic page not found", http.StatusNotFound, encoder, w, err)
		return
	}

	encoder.Encode(page)
}

func createClassicPage(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	body, err := BodyFromRequest(r, "title", "content")
	if err != nil {
		sendError("Invalid body", http.StatusBadRequest, encoder, w, err)
		return
	}

	artist := middleware.GetUserFromRequest(r)
	page := database.ClassicPage{
		CreatorId:   artist.Id,
		UpdatedById: artist.Id,
	}
	if title, ok := GetFromBody[string](body, "title"); ok {
		page.Title = *title
	}
	if content, ok := GetFromBody[string](body, "content"); ok {
		page.Content = *content
	}

	err = page.Create()
	if err != nil {
		switch detectErrorType(err) {
		case errorTypeUnique:
			sendError("Title is not unique", http.StatusConflict, encoder, w, err)
			return
		default:
			sendError("Failed to create classic page", http.StatusInternalServerError, encoder, w, err)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	encoder.Encode(&page)
}

func updateClassicPage(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Id is invalid", http.StatusBadRequest, encoder, w, err)
		return
	}

	page, err := database.FindClassicPageById(id)
	if err != nil {
		sendError("Classic page not found", http.StatusNotFound, encoder, w, err)
		return
	}

	body, err := BodyFromRequest(r)
	if err != nil {
		sendError("Invalid body", http.StatusBadRequest, encoder, w, err)
		return
	}

	if title, ok := GetFromBody[string](body, "title"); ok {
		page.Title = *title
	}
	if content, ok := GetFromBody[string](body, "content"); ok {
		page.Content = *content
	}

	artist := middleware.GetUserFromRequest(r)
	page.UpdatedById = artist.Id

	err = page.Update()
	if err != nil {
		switch detectErrorType(err) {
		case errorTypeUnique:
			sendError("Title is not unique", http.StatusConflict, encoder, w, err)
			return
		default:
			sendError("Failed to update classic page", http.StatusInternalServerError, encoder, w, err)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteClassicPage(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Id is invalid", http.StatusBadRequest, encoder, w, err)
		return
	}

	page, err := database.FindClassicPageById(id)
	if err != nil {
		sendError("Classic page not found", http.StatusNotFound, encoder, w, err)
		return
	}

	err = page.Delete()
	if err != nil {
		sendError("Failed to delete classic page", http.StatusInternalServerError, encoder, w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
