package database

import (
	"encoding/json"
)

type MenuItem struct {
	Id            int        `db:"id,primarykey,autoincrement"`
	Position      int        `db:"position"`
	MenuId        *int       `db:"menu_id"`
	ParentId      *int       `db:"parent_id"`
	Title         string     `db:"title"`
	Route         *string    `db:"route"`
	Highlighted   bool       `db:"highlighted"`
	ArtistId      *int       `db:"artist_id"`
	GalleryId     *int       `db:"gallery_id"`
	ModernPageId  *int       `db:"modern_page_id"`
	ClassicPageId *int       `db:"classic_page_id"`
	CategoryId    *int       `db:"blog_category_id"`
	BlogHomePage  *bool      `db:"blog_home_page"`
	Children      []MenuItem `db:"-"`
}

func (item *MenuItem) MarshalJSON() ([]byte, error) {
	childItems, err := item.GetItems()
	if err != nil {
		return nil, err
	}
	if childItems == nil {
		childItems = []MenuItem{}
	}

	data := map[string]any{
		"id":           item.Id,
		"position":     item.Position,
		"highlighted":  item.Highlighted,
		"title":        item.Title,
		"route":        item.Route,
		"blogHomePage": item.BlogHomePage,
		"items":        childItems,
	}

	if item.ArtistId != nil {
		artist, err := item.GetArtist()
		if err != nil {
			return nil, err
		}
		if artist != nil {
			data["artist"] = map[string]any{
				"id":         artist.Id,
				"artistName": artist.ArtistName,
				"email":      artist.Email,
			}
		}
	} else if item.ClassicPageId != nil {
		classicPage, err := item.GetClassicPage()
		if err != nil {
			return nil, err
		}
		if classicPage != nil {
			data["classicPage"] = map[string]any{
				"id":    classicPage.Id,
				"title": classicPage.Title,
			}
		}
	} else if item.ModernPageId != nil {
		modernPage, err := item.GetModernPage()
		if err != nil {
			return nil, err
		}
		if modernPage != nil {
			data["modernPage"] = map[string]any{
				"id":   modernPage.Id,
				"name": modernPage.Name,
			}
		}
	} else if item.GalleryId != nil {
		gallery, err := item.GetGallery()
		if err != nil {
			return nil, err
		}
		if gallery != nil {
			data["gallery"] = map[string]any{
				"id":   gallery.Id,
				"name": gallery.Name,
			}
		}
	} else if item.CategoryId != nil {
		category, err := item.GetBlogCategory()
		if err != nil {
			return nil, err
		}
		if category != nil {
			data["blogCategory"] = map[string]any{
				"id":   category.Id,
				"name": category.Name,
			}
		}
	}

	return json.Marshal(data)
}

func (item *MenuItem) Create() error {
	return dbMap.Insert(item)
}

func (item *MenuItem) Update() error {
	_, err := dbMap.Update(item)
	return err
}

func (item *MenuItem) Delete() error {
	_, err := dbMap.Delete(item)
	return err
}

func FindMenuItemById(id int) (*MenuItem, error) {
	return dbMap.GetType[MenuItem](id)
}

func FindMenuItemsByMenu(menuId int) ([]MenuItem, error) {
	return dbMap.SelectType[MenuItem](`
select id,
       position,
       menu_id,
       parent_id,
       title,
       route,
       highlighted,
       artist_id,
       gallery_id,
       modern_page_id,
       classic_page_id,
       blog_category_id,
       blog_home_page
from menu_item
where menu_id = $1
order by position
`, menuId)
}

func FindMenuItemsByParent(parentId int) ([]MenuItem, error) {
	return dbMap.SelectType[MenuItem](`
select id,
       position,
       menu_id,
       parent_id,
       title,
       route,
       highlighted,
       artist_id,
       gallery_id,
       modern_page_id,
       classic_page_id,
       blog_category_id,
       blog_home_page
from menu_item
where parent_id = $1
order by position
`, parentId)
}

func FindMenuItemByRoute(route string) (*MenuItem, error) {
	return dbMap.SelectOneType[MenuItem](`
select id,
       position,
       menu_id,
       parent_id,
       title,
       route,
       highlighted,
       artist_id,
       gallery_id,
       modern_page_id,
       classic_page_id,
       blog_category_id,
       blog_home_page
from menu_item
where route = $1
   or route = $2
`, route, "/"+route)
}

func (item *MenuItem) GetParent() (*MenuItem, error) {
	if item.ParentId == nil {
		return nil, nil
	}
	return FindMenuItemById(*item.ParentId)
}

func (item *MenuItem) GetMenu() (*Menu, error) {
	if item.MenuId == nil {
		return nil, nil
	}
	return FindMenuById(*item.MenuId)
}

func (item *MenuItem) GetArtist() (*Artist, error) {
	if item.ArtistId == nil {
		return nil, nil
	}
	return FindArtistById(*item.ArtistId)
}

func (item *MenuItem) GetClassicPage() (*ClassicPage, error) {
	if item.ClassicPageId == nil {
		return nil, nil
	}
	return FindClassicPageById(*item.ClassicPageId)
}

func (item *MenuItem) GetModernPage() (*ModernPage, error) {
	if item.ModernPageId == nil {
		return nil, nil
	}
	return FindModernPageById(*item.ModernPageId)
}

func (item *MenuItem) GetGallery() (*Gallery, error) {
	if item.GalleryId == nil {
		return nil, nil
	}
	return FindGalleryById(*item.GalleryId)
}

func (item *MenuItem) GetBlogCategory() (*BlogCategory, error) {
	if item.CategoryId == nil {
		return nil, nil
	}
	return FindBlogCategoryById(*item.CategoryId)
}

func (item *MenuItem) GetItems() ([]MenuItem, error) {
	return FindMenuItemsByParent(item.Id)
}
