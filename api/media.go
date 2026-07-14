package api

import (
	"encoding/json"
	"net/http"

	"jinya.de/goji/database"
)

func getMedia(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	files, err := database.FindRootFiles()
	if err != nil {
		sendError("Failed to collect root files", http.StatusInternalServerError, encoder, w, err)
		return
	}

	folders, err := database.FindRootFolders()
	if err != nil {
		sendError("Failed to collect root folders", http.StatusInternalServerError, encoder, w, err)
		return
	}

	tags, err := database.FindAllFileTags()
	if err != nil {
		sendError("Failed to fetch tags", http.StatusInternalServerError, encoder, w, err)
		return
	}

	encoder.Encode(map[string]any{
		"files":   files,
		"folders": folders,
		"tags":    tags,
	})
}
