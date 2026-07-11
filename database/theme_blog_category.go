package database

type ThemeBlogCategory struct {
	ThemeId        int    `db:"theme_id,primarykey" json:"-"`
	Name           string `db:"name,primarykey" json:"name"`
	BlogCategoryId int    `db:"blog_category_id" json:"-"`
}
