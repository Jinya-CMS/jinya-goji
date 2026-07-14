package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"jinya.de/goji/database"
)

type fileTagBody struct {
	Name  string `json:"name"`
	Color string `json:"color"`
	Emoji string `json:"emoji"`
}

func createFileTag(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	body, err := ParseBody[fileTagBody](r)
	if err != nil {
		sendError("Invalid body", http.StatusBadRequest, encoder, w, err)
		return
	}

	fileTag := database.FileTag{
		Name:  body.Name,
		Color: body.Color,
		Emoji: body.Emoji,
	}
	err = fileTag.Create()
	if err != nil {
		errType := detectErrorType(err)
		switch errType {
		case errorTypeUnique:
			sendError("The name is not unique", http.StatusConflict, encoder, w, err)
			return
		default:
			sendError("Failed to create tag", http.StatusInternalServerError, encoder, w, err)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	encoder.Encode(&fileTag)
}

func updateFileTag(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Invalid id", http.StatusBadRequest, encoder, w, err)
		return
	}

	body, err := ParseBody[fileTagBody](r)
	if err != nil {
		sendError("Invalid body", http.StatusBadRequest, encoder, w, err)
		return
	}

	fileTag, err := database.FindFileTagById(id)
	if err != nil {
		sendError("File tag not found", http.StatusNotFound, encoder, w, err)
		return
	}

	fileTag.Name = body.Name
	fileTag.Color = body.Color
	fileTag.Emoji = body.Emoji

	err = fileTag.Update()
	if err != nil {
		errType := detectErrorType(err)
		switch errType {
		case errorTypeUnique:
			sendError("The name is not unique", http.StatusConflict, encoder, w, err)
			return
		default:
			sendError("Failed to update tag", http.StatusInternalServerError, encoder, w, err)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteFileTag(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError("Invalid id", http.StatusBadRequest, encoder, w, err)
		return
	}

	fileTag, err := database.FindFileTagById(id)
	if err != nil {
		sendError("Tag not found", http.StatusNotFound, encoder, w, err)
		return
	}

	err = fileTag.Delete()
	if err != nil {
		sendError("Failed to delete tag", http.StatusInternalServerError, encoder, w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
