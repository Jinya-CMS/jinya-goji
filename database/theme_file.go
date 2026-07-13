package database

import "encoding/json"

type ThemeFile struct {
	ThemeId int    `db:"theme_id,primarykey"`
	Name    string `db:"name,primarykey"`
	FileId  int    `db:"file_id"`
}

func (tf *ThemeFile) MarshalJSON() ([]byte, error) {
	file, err := tf.GetFile()
	if err != nil {
		return nil, err
	}

	return json.Marshal(struct {
		Name string `json:"name"`
		File *File  `json:"file"`
	}{
		Name: tf.Name,
		File: file,
	})
}

func (tf *ThemeFile) Create() error {
	return Insert(tf)
}

func (tf *ThemeFile) Update() error {
	return Update(tf)
}

func (tf *ThemeFile) Delete() error {
	return Delete(tf)
}

func FindThemeFileByThemeAndName(themeId int, name string) (*ThemeFile, error) {
	return SelectOne[ThemeFile](`
SELECT name,
       theme_id,
       file_id
FROM theme_file
WHERE theme_id = $1
  AND name = $2
`, themeId, name)
}

func (tf *ThemeFile) GetTheme() (*Theme, error) {
	return FindThemeById(tf.ThemeId)
}

func (tf *ThemeFile) GetFile() (*File, error) {
	return FindFileById(tf.FileId)
}

func FindThemeFilesByTheme(themeId int) ([]ThemeFile, error) {
	return Select[ThemeFile](`
select name,
       theme_id,
       file_id
from theme_file
where theme_id = $1
`, themeId)
}
