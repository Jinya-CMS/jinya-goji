package database

import "encoding/json"

type BlogCategory struct {
	Id             int     `db:"id,primarykey,autoincrement"`
	Name           string  `db:"name"`
	Description    *string `db:"description"`
	ParentId       *int    `db:"parent_id"`
	WebhookEnabled bool    `db:"webhook_enabled"`
	WebhookUrl     *string `db:"webhook_url"`
}

func (category *BlogCategory) MarshalJSON() ([]byte, error) {
	parent, err := category.GetParent()
	if err != nil {
		return nil, err
	}
	return json.Marshal(struct {
		Id             int           `json:"id"`
		Name           string        `json:"name"`
		Description    *string       `json:"description"`
		Parent         *BlogCategory `json:"parent"`
		WebhookEnabled bool          `json:"webhookEnabled"`
		WebhookUrl     *string       `json:"webhookUrl"`
	}{
		Id:             category.Id,
		Name:           category.Name,
		Description:    category.Description,
		Parent:         parent,
		WebhookEnabled: category.WebhookEnabled,
		WebhookUrl:     category.WebhookUrl,
	})
}

func (category *BlogCategory) Create() error {
	return Insert(category)
}

func (category *BlogCategory) Update() error {
	return Update(category)
}

func (category *BlogCategory) Delete() error {
	return Delete(category)
}

func FindBlogCategoryById(id int) (*BlogCategory, error) {
	return Get[BlogCategory](id)
}

func FindAllBlogCategories() ([]BlogCategory, error) {
	return Select[BlogCategory](`
select id,
       name,
       description,
       parent_id,
       webhook_enabled,
       webhook_url
from blog_category
order by id
`)
}

func (category *BlogCategory) GetParent() (*BlogCategory, error) {
	if category.ParentId == nil {
		return nil, nil
	}
	return FindBlogCategoryById(*category.ParentId)
}

func (category *BlogCategory) GetBlogPosts(includeChildCategories bool, onlyPublic bool) ([]BlogPost, error) {
	var query string
	if includeChildCategories {
		query = `
with recursive category_cte as (select bc.*
                                from blog_category as bc
                                where bc.id = $1
                                union
                                select bc2.*
                                from category_cte
                                         inner join blog_category as bc2 on category_cte.id = bc2.parent_id)
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
where category_id in (select id from category_cte)
order by id`
	} else {
		query = `
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
where category_id = $1
order by id
`
	}

	if onlyPublic {
		query += " and public = true"
	}

	return Select[BlogPost](query, category.Id)
}
