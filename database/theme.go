package database

import (
	"encoding/json"
	"log/slog"
	"strconv"
)

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
select id,
       name,
       display_name,
       description,
       scss_variables,
       configuration,
       has_api_theme
from theme
where name = $1
`, name)
}

func FindAllThemes() ([]Theme, error) {
	return Select[Theme](`
select id,
       name,
       display_name,
       description,
       scss_variables,
       configuration,
       has_api_theme
from theme
`)
}

func GetActiveTheme() (*Theme, error) {
	return SelectOne[Theme](`
select t.id,
       t.name,
       t.display_name,
       t.description,
       t.scss_variables,
       t.configuration,
       t.has_api_theme
from theme t
         join configuration c on c.current_frontend_theme_id = t.id
`)
}

func (theme *Theme) MakeActiveTheme() error {
	_, err := dbMap.Exec(`
update configuration
set current_frontend_theme_id = $1`, theme.Id)
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

func (theme *Theme) TriggerCompile() error {
	_, err := dbMap.Exec("select pg_notify($1, $2)", "compile_theme_"+strconv.Itoa(theme.Id), theme.Id)
	if err != nil {
		slog.Error("Failed to notify about registration")
	}

	return err
}
