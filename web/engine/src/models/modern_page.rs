use sea_orm::entity::prelude::*;
use serde::{Deserialize, Serialize};

#[derive(Clone, Debug, PartialEq, DeriveEntityModel, Eq, Serialize, Deserialize)]
#[sea_orm(table_name = "modern_page")]
#[serde(rename_all = "camelCase")]
pub struct Model {
    #[sea_orm(primary_key)]
    pub id: i32,
    pub creator_id: i32,
    pub updated_by_id: i32,
    pub created_at: DateTimeWithTimeZone,
    pub last_updated_at: DateTimeWithTimeZone,
    #[sea_orm(column_type = "Text", unique)]
    pub name: String,
}

#[derive(Copy, Clone, Debug, EnumIter, DeriveRelation)]
pub enum Relation {
    #[sea_orm(has_many = "super::menu_item::Entity")]
    MenuItem,
    #[sea_orm(has_many = "super::modern_page_section::Entity")]
    ModernPageSection,
    #[sea_orm(has_many = "super::theme_modern_page::Entity")]
    ThemeModernPage,
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

impl Related<super::menu_item::Entity> for Entity {
    fn to() -> RelationDef {
        Relation::MenuItem.def()
    }
}

impl Related<super::modern_page_section::Entity> for Entity {
    fn to() -> RelationDef {
        Relation::ModernPageSection.def()
    }
}

impl Related<super::theme_modern_page::Entity> for Entity {
    fn to() -> RelationDef {
        Relation::ThemeModernPage.def()
    }
}

impl Related<super::theme::Entity> for Entity {
    fn to() -> RelationDef {
        super::theme_modern_page::Relation::Theme.def()
    }
    fn via() -> Option<RelationDef> {
        Some(super::theme_modern_page::Relation::ModernPage.def().rev())
    }
}

impl ActiveModelBehavior for ActiveModel {}

impl From<Model> for crate::theming::ModernPage {
    fn from(value: Model) -> Self {
        Self {
            id: value.id,
            creator_id: value.creator_id,
            created_at: value.created_at.to_rfc3339(),
            updated_by_id: value.updated_by_id,
            last_updated_at: value.last_updated_at.to_rfc3339(),
            name: value.name,
        }
    }
}
