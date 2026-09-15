package database

import "encoding/json"

type ThemeModernPage struct {
	ThemeId      int    `db:"theme_id,primarykey"`
	Name         string `db:"name,primarykey"`
	ModernPageId int    `db:"modern_page_id"`
}

func (tmp *ThemeModernPage) MarshalJSON() ([]byte, error) {
	modernPage, err := tmp.GetModernPage()
	if err != nil {
		return nil, err
	}

	return json.Marshal(struct {
		Name       string      `json:"name"`
		ModernPage *ModernPage `json:"modernPage"`
	}{
		Name:       tmp.Name,
		ModernPage: modernPage,
	})
}

func (tmp *ThemeModernPage) Create() error {
	return dbMap.Insert(tmp)
}

func (tmp *ThemeModernPage) Update() error {
	_, err := dbMap.Update(tmp)
	return err
}

func (tmp *ThemeModernPage) Delete() error {
	_, err := dbMap.Delete(tmp)
	return err
}

func (tmp *ThemeModernPage) GetTheme() (*Theme, error) {
	return FindThemeById(tmp.ThemeId)
}

func (tmp *ThemeModernPage) GetModernPage() (*ModernPage, error) {
	return FindModernPageById(tmp.ModernPageId)
}

func FindThemeModernPagesByTheme(themeId int) ([]ThemeModernPage, error) {
	return dbMap.SelectType[ThemeModernPage](`
select name,
       theme_id,
       modern_page_id
from theme_modern_page
where theme_id = $1
`, themeId)
}
