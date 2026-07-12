package database

import "encoding/json"

type ThemeClassicPage struct {
	ThemeId       int    `db:"theme_id,primarykey"`
	Name          string `db:"name,primarykey"`
	ClassicPageId int    `db:"page_id"`
}

func (tcp *ThemeClassicPage) MarshalJSON() ([]byte, error) {
	classicPage, err := tcp.GetClassicPage()
	if err != nil {
		return nil, err
	}

	return json.Marshal(struct {
		Name        string       `json:"name"`
		ClassicPage *ClassicPage `json:"classicPage"`
	}{
		Name:        tcp.Name,
		ClassicPage: classicPage,
	})
}

func (tcp *ThemeClassicPage) Create() error {
	return Insert(tcp)
}

func (tcp *ThemeClassicPage) Update() error {
	return Update(tcp)
}

func (tcp *ThemeClassicPage) Delete() error {
	return Delete(tcp)
}

func FindThemeClassicPageByThemeAndName(themeId int, name string) (*ThemeClassicPage, error) {
	return SelectOne[ThemeClassicPage](`
SELECT
    name,
    theme_id,
    page_id
FROM
    theme_classic_page
WHERE
    theme_id = ?
    AND name = ?
`, themeId, name)
}

func (tcp *ThemeClassicPage) GetTheme() (*Theme, error) {
	return FindThemeById(tcp.ThemeId)
}

func (tcp *ThemeClassicPage) GetClassicPage() (*ClassicPage, error) {
	return FindClassicPageById(tcp.ClassicPageId)
}

func FindThemeClassicPagesByTheme(themeId int) ([]ThemeClassicPage, error) {
	return Select[ThemeClassicPage](`
SELECT
    name,
    theme_id,
    page_id
FROM
    theme_classic_page
WHERE
    theme_id = ?
`, themeId)
}
