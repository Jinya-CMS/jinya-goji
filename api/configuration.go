package api

import (
	"encoding/json"
	"net/http"

	"jinya.de/goji/configuration"
)

func getConfiguration(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	encoder.Encode(map[string]any{
		"mailer": map[string]any{
			"encryption": configuration.GetString("mailer", "encryption", ""),
			"from":       configuration.GetString("mailer", "from", ""),
			"host":       configuration.GetString("mailer", "host", ""),
			"password":   configuration.GetString("mailer", "password", ""),
			"port":       configuration.GetInteger("mailer", "port", 25),
			"username":   configuration.GetString("mailer", "username", ""),
		},
		"jinya": map[string]any{
			"api_key_expiry":  configuration.GetInteger("jinya", "api_key_expiry", 86400),
			"ip_database_url": configuration.GetString("jinya", "ip_database_url", "https://download.db-ip.com/free/dbip-city-lite-{YEAR}-{MONTH}.csv.gz"),
		},
	})
}

func updateConfiguration(w http.ResponseWriter, r *http.Request) {
	type config struct {
		Mailer struct {
			Encryption *string `json:"encryption"`
			From       string  `json:"from"`
			Host       string  `json:"host"`
			Password   string  `json:"password"`
			Port       int     `json:"port"`
			Username   string  `json:"username"`
		} `json:"mailer"`
		Jinya struct {
			ApiKeyExpiry  int    `json:"api_key_expiry"`
			IpDatabaseUrl string `json:"ip_database_url"`
		} `json:"jinya"`
	}

	encoder := json.NewEncoder(w)
	decoder := json.NewDecoder(r.Body)

	var body config
	err := decoder.Decode(&body)
	if err != nil {
		sendError("Invalid body", http.StatusBadRequest, encoder, w, err)
		return
	}

	_ = configuration.SetInteger("mailer", "port", body.Mailer.Port)
	_ = configuration.SetString("mailer", "encryption", body.Mailer.Encryption)
	_ = configuration.SetString("mailer", "from", &body.Mailer.From)
	_ = configuration.SetString("mailer", "host", &body.Mailer.Host)
	_ = configuration.SetString("mailer", "password", &body.Mailer.Password)
	_ = configuration.SetString("mailer", "username", &body.Mailer.Username)
	_ = configuration.SetInteger("jinya", "api_key_expiry", body.Jinya.ApiKeyExpiry)
	_ = configuration.SetString("jinya", "ip_database_url", &body.Jinya.IpDatabaseUrl)

	w.WriteHeader(http.StatusNoContent)
}
