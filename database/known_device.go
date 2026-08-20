package database

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"jinya.de/goji/http_helper"
)

type KnownDevice struct {
	Id            int    `db:"id,primarykey,autoincrement"`
	UserId        int    `db:"user_id"`
	DeviceKey     string `db:"device_key"`
	UserAgent     string `db:"user_agent"`
	RemoteAddress string `db:"remote_address"`
}

func NewKnownDevice(artist *Artist, r *http.Request) (string, *KnownDevice, error) {
	ip := http_helper.GetClientIP(r)

	device := KnownDevice{
		Id:            0,
		UserId:        artist.Id,
		UserAgent:     r.UserAgent(),
		RemoteAddress: ip,
	}

	code, err := device.SetDeviceCode()
	if err != nil {
		return "", nil, err
	}

	return code, &device, err
}

func (device *KnownDevice) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		RemoteAddress string `json:"remoteAddress"`
		UserAgent     string `json:"userAgent"`
		Key           string `json:"key"`
	}{
		RemoteAddress: device.RemoteAddress,
		UserAgent:     device.UserAgent,
		Key:           device.DeviceKey,
	})
}

func (device *KnownDevice) SetDeviceCode() (string, error) {
	bytes := make([]byte, 20)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	plainKey := hex.EncodeToString(bytes)
	hash := sha512.Sum512([]byte(plainKey))
	device.DeviceKey = hex.EncodeToString(hash[:])
	return plainKey, nil
}

func (device *KnownDevice) Create() error {
	return Insert(device)
}

func (device *KnownDevice) Update() error {
	return Update(device)
}

func (device *KnownDevice) Delete() error {
	return Delete(device)
}

func FindKnownDevicesByUser(userId int) ([]KnownDevice, error) {
	return Select[KnownDevice](`
select id,
       user_id,
       device_key,
       user_agent,
       remote_address
from known_device
where user_id = $1
`, userId)
}

func FindKnownDeviceByKey(key string) (*KnownDevice, error) {
	hash := sha512.Sum512([]byte(key))
	hashedKey := hex.EncodeToString(hash[:])

	return SelectOne[KnownDevice](`
select id,
       user_id,
       device_key,
       user_agent,
       remote_address
from known_device
where device_key = $1
order by id
`, hashedKey)
}

func FindKnownDeviceByHashed(key string) (*KnownDevice, error) {
	return SelectOne[KnownDevice](`
select id,
       user_id,
       device_key,
       user_agent,
       remote_address
from known_device
where device_key = $1
order by id
`, key)
}
