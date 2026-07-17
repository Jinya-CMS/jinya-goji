use anyhow::anyhow;
use sea_orm::entity::prelude::*;
use std::collections::HashMap;
use sea_orm::entity::prelude::*;
use serde::{Deserialize, Serialize};

#[derive(Clone, Debug, PartialEq, DeriveEntityModel, Eq, Serialize, Deserialize)]
#[sea_orm(table_name = "theme")]
pub struct Model {
    #[sea_orm(primary_key)]
    pub id: i32,
    pub name: String,
    pub path: String,
    pub scss_variables: String,
    pub configuration: String,
}

#[derive(Copy, Clone, Debug, EnumIter, DeriveRelation)]
pub enum Relation {
    #[sea_orm(
        belongs_to = "crate::models::configuration::Entity",
        from = "crate::models::theme::Column::Id",
        to = "crate::models::configuration::Column::CurrentFrontendThemeId"
    )]
    Configuration,
}

impl ActiveModelBehavior for ActiveModel {}

#[derive(Debug, Clone, Default, Deserialize, Serialize)]
pub struct Theme {
    pub id: i32,
    pub name: String,
    pub path: String,
    pub scss_variables: HashMap<String, String>,
    pub configuration: HashMap<String, HashMap<String, String>>,
}

impl TryInto<Theme> for Model {
    type Error = anyhow::Error;

    fn try_into(self) -> Result<Theme, Self::Error> {
        let variables = serde_json::from_str(&self.scss_variables).map_err(|err| anyhow!(err))?;
        let configuration =
            serde_json::from_str(&self.configuration).map_err(|err| anyhow!(err))?;

        Ok(Theme {
            id: self.id,
            name: self.name,
            path: self.path,
            scss_variables: variables,
            configuration,
        })
    }
}