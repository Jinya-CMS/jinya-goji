package database

type Folder struct {
	Id       int    `db:"id,primarykey,autoincrement" json:"id"`
	Name     string `db:"name" json:"name"`
	ParentId *int   `db:"parent_id" json:"-"`
}
