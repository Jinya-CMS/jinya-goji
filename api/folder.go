package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"jinya.de/goji/database"
)

type folderBody struct {
	Name     string `json:"name"`
	ParentId *int   `json:"parentId"`
}

func createFolder(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	body, err := ParseBody[folderBody](r)
	if err != nil {
		sendError("Invalid body", http.StatusBadRequest, encoder, w)
		return
	}

	folder := database.Folder{
		Name:     body.Name,
		ParentId: body.ParentId,
	}
	err = folder.Create()
	if err != nil {
		errType := detectErrorType(err)
		switch errType {
		case errorTypeUnique:
			sendError("The name is not unique", http.StatusConflict, encoder, w)
			return
		case errorTypeReference:
			sendError("The parent folder does not exist", http.StatusConflict, encoder, w)
			return
		default:
			sendError("Failed to create folder", http.StatusInternalServerError, encoder, w)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	encoder.Encode(&folder)
}

func updateFolder(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Invalid id", http.StatusBadRequest, encoder, w)
		return
	}

	body, err := BodyFromRequest(r)
	if err != nil {
		sendError("Invalid body", http.StatusBadRequest, encoder, w)
		return
	}

	folder, err := database.FindFolderById(id)
	if err != nil {
		sendError("Folder not found", http.StatusNotFound, encoder, w)
		return
	}

	if name, ok := GetFromBody[string](body, "name"); ok {
		folder.Name = *name
	}
	if parentId, ok := GetIntFromBody(body, "parentId"); ok {
		folder.ParentId = &parentId
	}

	err = folder.Update()
	if err != nil {
		errType := detectErrorType(err)
		switch errType {
		case errorTypeUnique:
			sendError("The name is not unique", http.StatusConflict, encoder, w)
			return
		case errorTypeReference:
			sendError("The parent folder does not exist", http.StatusConflict, encoder, w)
			return
		default:
			sendError("Failed to update folder", http.StatusInternalServerError, encoder, w)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteFolder(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Invalid id", http.StatusBadRequest, encoder, w)
		return
	}

	folder, err := database.FindFolderById(id)
	if err != nil {
		sendError("Folder not found", http.StatusNotFound, encoder, w)
		return
	}

	err = folder.Delete()
	if err != nil {
		sendError("Failed to delete folder", http.StatusInternalServerError, encoder, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
