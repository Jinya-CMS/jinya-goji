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
	return Insert(asset)
}

func (asset *ThemeAsset) Update() error {
	return Update(asset)
}

func (asset *ThemeAsset) Delete() error {
	return Delete(asset)
}

func FindThemeAssetByThemeAndName(themeId int, name string) (*ThemeAsset, error) {
	return SelectOne[ThemeAsset](`
select name,
       theme_id,
       path
from theme_asset
where theme_id = $1
  and name = $2
`, themeId, name)
}

func (asset *ThemeAsset) GetTheme() (*Theme, error) {
	return FindThemeById(asset.ThemeId)
}

func FindThemeAssetsByTheme(themeId int) ([]ThemeAsset, error) {
	return Select[ThemeAsset](`
select name,
       theme_id,
       path
from theme_asset
where theme_id = $1
`, themeId)
}
