package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"jinya.de/goji/api/middleware"
	"jinya.de/goji/database"
	"jinya.de/goji/storage"
)

type fileBody struct {
	Name     string   `json:"name"`
	Tags     []string `json:"tags"`
	FolderId *int     `json:"folderId"`
}

func getFile(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Invalid id", http.StatusBadRequest, encoder, w)
		return
	}

	file, err := database.FindFileById(id)
	if err != nil {
		sendError("File not found", http.StatusNotFound, encoder, w)
		return
	}

	encoder.Encode(file)
}

func createFile(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	body, err := ParseBody[fileBody](r)
	if err != nil {
		sendError("Invalid body", http.StatusBadRequest, encoder, w)
		return
	}

	if body.FolderId != nil {
		_, err = database.FindFolderById(*body.FolderId)
		if err != nil {
			sendError("Parent folder not found", http.StatusConflict, encoder, w)
			return
		}
	}

	file := database.File{
		Name:     body.Name,
		FolderId: body.FolderId,
	}

	tags := make([]database.FileTag, len(body.Tags))
	for i, tag := range body.Tags {
		t, err := database.FindFileTagByName(tag)
		if err != nil {
			sendError(fmt.Sprintf("Tag %s not found", tag), http.StatusConflict, encoder, w)
			return
		}

		tags[i] = *t
	}

	artist := middleware.GetUserFromRequest(r)
	file.CreatorId = artist.Id
	file.UpdatedById = artist.Id

	err = file.Create(tags)
	if err != nil {
		errType := detectErrorType(err)
		switch errType {
		case errorTypeUnique:
			sendError("The name is not unique", http.StatusConflict, encoder, w)
			return
		default:
			sendError("Failed to create file", http.StatusInternalServerError, encoder, w)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	encoder.Encode(&file)
}

func updateFile(w http.ResponseWriter, r *http.Request) {
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

	file, err := database.FindFileById(id)
	if err != nil {
		sendError("File not found", http.StatusNotFound, encoder, w)
		return
	}

	if folderId, ok := GetIntFromBody(body, "folderId"); ok {
		_, err = database.FindFolderById(folderId)
		if err != nil {
			sendError("Parent folder not found", http.StatusConflict, encoder, w)
			return
		}
		file.FolderId = &folderId
	}

	if name, ok := GetFromBody[string](body, "name"); ok {
		file.Name = *name
	}

	artist := middleware.GetUserFromRequest(r)
	file.UpdatedById = artist.Id
	err = file.Update()
	if err != nil {
		errType := detectErrorType(err)
		switch errType {
		case errorTypeUnique:
			sendError("The name is not unique", http.StatusConflict, encoder, w)
			return
		default:
			sendError("Failed to update file", http.StatusInternalServerError, encoder, w)
			return
		}
	}

	if newTags, ok := GetArrayFromBody[string](body, "tags"); ok {
		tags := make([]database.FileTag, len(newTags))
		for i, tag := range newTags {
			t, err := database.FindFileTagByName(tag)
			if err != nil {
				sendError(fmt.Sprintf("Tag %s not found", tag), http.StatusConflict, encoder, w)
				return
			}

			tags[i] = *t
		}

		err = file.Tag(tags)
		if err != nil {
			sendError("Failed to tag file", http.StatusInternalServerError, encoder, w)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteFile(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Invalid id", http.StatusBadRequest, encoder, w)
		return
	}

	file, err := database.FindFileById(id)
	if err != nil {
		sendError("File not found", http.StatusNotFound, encoder, w)
		return
	}

	err = file.Delete()
	if err != nil {
		sendError("Failed to delete file", http.StatusInternalServerError, encoder, w)
		return
	}

	pathExists := database.FilePathExists(file.Path)
	if pathExists {
		_ = storage.DeleteFile(file.Path)
	}

	w.WriteHeader(http.StatusNoContent)
}

func startUpload(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Invalid id", http.StatusBadRequest, encoder, w)
		return
	}

	file, err := database.FindFileById(id)
	if err != nil {
		sendError("File not found", http.StatusNotFound, encoder, w)
		return
	}

	err = storage.StartFileUpload(file.Id)
	if err != nil {
		sendError("Failed to start upload", http.StatusInternalServerError, encoder, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func uploadChunk(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Invalid id", http.StatusBadRequest, encoder, w)
		return
	}

	position, err := strconv.Atoi(vars["position"])
	if err != nil {
		sendError("Invalid position", http.StatusBadRequest, encoder, w)
		return
	}

	file, err := database.FindFileById(id)
	if err != nil {
		sendError("File not found", http.StatusNotFound, encoder, w)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		sendError("Invalid chunk data", http.StatusBadRequest, encoder, w)
		return
	}

	err = storage.SaveFileChunk(file.Id, position, data)
	if err != nil {
		sendError("Failed to start upload", http.StatusInternalServerError, encoder, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func finishUpload(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Invalid id", http.StatusBadRequest, encoder, w)
		return
	}

	file, err := database.FindFileById(id)
	if err != nil {
		sendError("File not found", http.StatusNotFound, encoder, w)
		return
	}

	err = storage.FinishFileUpload(file.Id)
	if err != nil {
		sendError("Failed to start upload", http.StatusInternalServerError, encoder, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
