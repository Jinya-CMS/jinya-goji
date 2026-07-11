package database

type ThemeAsset struct {
	ThemeId    int    `db:"theme_id,primarykey" json:"-"`
	Name       string `db:"name,primarykey" json:"name"`
	PublicPath string `db:"public_path" json:"publicPath"`
}
