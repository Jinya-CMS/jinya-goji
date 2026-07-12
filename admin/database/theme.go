package database

import "encoding/json"

type Theme struct {
	Id            int               `db:"id,primarykey,autoincrement"`
	Configuration map[string]any    `db:"configuration"`
	Description   map[string]string `db:"description"`
	Name          string            `db:"name"`
	DisplayName   string            `db:"display_name"`
	ScssVariables map[string]string `db:"scss_variables"`
	HasApiTheme   bool              `db:"has_api_theme"`
}

func (theme *Theme) MarshalJSON() ([]byte, error) {
	configuration := theme.Configuration
	if configuration == nil {
		configuration = make(map[string]any)
	}
	scssVariables := theme.ScssVariables
	if scssVariables == nil {
		scssVariables = make(map[string]string)
	}

	return json.Marshal(struct {
		Configuration map[string]any    `json:"configuration"`
		Description   map[string]string `json:"description"`
		Name          string            `json:"name"`
		DisplayName   string            `json:"displayName"`
		ScssVariables map[string]string `json:"scssVariables"`
		Id            int               `json:"id"`
		HasApi        bool              `json:"hasApi"`
	}{
		Configuration: configuration,
		Description:   theme.Description,
		Name:          theme.Name,
		DisplayName:   theme.DisplayName,
		ScssVariables: scssVariables,
		Id:            theme.Id,
		HasApi:        theme.HasApiTheme,
	})
}

func (theme *Theme) Create() error {
	return Insert(theme)
}

func (theme *Theme) Update() error {
	return Update(theme)
}

func (theme *Theme) Delete() error {
	return Delete(theme)
}

func FindThemeById(id int) (*Theme, error) {
	return Get[Theme](id)
}

func FindThemeByName(name string) (*Theme, error) {
	return SelectOne[Theme](`
SELECT
    id,
    name,
    display_name,
    description,
    scss_variables,
    configuration,
    has_api_theme
FROM
    theme
WHERE
    name = ?
`, name)
}

func FindAllThemes() ([]Theme, error) {
	return Select[Theme](`
SELECT
    id,
    name,
    display_name,
    description,
    scss_variables,
    configuration,
    has_api_theme
FROM
    theme
`)
}

func GetActiveTheme() (*Theme, error) {
	return SelectOne[Theme](`
SELECT
    t.id,
    t.name,
    t.display_name,
    t.description,
    t.scss_variables,
    t.configuration,
    t.has_api_theme
FROM
    theme t
    JOIN configuration c ON c.active_theme_id = t.id
`)
}

func (theme *Theme) MakeActiveTheme() error {
	_, err := dbMap.Exec("UPDATE configuration SET active_theme_id = ?", theme.Id)
	return err
}

func (theme *Theme) GetFiles() ([]ThemeFile, error) {
	return FindThemeFilesByTheme(theme.Id)
}

func (theme *Theme) GetCategories() ([]ThemeBlogCategory, error) {
	return FindThemeBlogCategoriesByTheme(theme.Id)
}

func (theme *Theme) GetGalleries() ([]ThemeGallery, error) {
	return FindThemeGalleriesByTheme(theme.Id)
}

func (theme *Theme) GetForms() ([]ThemeForm, error) {
	return FindThemeFormsByTheme(theme.Id)
}

func (theme *Theme) GetMenus() ([]ThemeMenu, error) {
	return FindThemeMenusByTheme(theme.Id)
}

func (theme *Theme) GetModernPages() ([]ThemeModernPage, error) {
	return FindThemeModernPagesByTheme(theme.Id)
}

func (theme *Theme) GetClassicPages() ([]ThemeClassicPage, error) {
	return FindThemeClassicPagesByTheme(theme.Id)
}

func (theme *Theme) GetAssets() ([]ThemeAsset, error) {
	return FindThemeAssetsByTheme(theme.Id)
}
