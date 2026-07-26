use sea_orm::entity::prelude::*;
use serde::{Deserialize, Serialize};

#[derive(Clone, Debug, PartialEq, DeriveEntityModel, Eq, Serialize, Deserialize)]
#[sea_orm(table_name = "file")]
pub struct Model {
    #[sea_orm(primary_key)]
    pub id: i32,
    pub creator_id: i32,
    pub updated_by_id: i32,
    pub created_at: DateTimeWithTimeZone,
    pub last_updated_at: DateTimeWithTimeZone,
    #[sea_orm(column_type = "Text")]
    pub path: String,
    #[sea_orm(column_type = "Text")]
    pub name: String,
    pub r#type: String,
}

#[derive(Copy, Clone, Debug, EnumIter, DeriveRelation)]
pub enum Relation {
    #[sea_orm(has_many = "super::blog_post::Entity")]
    BlogPost,
    #[sea_orm(has_many = "super::blog_post_section::Entity")]
    BlogPostSection,
    #[sea_orm(has_many = "super::gallery_file_position::Entity")]
    GalleryFilePosition,
    #[sea_orm(has_many = "super::menu::Entity")]
    Menu,
    #[sea_orm(has_many = "super::modern_page_section::Entity")]
    ModernPageSection,
    #[sea_orm(has_many = "super::theme_file::Entity")]
    ThemeFile,
    #[sea_orm(
        belongs_to = "super::user::Entity",
        from = "Column::CreatorId",
        to = "super::user::Column::Id",
        on_update = "NoAction",
        on_delete = "Cascade"
    )]
    Creator,
    #[sea_orm(
        belongs_to = "super::user::Entity",
        from = "Column::UpdatedById",
        to = "super::user::Column::Id",
        on_update = "NoAction",
        on_delete = "Cascade"
    )]
    UpdatedBy,
}

impl Related<super::blog_post::Entity> for Entity {
    fn to() -> RelationDef {
        Relation::BlogPost.def()
    }
}

impl Related<super::blog_post_section::Entity> for Entity {
    fn to() -> RelationDef {
        Relation::BlogPostSection.def()
    }
}

impl Related<super::gallery_file_position::Entity> for Entity {
    fn to() -> RelationDef {
        Relation::GalleryFilePosition.def()
    }
}

impl Related<super::menu::Entity> for Entity {
    fn to() -> RelationDef {
        Relation::Menu.def()
    }
}

impl Related<super::modern_page_section::Entity> for Entity {
    fn to() -> RelationDef {
        Relation::ModernPageSection.def()
    }
}

impl Related<super::theme_file::Entity> for Entity {
    fn to() -> RelationDef {
        Relation::ThemeFile.def()
    }
}

impl Related<super::theme::Entity> for Entity {
    fn to() -> RelationDef {
        super::theme_file::Relation::Theme.def()
    }
    fn via() -> Option<RelationDef> {
        Some(super::theme_file::Relation::File.def().rev())
    }
}

impl ActiveModelBehavior for ActiveModel {}
