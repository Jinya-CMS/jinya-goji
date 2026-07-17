use crate::dbal::db;
use crate::models::theme::{Entity, Relation, Theme};
use anyhow::anyhow;
use sea_orm::{EntityTrait, JoinType, QuerySelect, RelationTrait};

pub async fn get_active_theme() -> anyhow::Result<Theme> {
    let db = db().await.map_err(|err| anyhow!(err))?;

    Entity::find()
        .join(JoinType::InnerJoin, Relation::Configuration.def())
        .one(db)
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
