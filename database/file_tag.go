package database

import "encoding/json"

type FileTag struct {
	Id    int     `db:"id,primarykey,autoincrement"`
	Name  string  `db:"name,unique"`
	Color *string `db:"color"`
	Emoji *string `db:"emoji"`
}

func (tag *FileTag) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Id    int     `json:"id"`
		Name  string  `json:"name"`
		Color *string `json:"color"`
		Emoji *string `json:"emoji"`
	}{
		Id:    tag.Id,
		Name:  tag.Name,
		Color: tag.Color,
		Emoji: tag.Emoji,
	})
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
