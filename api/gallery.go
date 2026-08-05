package api

import (
	"encoding/json"
	"net/http"
	"slices"
	"strconv"

	"github.com/gorilla/mux"
	"jinya.de/goji/api/middleware"
	"jinya.de/goji/database"
)

func getAllGalleries(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	galleries, err := database.FindAllGalleries()
	if err != nil {
		sendError("Failed to find galleries", http.StatusInternalServerError, encoder, w, err)
		return
	}

	encoder.Encode(galleries)
}

func getGalleryById(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Id is invalid", http.StatusBadRequest, encoder, w, err)
		return
	}

	gallery, err := database.FindGalleryById(id)
	if err != nil {
		sendError("Gallery not found", http.StatusNotFound, encoder, w, err)
		return
	}

	encoder.Encode(gallery)
}

func createGallery(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	body, err := BodyFromRequest(r, "name")
	if err != nil {
		sendError("Body is invalid", http.StatusBadRequest, encoder, w, err)
		return
	}

	artist := middleware.GetUserFromRequest(r)
	gallery := database.Gallery{
		CreatorId:   artist.Id,
		UpdatedById: artist.Id,
	}

	if orientation, ok := GetFromBody[string](body, "orientation"); ok && *orientation == "vertical" {
		gallery.Orientation = "vertical"
	} else {
		gallery.Orientation = "horizontal"
	}
	if galleryType, ok := GetFromBody[string](body, "type"); ok && *galleryType == "masonry" {
		gallery.Type = "masonry"
	} else {
		gallery.Type = "sequence"
	}
	if description, ok := GetFromBody[string](body, "description"); ok {
		gallery.Description = *description
	}
	if group, ok := GetFromBody[string](body, "group"); ok {
		gallery.Group = *group
	}

	name, _ := GetFromBody[string](body, "name")
	gallery.Name = *name

	err = gallery.Create()
	if err != nil {
		errType := detectErrorType(err)
		switch errType {
		case errorTypeUnique:
			sendError("The name is not unique", http.StatusConflict, encoder, w, err)
			return
		default:
			sendError("Failed to create gallery", http.StatusInternalServerError, encoder, w, err)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	encoder.Encode(&gallery)
}

func updateGallery(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Id is invalid", http.StatusBadRequest, encoder, w, err)
		return
	}

	gallery, err := database.FindGalleryById(id)
	if err != nil {
		sendError("Gallery not found", http.StatusNotFound, encoder, w, err)
		return
	}

	body, err := BodyFromRequest(r)
	if err != nil {
		sendError("Body is invalid", http.StatusBadRequest, encoder, w, err)
		return
	}

	artist := middleware.GetUserFromRequest(r)

	gallery.UpdatedById = artist.Id

	if orientation, ok := GetFromBody[string](body, "orientation"); ok && slices.Contains([]string{"vertical", "horizontal"}, *orientation) {
		gallery.Orientation = *orientation
	}
	if galleryType, ok := GetFromBody[string](body, "type"); ok && slices.Contains([]string{"masonry", "sequence"}, *galleryType) {
		gallery.Type = *galleryType
	}
	if description, ok := GetFromBody[string](body, "description"); ok {
		gallery.Description = *description
	}
	if name, ok := GetFromBody[string](body, "name"); ok {
		gallery.Name = *name
	}
	if group, ok := GetFromBody[string](body, "group"); ok {
		gallery.Group = *group
	}

	err = gallery.Update()
	if err != nil {
		errType := detectErrorType(err)
		switch errType {
		case errorTypeUnique:
			sendError("The name is not unique", http.StatusConflict, encoder, w, err)
			return
		default:
			sendError("Failed to create gallery", http.StatusInternalServerError, encoder, w, err)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteGallery(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Id is invalid", http.StatusBadRequest, encoder, w, err)
		return
	}

	gallery, err := database.FindGalleryById(id)
	if err != nil {
		sendError("Gallery not found", http.StatusNotFound, encoder, w, err)
		return
	}

	err = gallery.Delete()
	if err != nil {
		sendError("Failed to delete gallery", http.StatusInternalServerError, encoder, w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getGalleryFiles(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Id is invalid", http.StatusBadRequest, encoder, w, err)
		return
	}

	gallery, err := database.FindGalleryById(id)
	if err != nil {
		sendError("Gallery not found", http.StatusNotFound, encoder, w, err)
		return
	}

	files, err := gallery.GetFiles()
	if err != nil {
		sendError("Gallery not found", http.StatusNotFound, encoder, w, err)
		return
	}

	encoder.Encode(files)
}

func addFileToGallery(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Id is invalid", http.StatusBadRequest, encoder, w, err)
		return
	}

	gallery, err := database.FindGalleryById(id)
	if err != nil {
		sendError("Gallery not found", http.StatusNotFound, encoder, w, err)
		return
	}

	body, err := BodyFromRequest(r, "file", "position")
	if err != nil {
		sendError("Invalid body", http.StatusBadRequest, encoder, w, err)
		return
	}

	galleryFilePosition := database.GalleryFilePosition{
		GalleryId: gallery.Id,
	}

	if file, ok := GetIntFromBody(body, "file"); ok {
		file, err := database.FindFileById(file)
		if err != nil {
			sendError("File not found", http.StatusNotFound, encoder, w, err)
			return
		}

		galleryFilePosition.FileId = file.Id
	} else {
		sendError("File must be a number", http.StatusBadRequest, encoder, w, err)
		return
	}

	if position, ok := GetIntFromBody(body, "position"); ok {
		galleryFilePosition.Position = position
	} else {
		sendError("Position must be a number", http.StatusBadRequest, encoder, w, err)
		return
	}

	err = galleryFilePosition.Create()
	if err != nil {
		sendError("Failed to add file to gallery", http.StatusInternalServerError, encoder, w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	encoder.Encode(&galleryFilePosition)
}

func updateGalleryFile(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Id is invalid", http.StatusBadRequest, encoder, w, err)
		return
	}

	position, err := strconv.Atoi(vars["position"])
	if err != nil {
		sendError("Position is invalid", http.StatusBadRequest, encoder, w, err)
		return
	}

	gallery, err := database.FindGalleryById(id)
	if err != nil {
		sendError("Gallery not found", http.StatusNotFound, encoder, w, err)
		return
	}

	galleryFilePosition, err := database.FindGalleryFilePositionByPosition(gallery.Id, position)
	if err != nil {
		sendError("Gallery file not found", http.StatusNotFound, encoder, w, err)
		return
	}

	body, err := BodyFromRequest(r)
	if err != nil {
		sendError("Body is invalid", http.StatusBadRequest, encoder, w, err)
		return
	}

	if position, ok := GetIntFromBody(body, "newPosition"); ok {
		galleryFilePosition.Position = position
	}
	if file, ok := GetIntFromBody(body, "file"); ok {
		galleryFilePosition.FileId = file
	}

	err = galleryFilePosition.Update()
	if err != nil {
		sendError("Failed to update file in gallery", http.StatusInternalServerError, encoder, w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteGalleryFile(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Id is invalid", http.StatusBadRequest, encoder, w, err)
		return
	}

	position, err := strconv.Atoi(vars["position"])
	if err != nil {
		sendError("Position is invalid", http.StatusBadRequest, encoder, w, err)
		return
	}

	gallery, err := database.FindGalleryById(id)
	if err != nil {
		sendError("Gallery not found", http.StatusNotFound, encoder, w, err)
		return
	}

	galleryFilePosition, err := database.FindGalleryFilePositionByPosition(gallery.Id, position)
	if err != nil {
		sendError("Gallery file not found", http.StatusNotFound, encoder, w, err)
		return
	}

	err = galleryFilePosition.Delete()
	if err != nil {
		sendError("Failed to delete file from gallery", http.StatusInternalServerError, encoder, w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
