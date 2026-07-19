use crate::dbal::db;
use crate::models::BlogPost;
use crate::{find_by_parent, find_entity, models};
use sea_orm::QueryFilter;
use sea_orm::{ColumnTrait, EntityTrait};

find_entity!(blog_post, BlogPost);

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
