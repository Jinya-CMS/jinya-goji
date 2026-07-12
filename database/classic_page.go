package database

import (
	"encoding/json"
	"time"
)

type ClassicPage struct {
	Id            int       `db:"id,primarykey,autoincrement"`
	CreatorId     int       `db:"creator_id"`
	UpdatedById   int       `db:"updated_by_id"`
	CreatedAt     time.Time `db:"created_at"`
	LastUpdatedAt time.Time `db:"last_updated_at"`
	Content       string    `db:"content"`
	Title         string    `db:"title,unique"`
}

func (page *ClassicPage) MarshalJSON() ([]byte, error) {
	creator, err := page.GetCreator()
	if err != nil {
		return nil, err
	}
	updatedBy, err := page.GetUpdatedBy()
	if err != nil {
		return nil, err
	}

	return json.Marshal(struct {
		Id      int    `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
		Created any    `json:"created"`
		Updated any    `json:"updated"`
	}{
		Id:      page.Id,
		Title:   page.Title,
		Content: page.Content,
		Created: struct {
			At string `json:"at"`
			By any    `json:"by"`
		}{
			At: page.CreatedAt.Format(time.RFC3339),
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
			At: page.LastUpdatedAt.Format(time.RFC3339),
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

func (page *ClassicPage) Create() error {
	page.CreatedAt = time.Now()
	page.LastUpdatedAt = time.Now()
	return Insert(page)
}

func (page *ClassicPage) Update() error {
	page.LastUpdatedAt = time.Now()
	return Update(page)
}

func (page *ClassicPage) Delete() error {
	return Delete(page)
}

func FindClassicPageById(id int) (*ClassicPage, error) {
	return Get[ClassicPage](id)
}

func FindAllClassicPages() ([]ClassicPage, error) {
	return Select[ClassicPage](`
SELECT
    id,
    creator_id,
    updated_by_id,
    created_at,
    last_updated_at,
    content,
    title
FROM
    classic_page
`)
}

func (page *ClassicPage) GetCreator() (*Artist, error) {
	return FindArtistById(page.CreatorId)
}

func (page *ClassicPage) GetUpdatedBy() (*Artist, error) {
	return FindArtistById(page.UpdatedById)
}
