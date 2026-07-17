use crate::dbal::db;
use crate::models;
use crate::models::BlogCategory;
use sea_orm::EntityTrait;

pub async fn find_by_id(id: i32) -> Option<BlogCategory> {
    models::blog_category::Entity::find_by_id(id)
        .one(db().await.ok()?)
        .await
        .ok()?
}
