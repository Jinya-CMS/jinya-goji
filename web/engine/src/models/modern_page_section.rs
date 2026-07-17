use sea_orm::entity::prelude::*;
use serde::{Deserialize, Serialize};

#[derive(Clone, Debug, PartialEq, DeriveEntityModel, Eq, Serialize, Deserialize)]
#[sea_orm(table_name = "modern_page_section")]
pub struct Model {
    #[sea_orm(primary_key)]
    pub id: i32,
    pub position: Option<i32>,
    pub gallery_id: Option<i32>,
    pub file_id: Option<i32>,
    #[sea_orm(column_type = "Text", nullable)]
    pub html: Option<String>,
    pub page_id: Option<i32>,
    #[sea_orm(column_type = "Text", nullable)]
    pub action: Option<String>,
    #[sea_orm(column_type = "Text", nullable)]
    pub script: Option<String>,
    #[sea_orm(column_type = "Text", nullable)]
    pub target: Option<String>,
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
        belongs_to = "super::gallery::Entity",
        from = "Column::GalleryId",
        to = "super::gallery::Column::Id",
        on_update = "NoAction",
        on_delete = "Cascade"
    )]
    Gallery,
    #[sea_orm(
        belongs_to = "super::modern_page::Entity",
        from = "Column::PageId",
        to = "super::modern_page::Column::Id",
        on_update = "NoAction",
        on_delete = "Cascade"
    )]
    ModernPage,
}

impl Related<super::file::Entity> for Entity {
    fn to() -> RelationDef {
        Relation::File.def()
    }
}

impl Related<super::gallery::Entity> for Entity {
    fn to() -> RelationDef {
        Relation::Gallery.def()
    }
}

impl Related<super::modern_page::Entity> for Entity {
    fn to() -> RelationDef {
        Relation::ModernPage.def()
    }
}

impl ActiveModelBehavior for ActiveModel {}
