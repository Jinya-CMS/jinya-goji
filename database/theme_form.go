package database

type ThemeForm struct {
	ThemeId int    `db:"theme_id,primarykey" json:"-"`
	Name    string `db:"name,primarykey" json:"name"`
	FormId  int    `db:"form_id" json:"-"`
}
