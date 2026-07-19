use sea_orm::entity::prelude::*;
use serde::{Deserialize, Serialize};

#[derive(Clone, Debug, PartialEq, DeriveEntityModel, Eq, Serialize, Deserialize)]
#[sea_orm(table_name = "theme_file")]
pub struct Model {
    #[sea_orm(primary_key, auto_increment = false)]
    pub theme_id: i32,
    #[sea_orm(primary_key, auto_increment = false, column_type = "Text")]
    pub name: String,
    pub file_id: i32,
}

#[derive(Copy, Clone, Debug, EnumIter, DeriveRelation)]
pub enum Relation {
    #[sea_orm(
        belongs_to = "super::file::Entity",
        from = "Column::FileId",
        to = "super::file::Column::Id",
        on_update = "NoAction",
        on_delete = "Cascade"
    )]
    File,
    #[sea_orm(
        belongs_to = "super::theme::Entity",
        from = "Column::ThemeId",
        to = "super::theme::Column::Id",
        on_update = "NoAction",
        on_delete = "Cascade"
    )]
    Theme,
}

impl Related<super::file::Entity> for Entity {
    fn to() -> RelationDef {
        Relation::File.def()
    }
}

impl Related<super::theme::Entity> for Entity {
    fn to() -> RelationDef {
        Relation::Theme.def()
    }
}

impl ActiveModelBehavior for ActiveModel {}
