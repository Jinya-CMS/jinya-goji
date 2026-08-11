use crate::dbal::db;
use anyhow::anyhow;
use sea_orm::entity::prelude::*;
use std::collections::HashMap;

#[derive(Clone, Debug, PartialEq, Eq, DeriveEntityModel)]
#[sea_orm(table_name = "theme")]
pub struct Model {
    #[sea_orm(primary_key)]
    pub id: i32,
    pub name: String,
    pub scss_variables: String,
}

#[derive(Copy, Clone, Debug, EnumIter, DeriveRelation)]
pub enum Relation {}

impl ActiveModelBehavior for ActiveModel {}

#[derive(Debug, Clone)]
pub struct Theme {
    pub name: String,
    pub scss_variables: HashMap<String, String>,
}

impl TryInto<Theme> for Model {
    type Error = anyhow::Error;

    fn try_into(self) -> Result<Theme, Self::Error> {
        let variables =
            serde_json::from_str(&self.scss_variables).map_err(|err| anyhow!(err))?;

        Ok(Theme {
            name: self.name,
            scss_variables: variables,
        })
    }
}
pub async fn get_theme(theme_id: i32) -> anyhow::Result<Theme> {
    let db = db().await.map_err(|err| anyhow!(err))?;

    Entity::find_by_id(theme_id)
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
