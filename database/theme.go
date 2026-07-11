package database

type Theme struct {
	Id            int               `db:"id,primarykey,autoincrement" json:"id"`
	Configuration map[string]any    `db:"configuration" json:"configuration"`
	Description   map[string]string `db:"description" json:"description"`
	Name          string            `db:"name" json:"name"`
	DisplayName   string            `db:"display_name" json:"displayName"`
	ScssVariables map[string]string `db:"scss_variables" json:"scssVariables"`
	HasApiTheme   bool              `db:"has_api_theme" json:"hasApi"`
}
