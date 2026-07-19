use sea_orm::entity::prelude::*;
use sea_orm::FromQueryResult;
use serde::{Deserialize, Serialize};

#[derive(Clone, Debug, PartialEq, DeriveEntityModel, Eq, Serialize, Deserialize)]
#[sea_orm(table_name = "gallery_file_position")]
pub struct Model {
    #[sea_orm(primary_key)]
    pub id: i32,
    pub position: i32,
    pub gallery_id: i32,
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
        belongs_to = "super::gallery::Entity",
        from = "Column::GalleryId",
        to = "super::gallery::Column::Id",
        on_update = "NoAction",
        on_delete = "Cascade"
    )]
    Gallery,
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

impl ActiveModelBehavior for ActiveModel {}

#[derive(FromQueryResult)]
pub struct GalleryFileWithFile {
    pub id: i32,
    pub path: String,
    pub name: String,
    pub position: i32,
}