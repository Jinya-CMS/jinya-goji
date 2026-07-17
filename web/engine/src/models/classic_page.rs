use sea_orm::entity::prelude::*;
use serde::{Deserialize, Serialize};

#[derive(Clone, Debug, PartialEq, DeriveEntityModel, Eq, Serialize, Deserialize)]
#[sea_orm(table_name = "classic_page")]
pub struct Model {
    #[sea_orm(primary_key)]
    pub id: i32,
    pub creator_id: Option<i32>,
    pub updated_by_id: Option<i32>,
    pub created_at: Option<DateTimeWithTimeZone>,
    pub last_updated_at: Option<DateTimeWithTimeZone>,
    #[sea_orm(column_type = "Text", nullable)]
    pub content: Option<String>,
    #[sea_orm(column_type = "Text", nullable, unique)]
    pub title: Option<String>,
}

#[derive(Copy, Clone, Debug, EnumIter, DeriveRelation)]
pub enum Relation {
    #[sea_orm(has_many = "super::menu_item::Entity")]
    MenuItem,
    #[sea_orm(has_many = "super::theme_classic_page::Entity")]
    ThemeClassicPage,
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

impl Related<super::theme_classic_page::Entity> for Entity {
    fn to() -> RelationDef {
        Relation::ThemeClassicPage.def()
    }
}

impl Related<super::theme::Entity> for Entity {
    fn to() -> RelationDef {
        super::theme_classic_page::Relation::Theme.def()
    }
    fn via() -> Option<RelationDef> {
        Some(super::theme_classic_page::Relation::ClassicPage.def().rev())
    }
}

impl ActiveModelBehavior for ActiveModel {}
