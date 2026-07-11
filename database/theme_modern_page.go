package database

type ThemeModernPage struct {
	ThemeId      int    `db:"theme_id,primarykey" json:"-"`
	Name         string `db:"name,primarykey" json:"name"`
	ModernPageId int    `db:"segment_page_id" json:"-"`
}
