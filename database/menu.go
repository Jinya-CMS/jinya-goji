package database

import "encoding/json"

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
	return Insert(menu)
}

func (menu *Menu) Update() error {
	return Update(menu)
}

func (menu *Menu) Delete() error {
	return Delete(menu)
}

func FindMenuById(id int) (*Menu, error) {
	return Get[Menu](id)
}

func FindAllMenus() ([]Menu, error) {
	return Select[Menu](`
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

func (menu *Menu) GetLogo() (*File, error) {
	if menu.Logo == nil {
		return nil, nil
	}
	return FindFileById(*menu.Logo)
}
