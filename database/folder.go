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
	return dbMap.Insert(folder)
}

func (folder *Folder) Update() error {
	_, err := dbMap.Update(folder)
	return err
}

func (folder *Folder) Delete() error {
	_, err := dbMap.Delete(folder)
	return err
}

func FindFolderById(id int) (*Folder, error) {
	return dbMap.GetType[Folder](id)
}

func FindRootFolders() ([]Folder, error) {
	return dbMap.SelectType[Folder](`
select id,
       name,
       parent_id
from folder
where parent_id is null
`)
}

func (folder *Folder) GetFiles() ([]File, error) {
	return dbMap.SelectType[File](`
select id,
       creator_id,
       updated_by_id,
       created_at,
       last_updated_at,
       path,
       name,
       folder_id,
       type
from file
where folder_id = $1
`, folder.Id)
}

func (folder *Folder) GetFolders() ([]Folder, error) {
	return dbMap.SelectType[Folder](`
select id,
       name,
       parent_id
from folder
where parent_id = $1
`, folder.Id)
}
