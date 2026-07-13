package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"jinya.de/goji/database"
)

func getAllKnownDevices(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func deleteKnownDevice(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
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
