package database

type BlogPostSection struct {
	Id         int     `db:"id,primarykey,autoincrement"`
	GalleryId  *int    `db:"gallery_id"`
	FileId     *int    `db:"file_id"`
	Html       *string `db:"html"`
	Link       *string `db:"link"`
	BlogPostId int     `db:"blog_post_id"`
	Position   int     `db:"position"`
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

func FindBlogPostSectionById(id int) (*BlogPostSection, error) {
	return Get[BlogPostSection](id)
}

func FindBlogPostSectionsByPost(postId int) ([]BlogPostSection, error) {
	return Select[BlogPostSection](`
SELECT
    id,
    gallery_id,
    file_id,
    html,
    link,
    blog_post_id,
    position
FROM
    blog_post_section
WHERE
    blog_post_id = ?
ORDER BY
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
