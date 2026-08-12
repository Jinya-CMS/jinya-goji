use crate::db;
use crate::models::theme::{Entity, Relation, Theme};
use anyhow::anyhow;
use sea_orm::{EntityTrait, JoinType, QuerySelect, RelationTrait};

pub async fn get_active_theme() -> anyhow::Result<Theme> {
    Entity::find()
        .join(JoinType::InnerJoin, Relation::Configuration.def())
        .one(db().await?)
        .await
        .map_err(|err| anyhow!(err))
        .map(|res| {
            if let Some(m) = res {
                m.try_into()
            } else {
                Err(anyhow!("theme not found"))
            }
        })?
}

pub async fn find_theme_by_id(id: i32) -> anyhow::Result<Theme> {
    Entity::find_by_id(id)
        .one(db().await?)
        .await
        .map_err(|err| anyhow!(err))
        .map(|res| {
            if let Some(m) = res {
                m.try_into()
            } else {
                Err(anyhow!("theme not found"))
            }
        })?
}
