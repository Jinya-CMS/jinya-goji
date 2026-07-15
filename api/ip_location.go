package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"jinya.de/goji/locate"
)

func getIpLocation(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)
	ip := vars["ip"]
	location, err := locate.LocalizeIp(ip)
	if err != nil {
		sendError("Failed to locate ip", http.StatusInternalServerError, encoder, w, err)
		return
	}

	encoder.Encode(map[string]string{
		"city":      location.City,
		"country":   location.Country,
		"poweredBy": "https://db-ip.com",
	})
}
