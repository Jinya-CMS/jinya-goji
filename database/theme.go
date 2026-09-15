package database

import (
	"encoding/json"

	"github.com/DerKnerd/gorp/v4"
)

type Theme struct {
	Id                 int                       `db:"id,primarykey,autoincrement"`
	ConfigurationPlain string                    `db:"configuration"`
	Configuration      map[string]map[string]any `db:"-"`
	DescriptionPlain   string                    `db:"description"`
	Description        map[string]string         `db:"-"`
	Name               string                    `db:"name"`
	DisplayName        string                    `db:"display_name"`
	ScssVariablesPlain string                    `db:"scss_variables"`
	ScssVariables      map[string]string         `db:"-"`
	ErrorBehavior      string                    `db:"error_behavior"`
	Path               string                    `db:"path"`
}

func (theme *Theme) PreInsert(_ gorp.SqlExecutor) error {
	config, err := json.Marshal(theme.Configuration)
	if err != nil {
		return err
	}

	description, err := json.Marshal(theme.Description)
	if err != nil {
		return err
	}

	scssVariables, err := json.Marshal(theme.ScssVariables)
	if err != nil {
		return err
	}

	theme.ConfigurationPlain = string(config)
	theme.DescriptionPlain = string(description)
	theme.ScssVariablesPlain = string(scssVariables)

	return nil
}

func (theme *Theme) PreUpdate(_ gorp.SqlExecutor) error {
	config, err := json.Marshal(theme.Configuration)
	if err != nil {
		return err
	}

	description, err := json.Marshal(theme.Description)
	if err != nil {
		return err
	}

	scssVariables, err := json.Marshal(theme.ScssVariables)
	if err != nil {
		return err
	}

	theme.ConfigurationPlain = string(config)
	theme.DescriptionPlain = string(description)
	theme.ScssVariablesPlain = string(scssVariables)

	return nil
}

func (theme *Theme) PostGet(_ gorp.SqlExecutor) error {
	err := json.Unmarshal([]byte(theme.ConfigurationPlain), &theme.Configuration)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(theme.DescriptionPlain), &theme.Description)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(theme.ScssVariablesPlain), &theme.ScssVariables)
	if err != nil {
		return err
	}

	return nil
}

func (theme *Theme) MarshalJSON() ([]byte, error) {
	configuration := theme.Configuration
	if configuration == nil {
		configuration = map[string]map[string]any{}
	}
	scssVariables := theme.ScssVariables
	if scssVariables == nil {
		scssVariables = map[string]string{}
	}

	return json.Marshal(struct {
		Configuration map[string]map[string]any `json:"configuration"`
		Description   map[string]string         `json:"description"`
		Name          string                    `json:"name"`
		DisplayName   string                    `json:"displayName"`
		ScssVariables map[string]string         `json:"scssVariables"`
		Id            int                       `json:"id"`
	}{
		Configuration: configuration,
		Description:   theme.Description,
		Name:          theme.Name,
		DisplayName:   theme.DisplayName,
		ScssVariables: scssVariables,
		Id:            theme.Id,
	})
}

func (theme *Theme) Create() error {
	return dbMap.Insert(theme)
}

func (theme *Theme) Update() error {
	_, err := dbMap.Update(theme)
	return err
}

func (theme *Theme) Delete() error {
	_, err := dbMap.Delete(theme)
	return err
}

func FindThemeById(id int) (*Theme, error) {
	return dbMap.GetType[Theme](id)
}

func FindThemeByName(name string) (*Theme, error) {
	return dbMap.SelectOneType[Theme](`
select id,
       name,
       display_name,
       description,
       scss_variables,
       configuration,
       error_behavior
from theme
where name = $1
`, name)
}

func FindAllThemes() ([]Theme, error) {
	return dbMap.SelectType[Theme](`
select id,
       name,
       display_name,
       description,
       scss_variables,
       configuration,
       error_behavior
from theme
order by id
`)
}

func GetActiveTheme() (*Theme, error) {
	return dbMap.SelectOneType[Theme](`
select t.id,
       t.name,
       t.display_name,
       t.description,
       t.scss_variables,
       t.configuration,
       t.error_behavior
from theme t
         join configuration c on c.current_frontend_theme_id = t.id
`)
}

func (theme *Theme) MakeActiveTheme() error {
	_, err := dbMap.SelectOneType[int](`
select current_frontend_theme_id
from configuration`)
	if err != nil {
		_, err = dbMap.Exec(`
insert into configuration (current_frontend_theme_id) 
values ($1)`, theme.Id)
	} else {
		_, err = dbMap.Exec(`
update configuration
set current_frontend_theme_id = $1`, theme.Id)
	}

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

func (theme *Theme) SetBlogCategory(name string, category int) error {
	_, err := dbMap.Exec(`
insert into theme_blog_category(name, theme_id, blog_category_id)
values ($1, $2, $3)
on conflict (name, theme_id) do update set blog_category_id=$3`, name, theme.Id, category)

	return err
}

func (theme *Theme) SetFile(name string, file int) error {
	_, err := dbMap.Exec(`
insert into theme_file(name, theme_id, file_id)
values ($1, $2, $3)
on conflict (name, theme_id) do update set file_id=$3`, name, theme.Id, file)

	return err
}

func (theme *Theme) SetMenu(name string, menu int) error {
	_, err := dbMap.Exec(`
insert into theme_menu(name, theme_id, menu_id)
values ($1, $2, $3)
on conflict (name, theme_id) do update set menu_id=$3`, name, theme.Id, menu)

	return err
}

func (theme *Theme) SetGallery(name string, gallery int) error {
	_, err := dbMap.Exec(`
insert into theme_gallery(name, theme_id, gallery_id)
values ($1, $2, $3)
on conflict (name, theme_id) do update set gallery_id=$3`, name, theme.Id, gallery)

	return err
}

func (theme *Theme) SetModernPage(name string, modernPage int) error {
	_, err := dbMap.Exec(`
insert into theme_modern_page(name, theme_id, modern_page_id)
values ($1, $2, $3)
on conflict (name, theme_id) do update set modern_page_id=$3`, name, theme.Id, modernPage)

	return err
}

func (theme *Theme) SetClassicPage(name string, classicPage int) error {
	_, err := dbMap.Exec(`
insert into theme_classic_page(name, theme_id, page_id)
values ($1, $2, $3)
on conflict (name, theme_id) do update set page_id=$3`, name, theme.Id, classicPage)

	return err
}
