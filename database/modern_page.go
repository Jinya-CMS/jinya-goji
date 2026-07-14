package database

import (
	"encoding/json"
	"time"
)

type ModernPage struct {
	Id            int       `db:"id,primarykey,autoincrement"`
	CreatorId     int       `db:"creator_id"`
	UpdatedById   int       `db:"updated_by_id"`
	CreatedAt     time.Time `db:"created_at"`
	LastUpdatedAt time.Time `db:"last_updated_at"`
	Name          string    `db:"name,unique"`
}

func (page *ModernPage) MarshalJSON() ([]byte, error) {
	creator, err := page.GetCreator()
	if err != nil {
		return nil, err
	}
	updatedBy, err := page.GetUpdatedBy()
	if err != nil {
		return nil, err
	}
	sections, err := page.GetSections()
	if err != nil {
		return nil, err
	}

	return json.Marshal(struct {
		Id           int    `json:"id"`
		Name         string `json:"name"`
		SectionCount int    `json:"sectionCount"`
		Created      any    `json:"created"`
		Updated      any    `json:"updated"`
	}{
		Id:           page.Id,
		Name:         page.Name,
		SectionCount: len(sections),
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

func (page *ModernPage) Create() error {
	page.CreatedAt = time.Now()
	page.LastUpdatedAt = time.Now()
	return Insert(page)
}

func (page *ModernPage) Update() error {
	page.LastUpdatedAt = time.Now()
	return Update(page)
}

func (page *ModernPage) Delete() error {
	return Delete(page)
}

func FindModernPageById(id int) (*ModernPage, error) {
	return Get[ModernPage](id)
}

func FindAllModernPages() ([]ModernPage, error) {
	return Select[ModernPage](`
select id,
       creator_id,
       updated_by_id,
       created_at,
       last_updated_at,
       name
from modern_page
order by id
`)
}

func (page *ModernPage) GetCreator() (*Artist, error) {
	return FindArtistById(page.CreatorId)
}

func (page *ModernPage) GetUpdatedBy() (*Artist, error) {
	return FindArtistById(page.UpdatedById)
}

func (page *ModernPage) GetSections() ([]ModernPageSection, error) {
	return FindModernPageSectionsByPage(page.Id)
}

func (page *ModernPage) ReplaceSections(sections []ModernPageSection) error {
	tx, err := dbMap.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.Exec(`
delete
from modern_page_section
where page_id = ?`, page.Id)
	if err != nil {
		return err
	}

	for i, section := range sections {
		section.PageId = page.Id
		section.Position = i
		err = tx.Insert(&section)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
