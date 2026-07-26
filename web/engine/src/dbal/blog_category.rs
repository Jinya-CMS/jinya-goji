use crate::dbal::db;
use crate::models::blog_category::BlogCategoryPath;
use crate::{find_by_parent, find_by_theme, find_entity};
use sea_orm::{DbBackend, FromQueryResult, Statement};

find_entity!(blog_category, BlogCategory);

find_by_parent!(blog_post, BlogPost, CategoryId, find_posts);

find_by_theme!(
    blog_category,
    BlogCategory,
    theme_blog_category,
    ThemeBlogCategory,
    find_by_theme
);

pub async fn get_category_path(category_id: i32) -> Vec<BlogCategoryPath> {
    let query = r#"
with recursive category_path as (
    select bc.*, 0 as depth
    from blog_category bc
    where bc.id = $1

    union all

    -- recursive step: walk up to the parent
    select bc2.*, cp.depth + 1
    from category_path cp
             inner join blog_category bc2 on bc2.id = cp.parent_id)
select cp2.id as category_id, cp2.name as category_name, mi.id as menu_item_id, mi.route as route, mi.title as menu_item_title
from category_path cp2
         join menu_item mi on mi.category_id = cp2.id
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
