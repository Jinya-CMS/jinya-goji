package database

import (
	"encoding/json"
	"time"
)

type BlogPost struct {
	Id            int       `db:"id,primarykey,autoincrement"`
	Title         string    `db:"title,unique"`
	Slug          string    `db:"slug,unique"`
	HeaderImageId *int      `db:"header_image_id"`
	Public        bool      `db:"public"`
	CreatedAt     time.Time `db:"created_at"`
	CreatorId     int       `db:"creator_id"`
	UpdatedById   int       `db:"updated_by_id"`
	LastUpdatedAt time.Time `db:"last_updated_at"`
	CategoryId    *int      `db:"category_id"`
}

func (post *BlogPost) MarshalJSON() ([]byte, error) {
	creator, err := post.GetCreator()
	if err != nil {
		return nil, err
	}
	updatedBy, err := post.GetUpdatedBy()
	if err != nil {
		return nil, err
	}
	headerImage, err := post.GetHeaderImage()
	if err != nil {
		return nil, err
	}
	category, err := post.GetCategory()
	if err != nil {
		return nil, err
	}

	var headerFormat any
	if headerImage != nil {
		headerFormat = struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
			Path string `json:"path"`
		}{
			Id:   headerImage.Id,
			Name: headerImage.Name,
			Path: headerImage.Path,
		}
	}

	var categoryFormat any
	if category != nil {
		categoryFormat = struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		}{
			Id:   category.Id,
			Name: category.Name,
		}
	}

	return json.Marshal(struct {
		Id          int    `json:"id"`
		Title       string `json:"title"`
		Slug        string `json:"slug"`
		HeaderImage any    `json:"headerImage"`
		Category    any    `json:"category"`
		Public      bool   `json:"public"`
		Created     any    `json:"created"`
		Updated     any    `json:"updated"`
	}{
		Id:          post.Id,
		Title:       post.Title,
		Slug:        post.Slug,
		HeaderImage: headerFormat,
		Category:    categoryFormat,
		Public:      post.Public,
		Created: struct {
			At string `json:"at"`
			By any    `json:"by"`
		}{
			At: post.CreatedAt.Format(time.RFC3339),
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
			At: post.LastUpdatedAt.Format(time.RFC3339),
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

func (post *BlogPost) Create() error {
	post.CreatedAt = time.Now()
	post.LastUpdatedAt = time.Now()
	return Insert(post)
}

func (post *BlogPost) Update() error {
	post.LastUpdatedAt = time.Now()
	return Update(post)
}

func (post *BlogPost) Delete() error {
	return Delete(post)
}

func FindBlogPostById(id int) (*BlogPost, error) {
	return Get[BlogPost](id)
}

func FindBlogPostBySlug(slug string) (*BlogPost, error) {
	return SelectOne[BlogPost](`
select id,
       title,
       slug,
       header_image_id,
       public,
       created_at,
       creator_id,
       updated_by_id,
       last_updated_at,
       category_id
from blog_post
where slug = $1
`, slug)
}

func FindAllBlogPosts() ([]BlogPost, error) {
	return Select[BlogPost](`
select id,
       title,
       slug,
       header_image_id,
       public,
       created_at,
       creator_id,
       updated_by_id,
       last_updated_at,
       category_id
from blog_post
`)
}

func FindPublicBlogPosts() ([]BlogPost, error) {
	return Select[BlogPost](`
select id,
       title,
       slug,
       header_image_id,
       public,
       created_at,
       creator_id,
       updated_by_id,
       last_updated_at,
       category_id
from blog_post
where public is true
`)
}

func (post *BlogPost) GetCategory() (*BlogCategory, error) {
	if post.CategoryId == nil {
		return nil, nil
	}
	return FindBlogCategoryById(*post.CategoryId)
}

func (post *BlogPost) GetCreator() (*Artist, error) {
	return FindArtistById(post.CreatorId)
}

func (post *BlogPost) GetUpdatedBy() (*Artist, error) {
	return FindArtistById(post.UpdatedById)
}

func (post *BlogPost) GetHeaderImage() (*File, error) {
	if post.HeaderImageId == nil {
		return nil, nil
	}
	return FindFileById(*post.HeaderImageId)
}

func (post *BlogPost) GetSections() ([]BlogPostSection, error) {
	return FindBlogPostSectionsByPost(post.Id)
}

func (post *BlogPost) ReplaceSections(sections []BlogPostSection) error {
	tx, err := dbMap.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.Exec(`
delete
from blog_post_section
where blog_post_id = $1`, post.Id)
	if err != nil {
		return err
	}

	for i, section := range sections {
		section.BlogPostId = post.Id
		section.Position = i
		err = tx.Insert(&section)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
