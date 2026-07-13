package database

type FileTag struct {
	Id    int    `db:"id,primarykey,autoincrement" json:"id"`
	Name  string `db:"name,notnull,unique" json:"name"`
	Color string `db:"color,notnull" json:"color"`
	Emoji string `db:"emoji,notnull" json:"emoji"`
}

func (tag *FileTag) Create() error {
	return Insert(tag)
}

func (tag *FileTag) Update() error {
	return Update(tag)
}

func (tag *FileTag) Delete() error {
	return Delete(tag)
}

func FindFileTagById(id int) (*FileTag, error) {
	return Get[FileTag](id)
}

func FindFileTagByName(name string) (*FileTag, error) {
	return SelectOne[FileTag](`
select
    id,
    name,
    color,
    emoji
from
    file_tag
where
    name = $1
`, name)
}

func FindAllFileTags() ([]FileTag, error) {
	return Select[FileTag](`
select
    id,
    name,
    color,
    emoji
from
    file_tag
`)
}
