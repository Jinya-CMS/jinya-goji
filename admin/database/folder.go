package database

import "encoding/json"

type Folder struct {
	Id       int    `db:"id,primarykey,autoincrement"`
	Name     string `db:"name"`
	ParentId *int   `db:"parent_id"`
}

func (folder *Folder) MarshalJSON() ([]byte, error) {
	files, err := folder.GetFiles()
	if err != nil {
		return nil, err
	}
	folders, err := folder.GetFolders()
	if err != nil {
		return nil, err
	}

	if files == nil {
		files = []File{}
	}
	if folders == nil {
		folders = []Folder{}
	}

	return json.Marshal(struct {
		Id      int      `json:"id"`
		Name    string   `json:"name"`
		Files   []File   `json:"files"`
		Folders []Folder `json:"folders"`
	}{
		Id:      folder.Id,
		Name:    folder.Name,
		Files:   files,
		Folders: folders,
	})
}

func (folder *Folder) Create() error {
	return Insert(folder)
}

func (folder *Folder) Update() error {
	return Update(folder)
}

func (folder *Folder) Delete() error {
	return Delete(folder)
}

func FindFolderById(id int) (*Folder, error) {
	return Get[Folder](id)
}

func FindAllFolders() ([]Folder, error) {
	return Select[Folder](`
SELECT
    id,
    name,
    parent_id
FROM
    folder
`)
}

func FindRootFolders() ([]Folder, error) {
	return Select[Folder](`
SELECT
    id,
    name,
    parent_id
FROM
    folder
WHERE
    parent_id IS NULL
`)
}

func (folder *Folder) GetFiles() ([]File, error) {
	return Select[File](`
SELECT
    id,
    creator_id,
    updated_by_id,
    created_at,
    last_updated_at,
    path,
    name,
    folder_id,
    type
FROM
    file
WHERE
    folder_id = ?
`, folder.Id)
}

func (folder *Folder) GetFolders() ([]Folder, error) {
	return Select[Folder](`
SELECT
    id,
    name,
    parent_id
FROM
    folder
WHERE
    parent_id = ?
`, folder.Id)
}
