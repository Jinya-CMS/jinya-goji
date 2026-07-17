use sea_orm::ColumnTrait;
use sea_orm::QueryFilter;
use crate::dbal::db;
use crate::models;
use crate::models::BlogPost;
use sea_orm::EntityTrait;

pub async fn find_by_id(id: i32) -> Option<BlogPost> {
    models::blog_post::Entity::find_by_id(id)
        .one(db().await.ok()?)
        .await
        .ok()?
}

pub async fn find_by_slug(slug: &str) -> Option<BlogPost> {
    models::blog_post::Entity::find()
        .filter(models::blog_post::Column::Slug.eq(slug))
        .one(db().await.ok()?)
        .await
        .ok()?
}
