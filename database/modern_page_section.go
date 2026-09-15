package database

import "encoding/json"

type ModernPageSection struct {
	Id        int     `db:"id,primarykey,autoincrement"`
	Position  int     `db:"position"`
	GalleryId *int    `db:"gallery_id"`
	FileId    *int    `db:"file_id"`
	Html      *string `db:"html"`
	PageId    int     `db:"page_id"`
	Action    *string `db:"action"`
	Target    *string `db:"target"`
}

func (section *ModernPageSection) MarshalJSON() ([]byte, error) {
	data := map[string]any{
		"position": section.Position,
		"id":       section.Id,
	}

	if section.GalleryId != nil {
		gallery, err := section.GetGallery()
		if err != nil {
			return nil, err
		}
		if gallery != nil {
			data["gallery"] = map[string]any{
				"id":          gallery.Id,
				"name":        gallery.Name,
				"description": gallery.Description,
				"type":        gallery.Type,
				"orientation": gallery.Orientation,
			}
		}
	} else if section.FileId != nil {
		file, err := section.GetFile()
		if err != nil {
			return nil, err
		}
		if file != nil {
			data["file"] = map[string]any{
				"id":   file.Id,
				"name": file.Name,
				"type": file.Type,
				"path": file.Path,
			}
			if section.Action != nil && *section.Action == "link" {
				data["link"] = section.Target
			}
		}
	} else {
		data["html"] = section.Html
	}

	return json.Marshal(data)
}

func (section *ModernPageSection) Create() error {
	return dbMap.Insert(section)
}

func (section *ModernPageSection) Update() error {
	_, err := dbMap.Update(section)
	return err
}

func (section *ModernPageSection) Delete() error {
	_, err := dbMap.Delete(section)
	return err
}

func FindModernPageSectionsByPage(pageId int) ([]ModernPageSection, error) {
	return dbMap.SelectType[ModernPageSection](`
select id,
       position,
       gallery_id,
       file_id,
       html,
       page_id,
       action,
       target
from modern_page_section
where page_id = $1
order by position
`, pageId)
}

func (section *ModernPageSection) GetGallery() (*Gallery, error) {
	if section.GalleryId == nil {
		return nil, nil
	}
	return FindGalleryById(*section.GalleryId)
}

func (section *ModernPageSection) GetFile() (*File, error) {
	if section.FileId == nil {
		return nil, nil
	}
	return FindFileById(*section.FileId)
}

func (section *ModernPageSection) GetModernPage() (*ModernPage, error) {
	return FindModernPageById(section.PageId)
}
