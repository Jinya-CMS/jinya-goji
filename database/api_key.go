package database

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

type ApiKey struct {
	ApiKey        string    `db:"api_key,primarykey,unique"`
	UserId        int       `db:"user_id"`
	ValidSince    time.Time `db:"valid_since"`
	UserAgent     string    `db:"user_agent"`
	RemoteAddress string    `db:"remote_address"`
}

func (apiKey *ApiKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		RemoteAddress string `json:"remoteAddress"`
		ValidSince    string `json:"validSince"`
		UserAgent     string `json:"userAgent"`
		Key           string `json:"key"`
	}{
		RemoteAddress: apiKey.RemoteAddress,
		ValidSince:    apiKey.ValidSince.Format(time.RFC3339),
		UserAgent:     apiKey.UserAgent,
		Key:           apiKey.ApiKey,
	})
}

func (apiKey *ApiKey) Create() error {
	return Insert(apiKey)
}

func (apiKey *ApiKey) Update() error {
	return Update(apiKey)
}

func (apiKey *ApiKey) Delete() error {
	return Delete(apiKey)
}

func FindApiKeyByHashedKey(key string) (*ApiKey, error) {
	return Get[ApiKey](key)
}

func FindApiKeyByKey(key string) (*ApiKey, error) {
	hash := sha512.Sum512([]byte(key))
	return FindApiKeyByHashedKey(hex.EncodeToString(hash[:]))
}

func FindApiKeyByArtist(artistId int) ([]ApiKey, error) {
	return Select[ApiKey](`
SELECT
    api_key,
    user_id,
    valid_since,
    user_agent,
    remote_address
FROM
    api_key
WHERE
    user_id = ?
`, artistId)
}

func (apiKey *ApiKey) SetApiKey() (string, error) {
	bytes := make([]byte, 20)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	plainKey := fmt.Sprintf("jinya-api-token-%d-%s", apiKey.UserId, hex.EncodeToString(bytes))
	hash := sha512.Sum512([]byte(plainKey))
	apiKey.ApiKey = hex.EncodeToString(hash[:])
	return plainKey, nil
}

func (apiKey *ApiKey) GetArtist() (*Artist, error) {
	return FindArtistById(apiKey.UserId)
}
