package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"jinya.de/goji/api/middleware"
	"jinya.de/goji/database"
)

func getAllKnownDevices(w http.ResponseWriter, r *http.Request) {
	artist := middleware.GetUserFromRequest(r)
	encoder := json.NewEncoder(w)
	devices, err := database.FindKnownDevicesByUser(artist.Id)
	if err != nil {
		sendError("Failed to fetch known devices", http.StatusInternalServerError, encoder, w, err)
		return
	}

	encoder.Encode(devices)
}

func deleteKnownDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	encoder := json.NewEncoder(w)
	device, err := database.FindKnownDeviceByHashed(key)
	if err != nil {
		sendError("Known device not found", http.StatusNotFound, encoder, w, err)
		return
	}

	err = device.Delete()
	if err != nil {
		sendError("Known device not found", http.StatusNotFound, encoder, w, err)
	}

	w.WriteHeader(http.StatusNoContent)
}

func checkKnownDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	device, err := database.FindKnownDeviceByKey(key)
	if err != nil || device == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
