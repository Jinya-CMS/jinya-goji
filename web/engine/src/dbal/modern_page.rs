use crate::dbal::db;
use crate::models;
use crate::models::ModernPage;
use sea_orm::EntityTrait;

pub async fn find_by_id(id: i32) -> Option<ModernPage> {
    models::modern_page::Entity::find_by_id(id)
        .one(db().await.ok()?)
        .await
        .ok()?
}
