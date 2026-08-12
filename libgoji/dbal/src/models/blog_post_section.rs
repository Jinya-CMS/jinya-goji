use sea_orm::entity::prelude::*;
use serde::{Deserialize, Serialize};
use crate::models::{blog_post_section, BlogPostSection};

#[derive(Clone, Debug, PartialEq, DeriveEntityModel, Eq, Serialize, Deserialize)]
#[sea_orm(table_name = "blog_post_section")]
#[serde(rename_all = "camelCase")]
pub struct Model {
    #[sea_orm(primary_key)]
    pub id: i32,
    pub gallery_id: Option<i32>,
    pub file_id: Option<i32>,
    #[sea_orm(column_type = "Text", nullable)]
    pub html: Option<String>,
    #[sea_orm(column_type = "Text", nullable)]
    pub link: Option<String>,
    pub blog_post_id: Option<i32>,
    pub position: i32,
}

#[derive(Copy, Clone, Debug, EnumIter, DeriveRelation)]
pub enum Relation {
    #[sea_orm(
        belongs_to = "super::blog_post::Entity",
        from = "Column::BlogPostId",
        to = "super::blog_post::Column::Id",
        on_update = "NoAction",
        on_delete = "Cascade"
    )]
    BlogPost,
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

impl Related<super::blog_post::Entity> for Entity {
    fn to() -> RelationDef {
        Relation::BlogPost.def()
    }
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

impl From<Model> for wasi_bridge::BlogPostSection {
    fn from(value: Model) -> Self {
        Self {
            gallery_id: value.gallery_id,
            file_id: value.file_id,
            html: value.html,
            link: value.link,
            position: value.position,
        }
    }
}
