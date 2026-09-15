package database

type FileTag struct {
	Id    int    `db:"id,primarykey,autoincrement" json:"id"`
	Name  string `db:"name,notnull,unique" json:"name"`
	Color string `db:"color,notnull" json:"color"`
	Emoji string `db:"emoji,notnull" json:"emoji"`
}

func (tag *FileTag) Create() error {
	return dbMap.Insert(tag)
}

func (tag *FileTag) Update() error {
	_, err := dbMap.Update(tag)
	return err
}

func (tag *FileTag) Delete() error {
	_, err := dbMap.Delete(tag)
	return err
}

func FindFileTagById(id int) (*FileTag, error) {
	return dbMap.GetType[FileTag](id)
}

func FindFileTagByName(name string) (*FileTag, error) {
	return dbMap.SelectOneType[FileTag](`
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
	return dbMap.SelectType[FileTag](`
select
    id,
    name,
    color,
    emoji
from
    file_tag
`)
}
