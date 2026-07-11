package database

type BlogPostSection struct {
	Id         int     `db:"id,primarykey,autoincrement" json:"id"`
	GalleryId  *int    `db:"gallery_id" json:"-"`
	FileId     *int    `db:"file_id" json:"-"`
	Html       *string `db:"html" json:"html,omitempty"`
	Link       *string `db:"link" json:"link,omitempty"`
	BlogPostId int     `db:"blog_post_id" json:"-"`
	Position   int     `db:"position" json:"position"`
}
