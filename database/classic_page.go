package database

import "time"

type ClassicPage struct {
	Id            int       `db:"id,primarykey,autoincrement" json:"id"`
	CreatorId     int       `db:"creator_id" json:"-"`
	UpdatedById   int       `db:"updated_by_id" json:"-"`
	CreatedAt     time.Time `db:"created_at" json:"-"`
	LastUpdatedAt time.Time `db:"last_updated_at" json:"-"`
	Content       string    `db:"content" json:"content"`
	Title         string    `db:"title" json:"title"`
}
