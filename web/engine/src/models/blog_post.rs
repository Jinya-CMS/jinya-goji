use sea_orm::entity::prelude::*;
use serde::{Deserialize, Serialize};

#[derive(Clone, Debug, PartialEq, DeriveEntityModel, Eq, Serialize, Deserialize)]
#[sea_orm(table_name = "blog_post")]
pub struct Model {
    #[sea_orm(primary_key)]
    pub id: i32,
    #[sea_orm(column_type = "Text", unique)]
    pub title: String,
    #[sea_orm(column_type = "Text", unique)]
    pub slug: String,
    pub header_image_id: Option<i32>,
    pub public: bool,
    pub created_at: DateTimeWithTimeZone,
    pub creator_id: i32,
    pub updated_by_id: i32,
    pub last_updated_at: DateTimeWithTimeZone,
    pub category_id: Option<i32>,
}

#[derive(Copy, Clone, Debug, EnumIter, DeriveRelation)]
pub enum Relation {
    #[sea_orm(
        belongs_to = "super::blog_category::Entity",
        from = "Column::CategoryId",
        to = "super::blog_category::Column::Id",
        on_update = "NoAction",
        on_delete = "Cascade"
    )]
    BlogCategory,
    #[sea_orm(has_many = "super::blog_post_section::Entity")]
    BlogPostSection,
    #[sea_orm(
        belongs_to = "super::file::Entity",
        from = "Column::HeaderImageId",
        to = "super::file::Column::Id",
        on_update = "NoAction",
        on_delete = "Cascade"
    )]
    File,
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

impl Related<super::blog_category::Entity> for Entity {
    fn to() -> RelationDef {
        Relation::BlogCategory.def()
    }
}

impl Related<super::blog_post_section::Entity> for Entity {
    fn to() -> RelationDef {
        Relation::BlogPostSection.def()
    }
}

impl Related<super::file::Entity> for Entity {
    fn to() -> RelationDef {
        Relation::File.def()
    }
}

impl ActiveModelBehavior for ActiveModel {}
