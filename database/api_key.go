package database

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"
)

const ApiKeyPrefix = "jinya-api-token"

type ApiKey struct {
	ApiKey        string    `db:"api_key,primarykey,unique"`
	UserId        int       `db:"user_id"`
	ValidSince    time.Time `db:"valid_since"`
	UserAgent     string    `db:"user_agent"`
	RemoteAddress string    `db:"remote_address"`
}

func NewApiKey(artist *Artist, r *http.Request) (string, *ApiKey, error) {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)

	apiKey := ApiKey{
		UserId:        artist.Id,
		ValidSince:    time.Now(),
		UserAgent:     r.UserAgent(),
		RemoteAddress: ip,
	}

	key, err := apiKey.SetApiKey()
	if err != nil {
		return "", nil, err
	}

	return key, &apiKey, nil
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

func FindApiKeysByArtist(artistId int) ([]ApiKey, error) {
	return Select[ApiKey](`
select api_key,
       user_id,
       valid_since,
       user_agent,
       remote_address
from api_key
where user_id = $1
`, artistId)
}

func (apiKey *ApiKey) SetApiKey() (string, error) {
	bytes := make([]byte, 20)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	plainKey := fmt.Sprintf("%s-%d-%s", ApiKeyPrefix, apiKey.UserId, hex.EncodeToString(bytes))
	hash := sha512.Sum512([]byte(plainKey))
	apiKey.ApiKey = hex.EncodeToString(hash[:])
	return plainKey, nil
}

func (apiKey *ApiKey) GetArtist() (*Artist, error) {
	return FindArtistById(apiKey.UserId)
}
