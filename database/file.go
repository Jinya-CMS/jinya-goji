package database

import (
	"encoding/json"
	"time"
)

type File struct {
	Id            int       `db:"id,primarykey,autoincrement"`
	CreatorId     int       `db:"creator_id"`
	UpdatedById   int       `db:"updated_by_id"`
	CreatedAt     time.Time `db:"created_at"`
	LastUpdatedAt time.Time `db:"last_updated_at"`
	Path          string    `db:"path"`
	Name          string    `db:"name"`
	FolderId      *int      `db:"folder_id"`
	Type          string    `db:"type"`
}

func (file *File) MarshalJSON() ([]byte, error) {
	creator, err := file.GetCreator()
	if err != nil {
		return nil, err
	}
	updatedBy, err := file.GetUpdatedBy()
	if err != nil {
		return nil, err
	}
	tags, err := file.GetTags()
	if err != nil {
		return nil, err
	}

	if tags == nil {
		tags = []FileTag{}
	}

	return json.Marshal(struct {
		Id       int       `json:"id"`
		Name     string    `json:"name"`
		Type     string    `json:"type"`
		Path     string    `json:"path"`
		FolderId *int      `json:"folderId"`
		Tags     []FileTag `json:"tags"`
		Created  any       `json:"created"`
		Updated  any       `json:"updated"`
	}{
		Id:       file.Id,
		Name:     file.Name,
		Type:     file.Type,
		Path:     file.Path,
		FolderId: file.FolderId,
		Tags:     tags,
		Created: struct {
			At string `json:"at"`
			By any    `json:"by"`
		}{
			At: file.CreatedAt.Format(time.RFC3339),
			By: struct {
				ArtistName     string  `json:"artistName"`
				Email          string  `json:"email"`
				ProfilePicture *string `json:"profilePicture"`
			}{
				ArtistName:     creator.ArtistName,
				Email:          creator.Email,
				ProfilePicture: creator.ProfilePicture,
			},
		},
		Updated: struct {
			At string `json:"at"`
			By any    `json:"by"`
		}{
			At: file.LastUpdatedAt.Format(time.RFC3339),
			By: struct {
				ArtistName     string  `json:"artistName"`
				Email          string  `json:"email"`
				ProfilePicture *string `json:"profilePicture"`
			}{
				ArtistName:     updatedBy.ArtistName,
				Email:          updatedBy.Email,
				ProfilePicture: updatedBy.ProfilePicture,
			},
		},
	})
}

func (file *File) Create(tags []FileTag) error {
	tx, err := dbMap.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	file.CreatedAt = time.Now()
	file.LastUpdatedAt = time.Now()
	err = tx.Insert(file)
	if err != nil {
		return err
	}

	for _, tag := range tags {
		_, err = tx.Exec(`
insert into file_tag_file (file_tag_id, file_id) 
values ($1, $2)`, tag.Id, file.Id)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (file *File) Tag(tags []FileTag) error {
	tx, err := dbMap.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.Exec(`
delete
from file_tag_file
where file_id = $1`, file.Id)
	if err != nil {
		return err
	}

	for _, tag := range tags {
		_, err = tx.Exec(`
insert into file_tag_file (file_tag_id, file_id) 
values ($1, $2)`, tag.Id, file.Id)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (file *File) Update() error {
	file.LastUpdatedAt = time.Now()
	return Update(file)
}

func (file *File) Delete() error {
	return Delete(file)
}

func FindFileById(id int) (*File, error) {
	return Get[File](id)
}

func FilePathExists(path string) bool {
	fileCount, err := dbMap.SelectInt(`
select count(id)
from file
where path = $1`, path)

	if err == nil && fileCount > 0 {
		return true
	}

	profilePictures, err := dbMap.SelectInt(`
select count(id)
from users
where profile_picture = $1`, path)

	return err == nil && profilePictures > 0
}

func FindAllFiles() ([]File, error) {
	return Select[File](`
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
`)
}

func FindRootFiles() ([]File, error) {
	return Select[File](`
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
where folder_id is null
`)
}

func (file *File) GetCreator() (*Artist, error) {
	return FindArtistById(file.CreatorId)
}

func (file *File) GetUpdatedBy() (*Artist, error) {
	return FindArtistById(file.UpdatedById)
}

func (file *File) GetUploadChunks() ([]UploadingFileChunk, error) {
	uploadingFile, err := FindUploadingFileByFile(file.Id)
	if err != nil {
		return nil, err
	}
	if uploadingFile == nil {
		return []UploadingFileChunk{}, nil
	}
	return uploadingFile.GetChunks()
}

func (file *File) GetTags() ([]FileTag, error) {
	return Select[FileTag](`
select ft.id,
       ft.name,
       ft.color,
       ft.emoji
from file_tag ft
         join
     file_tag_file ftf on ft.id = ftf.file_tag_id
where ftf.file_id = $1
`, file.Id)
}
