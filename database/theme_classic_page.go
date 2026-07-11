package database

type ThemeClassicPage struct {
	ThemeId       int    `db:"theme_id,primarykey" json:"-"`
	Name          string `db:"name,primarykey" json:"name"`
	ClassicPageId int    `db:"page_id" json:"-"`
}
