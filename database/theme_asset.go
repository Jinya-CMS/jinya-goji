package database

import "encoding/json"

type ThemeAsset struct {
	ThemeId int    `db:"theme_id,primarykey"`
	Name    string `db:"name,primarykey"`
	Path    string `db:"path"`
}

func (asset *ThemeAsset) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name string `json:"name"`
	}{
		Name: asset.Name,
	})
}

func (asset *ThemeAsset) Create() error {
	return dbMap.Insert(asset)
}

func (asset *ThemeAsset) Update() error {
	_, err := dbMap.Update(asset)
	return err
}

func (asset *ThemeAsset) Delete() error {
	_, err := dbMap.Delete(asset)
	return err
}

func FindThemeAssetByThemeAndName(themeName, name string) (*ThemeAsset, error) {
	return dbMap.SelectOneType[ThemeAsset](`
select ta.name,
       ta.theme_id,
       ta.path
from theme_asset ta
         join theme t on t.id = ta.theme_id
where t.name = $1
  and ta.name = $2
`, themeName, name)
}

func (asset *ThemeAsset) GetTheme() (*Theme, error) {
	return FindThemeById(asset.ThemeId)
}

func FindThemeAssetsByTheme(themeId int) ([]ThemeAsset, error) {
	return dbMap.SelectType[ThemeAsset](`
select name,
       theme_id,
       path
from theme_asset
where theme_id = $1
`, themeId)
}
