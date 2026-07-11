package database

import "time"

type BlogPost struct {
	Id            int       `db:"id,primarykey,autoincrement" json:"id"`
	Title         string    `db:"title" json:"title"`
	Slug          string    `db:"slug" json:"slug"`
	HeaderImageId *int      `db:"header_image_id" json:"-"`
	Public        bool      `db:"public" json:"public"`
	CreatedAt     time.Time `db:"created_at" json:"-"`
	CreatorId     int       `db:"creator_id" json:"-"`
	UpdatedById   int       `db:"updated_by_id" json:"-"`
	LastUpdatedAt time.Time `db:"last_updated_at" json:"-"`
	CategoryId    *int      `db:"category_id" json:"-"`
}
