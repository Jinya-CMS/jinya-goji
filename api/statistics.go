package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"jinya.de/goji/database"
)

func getGroupStatisticsByMonth(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoder := json.NewEncoder(w)

	group := vars["group"]
	months, err := strconv.Atoi(vars["months"])
	if err != nil {
		sendError("Invalid months format", http.StatusBadRequest, encoder, w, err)
		return
	}

	statistics, err := database.GetAnalyticsPastIntervalGrouped(group, &months, nil, nil, true)
	if err != nil {
		sendError("Failed to get analytics", http.StatusInternalServerError, encoder, w, err)
		return
	}

	encoder.Encode(statistics)
}

func getGroupStatisticsByYear(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoder := json.NewEncoder(w)

	group := vars["group"]
	years, err := strconv.Atoi(vars["years"])
	if err != nil {
		sendError("Invalid years format", http.StatusBadRequest, encoder, w, err)
		return
	}

	statistics, err := database.GetAnalyticsPastIntervalGrouped(group, new(years*12), nil, nil, true)
	if err != nil {
		sendError("Failed to get analytics", http.StatusInternalServerError, encoder, w, err)
		return
	}

	encoder.Encode(statistics)
}

func entityTypeStringToEntityType(entityType string) database.EntityType {
	switch entityType {
	case "blog-post":
		return database.EntityTypeBlogPost
	case "blog-category":
		return database.EntityTypeBlogCategory
	case "form":
		return database.EntityTypeForm
	case "gallery":
		return database.EntityTypeGallery
	case "modern-page":
		return database.EntityTypeModernPage
	case "classic-page":
		return database.EntityTypeClassicPage
	case "artist":
		return database.EntityTypeArtist
	}

	return -1
}

func getEntityStatisticsByMonth(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoder := json.NewEncoder(w)

	entityType := entityTypeStringToEntityType(vars["entityType"])
	if entityType == -1 {
		sendError("Invalid entity type", http.StatusBadRequest, encoder, w, fmt.Errorf("the entity type is invalid"))
		return
	}

	entityId, err := strconv.Atoi(vars["entityId"])
	if err != nil {
		sendError("Invalid entity id format", http.StatusBadRequest, encoder, w, err)
		return
	}

	months, err := strconv.Atoi(vars["months"])
	if err != nil {
		sendError("Invalid months format", http.StatusBadRequest, encoder, w, err)
		return
	}

	statistics, err := database.GetAnalyticsPastInterval(&months, &entityType, &entityId, true)
	if err != nil {
		sendError("Failed to get analytics", http.StatusInternalServerError, encoder, w, err)
		return
	}

	encoder.Encode(statistics)
}

func getEntityStatisticsByYear(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoder := json.NewEncoder(w)

	entityType := entityTypeStringToEntityType(vars["entityType"])
	if entityType == -1 {
		sendError("Invalid entity type", http.StatusBadRequest, encoder, w, fmt.Errorf("the entity type is invalid"))
		return
	}

	entityId, err := strconv.Atoi(vars["entityId"])
	if err != nil {
		sendError("Invalid entity id format", http.StatusBadRequest, encoder, w, err)
		return
	}

	years, err := strconv.Atoi(vars["years"])
	if err != nil {
		sendError("Invalid years format", http.StatusBadRequest, encoder, w, err)
		return
	}

	statistics, err := database.GetAnalyticsPastInterval(new(years*12), &entityType, &entityId, true)
	if err != nil {
		sendError("Failed to get analytics", http.StatusInternalServerError, encoder, w, err)
		return
	}

	encoder.Encode(statistics)
}

func getVisitStatisticsByMonth(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoder := json.NewEncoder(w)

	months, err := strconv.Atoi(vars["months"])
	if err != nil {
		sendError("Invalid months format", http.StatusBadRequest, encoder, w, err)
		return
	}

	statistics, err := database.GetTotalAnalyticsPastInterval(&months)
	if err != nil {
		sendError("Failed to get analytics", http.StatusInternalServerError, encoder, w, err)
		return
	}

	encoder.Encode(statistics)
}

func getVisitStatisticsByYear(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoder := json.NewEncoder(w)

	years, err := strconv.Atoi(vars["years"])
	if err != nil {
		sendError("Invalid years format", http.StatusBadRequest, encoder, w, err)
		return
	}

	statistics, err := database.GetTotalAnalyticsPastInterval(new(years * 12))
	if err != nil {
		sendError("Failed to get analytics", http.StatusInternalServerError, encoder, w, err)
		return
	}

	encoder.Encode(statistics)
}
