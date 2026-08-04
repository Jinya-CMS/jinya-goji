package database

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

const ApiTokenPrefix = "jinya-api-persistent-token"

type ApiToken struct {
	Token      string     `db:"api_token,primarykey,unique"`
	Name       string     `db:"name"`
	UserId     int        `db:"user_id"`
	ExpiryDate *time.Time `db:"expiry_date"`
}

func NewApiToken(name string, expiryDate *time.Time, artist *Artist) (string, *ApiToken, error) {
	apiToken := ApiToken{
		UserId:     artist.Id,
		ExpiryDate: expiryDate,
		Name:       name,
	}

	key, err := apiToken.SetApiToken()
	if err != nil {
		return "", nil, err
	}

	return key, &apiToken, nil
}

func (apiToken *ApiToken) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ExpiryDate *time.Time `json:"expiryDate"`
		Token      string     `json:"token"`
		Name       string     `json:"name"`
	}{
		ExpiryDate: apiToken.ExpiryDate,
		Token:      apiToken.Token,
		Name:       apiToken.Name,
	})
}

func (apiToken *ApiToken) Create() error {
	return Insert(apiToken)
}

func (apiToken *ApiToken) Update() error {
	return Update(apiToken)
}

func (apiToken *ApiToken) Delete() error {
	return Delete(apiToken)
}

func FindApiTokenByHashedToken(key string) (*ApiToken, error) {
	return Get[ApiToken](key)
}

func FindApiTokenByToken(key string) (*ApiToken, error) {
	hash := sha512.Sum512([]byte(key))
	return FindApiTokenByHashedToken(hex.EncodeToString(hash[:]))
}

func FindApiTokensByArtist(artistId int) ([]ApiToken, error) {
	return Select[ApiToken](`
select api_token,
       user_id,
       name,
       expiry_date
from api_token
where user_id = $1
`, artistId)
}

func (apiToken *ApiToken) SetApiToken() (string, error) {
	bytes := make([]byte, 20)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	plainKey := fmt.Sprintf("%s-%d-%s", ApiTokenPrefix, apiToken.UserId, hex.EncodeToString(bytes))
	hash := sha512.Sum512([]byte(plainKey))
	apiToken.Token = hex.EncodeToString(hash[:])
	return plainKey, nil
}

func (apiToken *ApiToken) GetArtist() (*Artist, error) {
	return FindArtistById(apiToken.UserId)
}
