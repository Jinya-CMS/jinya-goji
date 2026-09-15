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
	return dbMap.Insert(tf)
}

func (tf *ThemeFile) Update() error {
	_, err := dbMap.Update(tf)
	return err
}

func (tf *ThemeFile) Delete() error {
	_, err := dbMap.Delete(tf)
	return err
}

func (tf *ThemeFile) GetTheme() (*Theme, error) {
	return FindThemeById(tf.ThemeId)
}

func (tf *ThemeFile) GetFile() (*File, error) {
	return FindFileById(tf.FileId)
}

func FindThemeFilesByTheme(themeId int) ([]ThemeFile, error) {
	return dbMap.SelectType[ThemeFile](`
select name,
       theme_id,
       file_id
from theme_file
where theme_id = $1
`, themeId)
}
