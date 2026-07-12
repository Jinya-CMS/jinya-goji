package database

import "encoding/json"

type ThemeAsset struct {
	ThemeId    int    `db:"theme_id,primarykey"`
	Name       string `db:"name,primarykey"`
	PublicPath string `db:"public_path"`
}

func (asset *ThemeAsset) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name       string `json:"name"`
		PublicPath string `json:"publicPath"`
	}{
		Name:       asset.Name,
		PublicPath: asset.PublicPath,
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
SELECT
    name,
    theme_id,
    public_path
FROM
    theme_asset
WHERE
    theme_id = ?
    AND name = ?
`, themeId, name)
}

func (asset *ThemeAsset) GetTheme() (*Theme, error) {
	return FindThemeById(asset.ThemeId)
}

func FindThemeAssetsByTheme(themeId int) ([]ThemeAsset, error) {
	return Select[ThemeAsset](`
SELECT
    name,
    theme_id,
    public_path
FROM
    theme_asset
WHERE
    theme_id = ?
`, themeId)
}
