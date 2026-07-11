package database

type KnownDevice struct {
	Id            int    `db:"id,primarykey,autoincrement" json:"-"`
	UserId        int    `db:"user_id" json:"-"`
	DeviceKey     string `db:"device_key" json:"key"`
	UserAgent     string `db:"user_agent" json:"userAgent"`
	RemoteAddress string `db:"remote_address" json:"remoteAddress"`
}
