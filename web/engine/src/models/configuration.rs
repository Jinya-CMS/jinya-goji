use sea_orm::entity::prelude::*;
use serde::{Deserialize, Serialize};

#[derive(Clone, Debug, PartialEq, DeriveEntityModel, Eq, Serialize, Deserialize)]
#[sea_orm(table_name = "configuration")]
pub struct Model {
    #[sea_orm(primary_key)]
    pub id: i32,
    #[sea_orm(unique)]
    pub current_frontend_theme_id: Option<i32>,
}

#[derive(Copy, Clone, Debug, EnumIter, DeriveRelation)]
pub enum Relation {
    #[sea_orm(
        belongs_to = "super::theme::Entity",
        from = "Column::CurrentFrontendThemeId",
        to = "super::theme::Column::Id",
        on_update = "NoAction",
        on_delete = "NoAction"
    )]
    Theme,
}

impl Related<super::theme::Entity> for Entity {
    fn to() -> RelationDef {
        Relation::Theme.def()
    }
}

impl ActiveModelBehavior for ActiveModel {}
