package database

type ThemeMenu struct {
	ThemeId int    `db:"theme_id,primarykey" json:"-"`
	Name    string `db:"name,primarykey" json:"name"`
	MenuId  int    `db:"menu_id" json:"-"`
}
