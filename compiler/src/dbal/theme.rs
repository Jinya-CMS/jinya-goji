use crate::dbal::{get_database, Error};
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
    type Error = Error;

    fn try_into(self) -> Result<Theme, Self::Error> {
        let variables = serde_json::from_str(&self.scss_variables).map_err(|err| Error {
            message: err.to_string(),
        })?;

        Ok(Theme {
            name: self.name,
            scss_variables: variables,
        })
    }
}
pub async fn get_theme(theme_id: i32) -> Result<Theme, Error> {
    let db = get_database().await.map_err(|err| Error {
        message: err.to_string(),
    })?;

    Entity::find_by_id(theme_id)
        .one(&db)
        .await
        .map_err(|err| Error {
            message: err.to_string(),
        })
        .map(|res| {
            if let Some(m) = res {
                m.try_into()
            } else {
                Err(Error {
                    message: "theme not found".to_string(),
                })
            }
        })?
}
