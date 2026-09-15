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

func (tbc *ThemeBlogCategory) Upsert() error {
	return dbMap.Insert(tbc)
}

func (tbc *ThemeBlogCategory) Update() error {
	_, err := dbMap.Update(tbc)
	return err
}

func (tbc *ThemeBlogCategory) Delete() error {
	_, err := dbMap.Delete(tbc)
	return err
}

func (tbc *ThemeBlogCategory) GetTheme() (*Theme, error) {
	return FindThemeById(tbc.ThemeId)
}

func (tbc *ThemeBlogCategory) GetBlogCategory() (*BlogCategory, error) {
	return FindBlogCategoryById(tbc.BlogCategoryId)
}

func FindThemeBlogCategoriesByTheme(themeId int) ([]ThemeBlogCategory, error) {
	return dbMap.SelectType[ThemeBlogCategory](`
select name,
       theme_id,
       blog_category_id
from theme_blog_category
where theme_id = $1
`, themeId)
}
