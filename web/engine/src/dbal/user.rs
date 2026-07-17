use crate::dbal::db;
use crate::models;
use crate::models::User;
use sea_orm::EntityTrait;

pub async fn find_by_id(id: i32) -> Option<User> {
    models::user::Entity::find_by_id(id)
        .one(db().await.ok()?)
        .await
        .ok()?
}
