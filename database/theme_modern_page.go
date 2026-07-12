package database

import "encoding/json"

type ThemeModernPage struct {
	ThemeId      int    `db:"theme_id,primarykey"`
	Name         string `db:"name,primarykey"`
	ModernPageId int    `db:"segment_page_id"`
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
	return Insert(tmp)
}

func (tmp *ThemeModernPage) Update() error {
	return Update(tmp)
}

func (tmp *ThemeModernPage) Delete() error {
	return Delete(tmp)
}

func FindThemeModernPageByThemeAndName(themeId int, name string) (*ThemeModernPage, error) {
	return SelectOne[ThemeModernPage](`
SELECT
    name,
    theme_id,
    segment_page_id
FROM
    theme_segment_page
WHERE
    theme_id = ?
    AND name = ?
`, themeId, name)
}

func (tmp *ThemeModernPage) GetTheme() (*Theme, error) {
	return FindThemeById(tmp.ThemeId)
}

func (tmp *ThemeModernPage) GetModernPage() (*ModernPage, error) {
	return FindModernPageById(tmp.ModernPageId)
}

func FindThemeModernPagesByTheme(themeId int) ([]ThemeModernPage, error) {
	return Select[ThemeModernPage](`
SELECT
    name,
    theme_id,
    segment_page_id
FROM
    theme_segment_page
WHERE
    theme_id = ?
`, themeId)
}
