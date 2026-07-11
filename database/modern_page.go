package database

import "time"

type ModernPage struct {
	Id            int       `db:"id,primarykey,autoincrement" json:"id"`
	CreatorId     int       `db:"creator_id" json:"-"`
	UpdatedById   int       `db:"updated_by_id" json:"-"`
	CreatedAt     time.Time `db:"created_at" json:"-"`
	LastUpdatedAt time.Time `db:"last_updated_at" json:"-"`
	Name          string    `db:"name" json:"name"`
}
