package database

import (
	"encoding/json"
	"time"
)

type Gallery struct {
	Id            int       `db:"id,primarykey,autoincrement"`
	CreatorId     int       `db:"creator_id"`
	UpdatedById   int       `db:"updated_by_id"`
	CreatedAt     time.Time `db:"created_at"`
	LastUpdatedAt time.Time `db:"last_updated_at"`
	Name          string    `db:"name,unique"`
	Description   string    `db:"description"`
	Type          string    `db:"type"`
	Orientation   string    `db:"orientation"`
	Group         string    `db:"group"`
}

func (gallery *Gallery) MarshalJSON() ([]byte, error) {
	creator, err := gallery.GetCreator()
	if err != nil {
		return nil, err
	}
	updatedBy, err := gallery.GetUpdatedBy()
	if err != nil {
		return nil, err
	}

	return json.Marshal(struct {
		Id          int    `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Type        string `json:"type"`
		Orientation string `json:"orientation"`
		Group       string `json:"group"`
		Created     any    `json:"created"`
		Updated     any    `json:"updated"`
	}{
		Id:          gallery.Id,
		Name:        gallery.Name,
		Description: gallery.Description,
		Type:        gallery.Type,
		Orientation: gallery.Orientation,
		Group:       gallery.Group,
		Created: struct {
			At string `json:"at"`
			By any    `json:"by"`
		}{
			At: gallery.CreatedAt.Format(time.RFC3339),
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
			At: gallery.LastUpdatedAt.Format(time.RFC3339),
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

func (gallery *Gallery) Create() error {
	gallery.CreatedAt = time.Now()
	gallery.LastUpdatedAt = time.Now()
	return Insert(gallery)
}

func (gallery *Gallery) Update() error {
	gallery.LastUpdatedAt = time.Now()
	return Update(gallery)
}

func (gallery *Gallery) Delete() error {
	return Delete(gallery)
}

func FindGalleryById(id int) (*Gallery, error) {
	return Get[Gallery](id)
}

func FindAllGalleries() ([]Gallery, error) {
	return Select[Gallery](`
select id,
       creator_id,
       updated_by_id,
       created_at,
       last_updated_at,
       name,
       description,
       type,
       orientation,
       group
from gallery
order by id
`)
}

func (gallery *Gallery) GetFiles() ([]GalleryFilePosition, error) {
	return FindGalleryFilePositionsByGallery(gallery.Id)
}

func (gallery *Gallery) GetCreator() (*Artist, error) {
	return FindArtistById(gallery.CreatorId)
}

func (gallery *Gallery) GetUpdatedBy() (*Artist, error) {
	return FindArtistById(gallery.UpdatedById)
}
