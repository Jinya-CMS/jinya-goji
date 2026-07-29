use crate::dbal::db;
use crate::models::BlogPost;
use crate::models::blog_category::BlogCategoryPath;
use crate::{find_by_theme, find_entity, models};
use sea_orm::ColumnTrait;
use sea_orm::QueryFilter;
use sea_orm::{DbBackend, EntityTrait, FromQueryResult, Statement};

find_entity!(blog_category, BlogCategory);

pub async fn find_posts(id: i32, with_children: bool) -> Option<Vec<BlogPost>> {
    if with_children {
        //language=postgresql
        let query = r#"
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
  and public is true
order by id
        "#;

     BlogPost::find_by_statement(Statement::from_sql_and_values(
            DbBackend::Postgres,
            query,
            vec![id.into()],
        ))
            .all(db().await.ok()?)
            .await
            .ok()
    } else {
        models::blog_post::Entity::find()
            .filter(models::blog_post::Column::CategoryId.eq(id))
            .all(db().await.ok()?)
            .await
            .ok()
    }
}

find_by_theme!(
    blog_category,
    BlogCategory,
    theme_blog_category,
    ThemeBlogCategory,
    find_by_theme
);

pub async fn get_category_path(category_id: i32) -> Vec<BlogCategoryPath> {
    //language=postgresql
    let query = r#"
with recursive category_path as (select bc.*, 0 as depth
                                 from blog_category bc
                                 where bc.id = $1

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
            vec![category_id.into()],
        ))
        .all(db)
        .await
        .unwrap_or_default()
    } else {
        vec![]
    }
}
