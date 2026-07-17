use crate::dbal::db;
use crate::models;
use crate::models::ClassicPage;
use sea_orm::EntityTrait;

pub async fn find_by_id(id: i32) -> Option<ClassicPage> {
    models::classic_page::Entity::find_by_id(id)
        .one(db().await.ok()?)
        .await
        .ok()?
}
