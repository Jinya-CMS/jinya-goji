package database

import (
	"encoding/json"
	"time"
)

type Form struct {
	Id            int       `db:"id,primarykey,autoincrement"`
	CreatorId     int       `db:"creator_id"`
	UpdatedById   int       `db:"updated_by_id"`
	CreatedAt     time.Time `db:"created_at"`
	LastUpdatedAt time.Time `db:"last_updated_at"`
	Title         string    `db:"title,unique"`
	Description   string    `db:"description"`
	ToAddress     string    `db:"to_address"`
}

func (form *Form) MarshalJSON() ([]byte, error) {
	creator, err := form.GetCreator()
	if err != nil {
		return nil, err
	}
	updatedBy, err := form.GetUpdatedBy()
	if err != nil {
		return nil, err
	}

	return json.Marshal(struct {
		Id          int    `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		ToAddress   string `json:"toAddress"`
		Created     any    `json:"created"`
		Updated     any    `json:"updated"`
	}{
		Id:          form.Id,
		Title:       form.Title,
		Description: form.Description,
		ToAddress:   form.ToAddress,
		Created: struct {
			At string `json:"at"`
			By any    `json:"by"`
		}{
			At: form.CreatedAt.Format(time.RFC3339),
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
			At: form.LastUpdatedAt.Format(time.RFC3339),
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

func (form *Form) Create() error {
	form.CreatedAt = time.Now()
	form.LastUpdatedAt = time.Now()
	return Insert(form)
}

func (form *Form) Update() error {
	form.LastUpdatedAt = time.Now()
	return Update(form)
}

func (form *Form) Delete() error {
	return Delete(form)
}

func FindFormById(id int) (*Form, error) {
	return Get[Form](id)
}

func FindAllForms() ([]Form, error) {
	return Select[Form](`
SELECT
    id,
    creator_id,
    updated_by_id,
    created_at,
    last_updated_at,
    title,
    description,
    to_address
FROM
    form
`)
}

func (form *Form) GetItems() ([]FormItem, error) {
	return FindFormItemsByForm(form.Id)
}

func (form *Form) GetCreator() (*Artist, error) {
	return FindArtistById(form.CreatorId)
}

func (form *Form) GetUpdatedBy() (*Artist, error) {
	return FindArtistById(form.UpdatedById)
}

func (form *Form) ReplaceItems(items []FormItem) error {
	tx, err := dbMap.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM form_item WHERE form_id = ?", form.Id)
	if err != nil {
		return err
	}

	for i, item := range items {
		item.FormId = form.Id
		item.Position = i
		err = tx.Insert(&item)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
