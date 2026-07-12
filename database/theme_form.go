package database

import "encoding/json"

type ThemeForm struct {
	ThemeId int    `db:"theme_id,primarykey"`
	Name    string `db:"name,primarykey"`
	FormId  int    `db:"form_id"`
}

func (tf *ThemeForm) MarshalJSON() ([]byte, error) {
	form, err := tf.GetForm()
	if err != nil {
		return nil, err
	}

	return json.Marshal(struct {
		Name string `json:"name"`
		Form *Form  `json:"form"`
	}{
		Name: tf.Name,
		Form: form,
	})
}

func (tf *ThemeForm) Create() error {
	return Insert(tf)
}

func (tf *ThemeForm) Update() error {
	return Update(tf)
}

func (tf *ThemeForm) Delete() error {
	return Delete(tf)
}

func FindThemeFormByThemeAndName(themeId int, name string) (*ThemeForm, error) {
	return SelectOne[ThemeForm](`
SELECT
    name,
    theme_id,
    form_id
FROM
    theme_form
WHERE
    theme_id = ?
    AND name = ?
`, themeId, name)
}

func (tf *ThemeForm) GetTheme() (*Theme, error) {
	return FindThemeById(tf.ThemeId)
}

func (tf *ThemeForm) GetForm() (*Form, error) {
	return FindFormById(tf.FormId)
}

func FindThemeFormsByTheme(themeId int) ([]ThemeForm, error) {
	return Select[ThemeForm](`
SELECT
    name,
    theme_id,
    form_id
FROM
    theme_form
WHERE
    theme_id = ?
`, themeId)
}
