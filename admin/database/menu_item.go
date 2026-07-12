package database

import "encoding/json"

type MenuItem struct {
	Id            int     `db:"id,primarykey,autoincrement"`
	Position      int     `db:"position"`
	MenuId        *int    `db:"menu_id"`
	ParentId      *int    `db:"parent_id"`
	Title         string  `db:"title"`
	Route         *string `db:"route"`
	Highlighted   bool    `db:"highlighted"`
	FormId        *int    `db:"form_id"`
	ArtistId      *int    `db:"artist_id"`
	GalleryId     *int    `db:"gallery_id"`
	ModernPageId  *int    `db:"segment_page_id"`
	ClassicPageId *int    `db:"page_id"`
	CategoryId    *int    `db:"category_id"`
	BlogHomePage  *bool   `db:"blog_home_page"`
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

	if item.FormId != nil {
		form, err := item.GetForm()
		if err != nil {
			return nil, err
		}
		if form != nil {
			data["form"] = map[string]any{
				"id":    form.Id,
				"title": form.Title,
			}
		}
	} else if item.ArtistId != nil {
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
	return Insert(item)
}

func (item *MenuItem) Update() error {
	return Update(item)
}

func (item *MenuItem) Delete() error {
	return Delete(item)
}

func FindMenuItemById(id int) (*MenuItem, error) {
	return Get[MenuItem](id)
}

func FindMenuItemsByMenu(menuId int) ([]MenuItem, error) {
	return Select[MenuItem](`
SELECT
    id,
    position,
    menu_id,
    parent_id,
    title,
    route,
    highlighted,
    form_id,
    artist_id,
    gallery_id,
    segment_page_id,
    page_id,
    category_id,
    blog_home_page
FROM
    menu_item
WHERE
    menu_id = ?
ORDER BY
    position
`, menuId)
}

func FindMenuItemsByParent(parentId int) ([]MenuItem, error) {
	return Select[MenuItem](`
SELECT
    id,
    position,
    menu_id,
    parent_id,
    title,
    route,
    highlighted,
    form_id,
    artist_id,
    gallery_id,
    segment_page_id,
    page_id,
    category_id,
    blog_home_page
FROM
    menu_item
WHERE
    parent_id = ?
ORDER BY
    position
`, parentId)
}

func FindMenuItemByRoute(route string) (*MenuItem, error) {
	return SelectOne[MenuItem](`
SELECT
    id,
    position,
    menu_id,
    parent_id,
    title,
    route,
    highlighted,
    form_id,
    artist_id,
    gallery_id,
    segment_page_id,
    page_id,
    category_id,
    blog_home_page
FROM
    menu_item
WHERE
    route = ? OR route = ?
`, route, "/"+route)
}

func FindMenuItemByMenuAndPosition(menuId int, position int) (*MenuItem, error) {
	return SelectOne[MenuItem](`
SELECT
    id,
    position,
    menu_id,
    parent_id,
    title,
    route,
    highlighted,
    form_id,
    artist_id,
    gallery_id,
    segment_page_id,
    page_id,
    category_id,
    blog_home_page
FROM
    menu_item
WHERE
    menu_id = ? AND position = ?
`, menuId, position)
}

func FindMenuItemByMenuItemAndPosition(parentId int, position int) (*MenuItem, error) {
	return SelectOne[MenuItem](`
SELECT
    id,
    position,
    menu_id,
    parent_id,
    title,
    route,
    highlighted,
    form_id,
    artist_id,
    gallery_id,
    segment_page_id,
    page_id,
    category_id,
    blog_home_page
FROM
    menu_item
WHERE
    parent_id = ? AND position = ?
`, parentId, position)
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

func (item *MenuItem) GetForm() (*Form, error) {
	if item.FormId == nil {
		return nil, nil
	}
	return FindFormById(*item.FormId)
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
