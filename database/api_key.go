package database

import "time"

type ApiKey struct {
	ApiKey        string    `db:"api_key,primarykey" json:"key"`
	UserId        int       `db:"user_id" json:"-"`
	ValidSince    time.Time `db:"valid_since" json:"validSince"`
	UserAgent     string    `db:"user_agent" json:"userAgent"`
	RemoteAddress string    `db:"remote_address" json:"remoteAddress"`
}
