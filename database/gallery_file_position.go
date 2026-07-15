package database

import "encoding/json"

type GalleryFilePosition struct {
	Id        int `db:"id,primarykey,autoincrement"`
	Position  int `db:"position"`
	GalleryId int `db:"gallery_id"`
	FileId    int `db:"file_id"`
}

func (position *GalleryFilePosition) MarshalJSON() ([]byte, error) {
	gallery, err := position.GetGallery()
	if err != nil {
		return nil, err
	}
	file, err := position.GetFile()
	if err != nil {
		return nil, err
	}

	return json.Marshal(struct {
		Gallery  any `json:"gallery"`
		File     any `json:"file"`
		Id       int `json:"id"`
		Position int `json:"position"`
	}{
		Id:       position.Id,
		Position: position.Position,
		Gallery: struct {
			Id          int    `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
		}{
			Id:          position.GalleryId,
			Name:        gallery.Name,
			Description: gallery.Description,
		},
		File: struct {
			Path string `json:"path"`
			Id   int    `json:"id"`
			Name string `json:"name"`
			Type string `json:"type"`
		}{
			Path: file.Path,
			Id:   file.Id,
			Name: file.Name,
			Type: file.Type,
		},
	})
}

func (position *GalleryFilePosition) Create() error {
	err := position.internalRearrange(position.GalleryId, position.Position)
	if err != nil {
		return err
	}

	err = Insert(position)
	if err != nil {
		return err
	}

	return position.resetOrder(position.GalleryId)
}

func (position *GalleryFilePosition) Update() error {
	return Update(position)
}

func (position *GalleryFilePosition) Delete() error {
	err := Delete(position)
	if err != nil {
		return err
	}

	err = position.internalRearrange(position.GalleryId, -1)
	if err != nil {
		return err
	}

	return position.resetOrder(position.GalleryId)
}

func (position *GalleryFilePosition) Move(newPosition int) error {
	if newPosition == position.Position {
		return nil
	}

	err := position.internalRearrange(position.GalleryId, newPosition)
	if err != nil {
		return err
	}

	err = position.Update()
	if err != nil {
		return err
	}

	return position.resetOrder(position.GalleryId)
}

func (position *GalleryFilePosition) internalRearrange(parentId int, newPosition int) error {
	target := newPosition
	if newPosition > position.Position {
		target++
	}

	data, err := Select[GalleryFilePosition](`
select id,
       position,
       gallery_id,
       file_id
from gallery_file_position
where position >= $1
  and gallery_id = $2
order by position
`, target, parentId)
	if err != nil {
		return err
	}

	previousPosition := target
	for _, item := range data {
		previousPosition++
		_, err = dbMap.Exec(`
update gallery_file_position
set position = $1
where id = $2`, previousPosition, item.Id)
		if err != nil {
			return err
		}
	}

	position.Position = target
	return nil
}

func (position *GalleryFilePosition) resetOrder(parentId int) error {
	data, err := Select[GalleryFilePosition](`
select id,
       position,
       gallery_id,
       file_id
from gallery_file_position
where gallery_id = $1
order by position
`, parentId)
	if err != nil {
		return err
	}

	for idx, item := range data {
		_, err = dbMap.Exec(`
update gallery_file_position
set position = $1
where id = $2`, idx, item.Id)
		if err != nil {
			return err
		}
	}

	return nil
}

func FindGalleryFilePositionByPosition(galleryId int, position int) (*GalleryFilePosition, error) {
	return SelectOne[GalleryFilePosition](`
select id,
       position,
       gallery_id,
       file_id
from gallery_file_position
where gallery_id = $1
  and position = $2
`, galleryId, position)
}

func FindGalleryFilePositionsByGallery(galleryId int) ([]GalleryFilePosition, error) {
	return Select[GalleryFilePosition](`
select id,
       position,
       gallery_id,
       file_id
from gallery_file_position
where gallery_id = $1
order by position
`, galleryId)
}

func (position *GalleryFilePosition) GetGallery() (*Gallery, error) {
	return FindGalleryById(position.GalleryId)
}

func (position *GalleryFilePosition) GetFile() (*File, error) {
	return FindFileById(position.FileId)
}
