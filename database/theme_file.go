package database

type ThemeFile struct {
	ThemeId int    `db:"theme_id,primarykey" json:"-"`
	Name    string `db:"name,primarykey" json:"name"`
	FileId  int    `db:"file_id" json:"-"`
}
