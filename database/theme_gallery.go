package database

import "encoding/json"

type ThemeGallery struct {
	ThemeId   int    `db:"theme_id,primarykey"`
	Name      string `db:"name,primarykey"`
	GalleryId int    `db:"gallery_id"`
}

func (tg *ThemeGallery) MarshalJSON() ([]byte, error) {
	gallery, err := tg.GetGallery()
	if err != nil {
		return nil, err
	}

	return json.Marshal(struct {
		Name    string   `json:"name"`
		Gallery *Gallery `json:"gallery"`
	}{
		Name:    tg.Name,
		Gallery: gallery,
	})
}

func (tg *ThemeGallery) Create() error {
	return dbMap.Insert(tg)
}

func (tg *ThemeGallery) Update() error {
	_, err := dbMap.Update(tg)
	return err
}

func (tg *ThemeGallery) Delete() error {
	_, err := dbMap.Delete(tg)
	return err
}

func (tg *ThemeGallery) GetTheme() (*Theme, error) {
	return FindThemeById(tg.ThemeId)
}

func (tg *ThemeGallery) GetGallery() (*Gallery, error) {
	return FindGalleryById(tg.GalleryId)
}

func FindThemeGalleriesByTheme(themeId int) ([]ThemeGallery, error) {
	return dbMap.SelectType[ThemeGallery](`
select name,
       theme_id,
       gallery_id
from theme_gallery
where theme_id = $1
`, themeId)
}
