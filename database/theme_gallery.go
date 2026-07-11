package database

type ThemeGallery struct {
	ThemeId   int    `db:"theme_id,primarykey" json:"-"`
	Name      string `db:"name,primarykey" json:"name"`
	GalleryId int    `db:"gallery_id" json:"-"`
}
