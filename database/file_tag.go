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
SELECT
    id,
    name,
    color,
    emoji
FROM
    file_tag
WHERE
    name = ?
`, name)
}

func FindAllFileTags() ([]FileTag, error) {
	return Select[FileTag](`
SELECT
    id,
    name,
    color,
    emoji
FROM
    file_tag
`)
}
