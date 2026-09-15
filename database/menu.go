package database

import (
	"encoding/json"

	"github.com/DerKnerd/gorp/v4"
)

type Menu struct {
	Id   int    `db:"id,primarykey,autoincrement"`
	Name string `db:"name"`
	Logo *int   `db:"logo"`
}

func (menu *Menu) MarshalJSON() ([]byte, error) {
	logo, err := menu.GetLogo()
	if err != nil {
		return nil, err
	}
	var logoData any
	if logo != nil {
		logoData = struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		}{
			Id:   logo.Id,
			Name: logo.Name,
		}
	}

	return json.Marshal(struct {
		Name string `json:"name"`
		Id   int    `json:"id"`
		Logo any    `json:"logo,omitzero"`
	}{
		Name: menu.Name,
		Id:   menu.Id,
		Logo: logoData,
	})
}

func (menu *Menu) Create() error {
	return dbMap.Insert(menu)
}

func (menu *Menu) Update() error {
	_, err := dbMap.Update(menu)
	return err
}

func (menu *Menu) Delete() error {
	_, err := dbMap.Delete(menu)
	return err
}

func FindMenuById(id int) (*Menu, error) {
	return dbMap.GetType[Menu](id)
}

func FindAllMenus() ([]Menu, error) {
	return dbMap.SelectType[Menu](`
select id,
       name,
       logo
from menu
order by id
`)
}

func (menu *Menu) GetItems() ([]MenuItem, error) {
	return FindMenuItemsByMenu(menu.Id)
}

func (menu *Menu) ReplaceItems(items []MenuItem) error {
	tx, err := dbMap.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.Exec(`
delete
from menu_item
where menu_id = $1`, menu.Id)

	err = internalReplaceItems(-1, items, tx)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func internalReplaceItems(parentId int, items []MenuItem, tx *gorp.Transaction) error {
	for i, item := range items {
		if parentId != -1 {
			item.ParentId = &parentId
		}

		item.Position = i

		err := tx.Insert(&item)
		if err != nil {
			return err
		}

		err = internalReplaceItems(item.Id, item.Children, tx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (menu *Menu) GetLogo() (*File, error) {
	if menu.Logo == nil {
		return nil, nil
	}
	return FindFileById(*menu.Logo)
}
