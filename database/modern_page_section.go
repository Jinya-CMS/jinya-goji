package database

type ModernPageSection struct {
	Id        int     `db:"id,primarykey,autoincrement" json:"id"`
	Position  int     `db:"position" json:"position"`
	GalleryId *int    `db:"gallery_id" json:"-"`
	FileId    *int    `db:"file_id" json:"-"`
	Html      *string `db:"html" json:"html,omitempty"`
	PageId    int     `db:"page_id" json:"-"`
	Action    *string `db:"action" json:"-"`
	Script    *string `db:"script" json:"script,omitempty"`
	Target    *string `db:"target" json:"-"`
}
