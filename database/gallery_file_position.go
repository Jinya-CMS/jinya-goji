package database

type GalleryFilePosition struct {
	Id        int `db:"id,primarykey,autoincrement" json:"id"`
	Position  int `db:"position" json:"position"`
	GalleryId int `db:"gallery_id" json:"-"`
	FileId    int `db:"file_id" json:"-"`
}
