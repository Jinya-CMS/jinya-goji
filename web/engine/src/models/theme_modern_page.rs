use sea_orm::entity::prelude::*;
use serde::{Deserialize, Serialize};

#[derive(Clone, Debug, PartialEq, DeriveEntityModel, Eq, Serialize, Deserialize)]
#[sea_orm(table_name = "theme_modern_page")]
pub struct Model {
    #[sea_orm(primary_key, auto_increment = false)]
    pub theme_id: i32,
    #[sea_orm(primary_key, auto_increment = false, column_type = "Text")]
    pub name: String,
    pub modern_page_id: Option<i32>,
}

#[derive(Copy, Clone, Debug, EnumIter, DeriveRelation)]
pub enum Relation {
    #[sea_orm(
        belongs_to = "super::modern_page::Entity",
        from = "Column::ModernPageId",
        to = "super::modern_page::Column::Id",
        on_update = "NoAction",
        on_delete = "Cascade"
    )]
    ModernPage,
    #[sea_orm(
        belongs_to = "super::theme::Entity",
        from = "Column::ThemeId",
        to = "super::theme::Column::Id",
        on_update = "NoAction",
        on_delete = "Cascade"
    )]
    Theme,
}

impl Related<super::modern_page::Entity> for Entity {
    fn to() -> RelationDef {
        Relation::ModernPage.def()
    }
}

impl Related<super::theme::Entity> for Entity {
    fn to() -> RelationDef {
        Relation::Theme.def()
    }
}

impl ActiveModelBehavior for ActiveModel {}
