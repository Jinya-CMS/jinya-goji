package database

type FileTag struct {
	Id    int     `db:"id,primarykey,autoincrement" json:"id"`
	Name  string  `db:"name" json:"name"`
	Color *string `db:"color" json:"color"`
	Emoji *string `db:"emoji" json:"emoji"`
}
