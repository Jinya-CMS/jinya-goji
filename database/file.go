package database

import "time"

type File struct {
	Id            int       `db:"id,primarykey,autoincrement" json:"id"`
	CreatorId     int       `db:"creator_id" json:"-"`
	UpdatedById   int       `db:"updated_by_id" json:"-"`
	CreatedAt     time.Time `db:"created_at" json:"-"`
	LastUpdatedAt time.Time `db:"last_updated_at" json:"-"`
	Path          string    `db:"path" json:"path"`
	Name          string    `db:"name" json:"name"`
	FolderId      *int      `db:"folder_id" json:"-"`
	Type          string    `db:"type" json:"type"`
}
