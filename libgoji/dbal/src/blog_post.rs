use crate::models::BlogPost;
use crate::models::blog_category::BlogCategoryPath;
use crate::{db, find_by_parent, models};
use sea_orm::{ColumnTrait, EntityTrait, FromQueryResult};
use sea_orm::{DbBackend, QueryFilter, Statement};

pub async fn find_by_id(id: i32) -> Option<BlogPost> {
    models::blog_post::Entity::find_by_id(id)
        .one(db().await.ok()?)
        .await
        .ok()?
}
pub async fn find_all() -> Option<Vec<BlogPost>> {
    models::blog_post::Entity::find()
        .filter(models::blog_post::Column::Public.eq(true))
        .all(db().await.ok()?)
        .await
        .ok()
}
pub async fn find_by_slug(slug: &str) -> Option<BlogPost> {
    models::blog_post::Entity::find()
        .filter(models::blog_post::Column::Slug.eq(slug))
        .one(db().await.ok()?)
        .await
        .ok()?
}

find_by_parent!(
    blog_post_section,
    BlogPostSection,
    BlogPostId,
    find_sections
);

pub async fn get_category_path(post_id: i32) -> Vec<BlogCategoryPath> {
    //language=postgresql
    let query = r#"
with recursive category_path as (select bc.*, 0 as depth
                                 from blog_category bc
                                 where bc.id = (select category_id from blog_post where id = $1)

                                 union all

                                 select bc2.*, cp.depth + 1
                                 from category_path cp
                                          inner join blog_category bc2 on bc2.id = cp.parent_id)
select cp2.id   as category_id,
       cp2.name as category_name,
       mi.id    as menu_item_id,
       mi.route as route,
       mi.title as menu_item_title
from category_path cp2
         join menu_item mi on mi.blog_category_id = cp2.id
order by depth desc;
"#;

    if let Ok(db) = db().await {
        BlogCategoryPath::find_by_statement(Statement::from_sql_and_values(
            DbBackend::Postgres,
            query,
            vec![post_id.into()],
        ))
        .all(db)
        .await
        .unwrap_or_default()
    } else {
        vec![]
    }
}
