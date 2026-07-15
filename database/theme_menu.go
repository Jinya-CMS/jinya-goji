package database

import "encoding/json"

type ThemeMenu struct {
	ThemeId int    `db:"theme_id,primarykey"`
	Name    string `db:"name,primarykey"`
	MenuId  int    `db:"menu_id"`
}

func (tm *ThemeMenu) MarshalJSON() ([]byte, error) {
	menu, err := tm.GetMenu()
	if err != nil {
		return nil, err
	}

	return json.Marshal(struct {
		Name string `json:"name"`
		Menu *Menu  `json:"menu"`
	}{
		Name: tm.Name,
		Menu: menu,
	})
}

func (tm *ThemeMenu) Create() error {
	return Insert(tm)
}

func (tm *ThemeMenu) Update() error {
	return Update(tm)
}

func (tm *ThemeMenu) Delete() error {
	return Delete(tm)
}

func (tm *ThemeMenu) GetTheme() (*Theme, error) {
	return FindThemeById(tm.ThemeId)
}

func (tm *ThemeMenu) GetMenu() (*Menu, error) {
	return FindMenuById(tm.MenuId)
}

func FindThemeMenusByTheme(themeId int) ([]ThemeMenu, error) {
	return Select[ThemeMenu](`
select name,
       theme_id,
       menu_id
from theme_menu
where theme_id = $1
`, themeId)
}
