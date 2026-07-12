package database

import "encoding/json"

type KnownDevice struct {
	Id            int    `db:"id,primarykey,autoincrement"`
	UserId        int    `db:"user_id"`
	DeviceKey     string `db:"device_key"`
	UserAgent     string `db:"user_agent"`
	RemoteAddress string `db:"remote_address"`
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

func (device *KnownDevice) Create() error {
	return Insert(device)
}

func (device *KnownDevice) Update() error {
	return Update(device)
}

func (device *KnownDevice) Delete() error {
	return Delete(device)
}

func FindKnownDeviceById(id int) (*KnownDevice, error) {
	return Get[KnownDevice](id)
}

func FindKnownDevicesByUser(userId int) ([]KnownDevice, error) {
	return Select[KnownDevice](`
SELECT
    id,
    user_id,
    device_key,
    user_agent,
    remote_address
FROM
    known_device
WHERE
    user_id = ?
`, userId)
}

func FindKnownDeviceByKey(key string) (*KnownDevice, error) {
	return SelectOne[KnownDevice](`
SELECT
    id,
    user_id,
    device_key,
    user_agent,
    remote_address
FROM
    known_device
WHERE
    device_key = ?
`, key)
}
