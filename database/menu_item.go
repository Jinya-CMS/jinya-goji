package database

type MenuItem struct {
	Id            int     `db:"id,primarykey,autoincrement" json:"id"`
	Position      int     `db:"position" json:"position"`
	MenuId        *int    `db:"menu_id" json:"-"`
	ParentId      *int    `db:"parent_id" json:"-"`
	Title         string  `db:"title" json:"title"`
	Route         *string `db:"route" json:"route"`
	Highlighted   bool    `db:"highlighted" json:"highlighted"`
	FormId        *int    `db:"form_id" json:"-"`
	ArtistId      *int    `db:"artist_id" json:"-"`
	GalleryId     *int    `db:"gallery_id" json:"-"`
	ModernPageId  *int    `db:"segment_page_id" json:"-"`
	ClassicPageId *int    `db:"page_id" json:"-"`
	CategoryId    *int    `db:"category_id" json:"-"`
	BlogHomePage  *bool   `db:"blog_home_page" json:"blogHomePage"`
}
