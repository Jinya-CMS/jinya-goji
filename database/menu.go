package database

type Menu struct {
	Id   int    `db:"id,primarykey,autoincrement" json:"id"`
	Name string `db:"name" json:"name"`
	Logo *int   `db:"logo" json:"-"`
}
