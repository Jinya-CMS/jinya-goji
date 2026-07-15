package database

import "encoding/json"

type BlogPostSection struct {
	Id         int     `db:"id,primarykey,autoincrement"`
	GalleryId  *int    `db:"gallery_id"`
	FileId     *int    `db:"file_id"`
	Html       *string `db:"html"`
	Link       *string `db:"link"`
	BlogPostId int     `db:"blog_post_id"`
	Position   int     `db:"position"`
}

func (section *BlogPostSection) MarshalJSON() ([]byte, error) {
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
			if section.Link != nil {
				data["link"] = section.Link
			}
		}
	} else {
		data["html"] = section.Html
	}

	return json.Marshal(data)
}

func (section *BlogPostSection) Create() error {
	return Insert(section)
}

func (section *BlogPostSection) Update() error {
	return Update(section)
}

func (section *BlogPostSection) Delete() error {
	return Delete(section)
}

func FindBlogPostSectionsByPost(postId int) ([]BlogPostSection, error) {
	return Select[BlogPostSection](`
select
    id,
    gallery_id,
    file_id,
    html,
    link,
    blog_post_id,
    position
from
    blog_post_section
where
    blog_post_id = $1
order by
    position
`, postId)
}

func (section *BlogPostSection) GetBlogPost() (*BlogPost, error) {
	return FindBlogPostById(section.BlogPostId)
}

func (section *BlogPostSection) GetFile() (*File, error) {
	if section.FileId == nil {
		return nil, nil
	}
	return FindFileById(*section.FileId)
}

func (section *BlogPostSection) GetGallery() (*Gallery, error) {
	if section.GalleryId == nil {
		return nil, nil
	}
	return FindGalleryById(*section.GalleryId)
}
