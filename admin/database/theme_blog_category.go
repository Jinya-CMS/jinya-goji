package database

import "encoding/json"

type ThemeBlogCategory struct {
	ThemeId        int    `db:"theme_id,primarykey"`
	Name           string `db:"name,primarykey"`
	BlogCategoryId int    `db:"blog_category_id"`
}

func (tbc *ThemeBlogCategory) MarshalJSON() ([]byte, error) {
	category, _ := tbc.GetBlogCategory()

	return json.Marshal(struct {
		Name         string        `json:"name"`
		BlogCategory *BlogCategory `json:"blogCategory"`
	}{
		Name:         tbc.Name,
		BlogCategory: category,
	})
}

func (tbc *ThemeBlogCategory) Create() error {
	return Insert(tbc)
}

func (tbc *ThemeBlogCategory) Update() error {
	return Update(tbc)
}

func (tbc *ThemeBlogCategory) Delete() error {
	return Delete(tbc)
}

func FindThemeBlogCategoryByThemeAndName(themeId int, name string) (*ThemeBlogCategory, error) {
	return SelectOne[ThemeBlogCategory](`
SELECT
    name,
    theme_id,
    blog_category_id
FROM
    theme_blog_category
WHERE
    theme_id = ?
    AND name = ?
`, themeId, name)
}

func (tbc *ThemeBlogCategory) GetTheme() (*Theme, error) {
	return FindThemeById(tbc.ThemeId)
}

func (tbc *ThemeBlogCategory) GetBlogCategory() (*BlogCategory, error) {
	return FindBlogCategoryById(tbc.BlogCategoryId)
}

func FindThemeBlogCategoriesByTheme(themeId int) ([]ThemeBlogCategory, error) {
	return Select[ThemeBlogCategory](`
SELECT
    name,
    theme_id,
    blog_category_id
FROM
    theme_blog_category
WHERE
    theme_id = ?
`, themeId)
}
