use sea_orm::entity::prelude::*;
use serde::{Deserialize, Serialize};

#[derive(Clone, Debug, PartialEq, DeriveEntityModel, Eq, Serialize, Deserialize)]
#[sea_orm(table_name = "menu")]
#[serde(rename_all = "camelCase")]
pub struct Model {
    #[sea_orm(primary_key)]
    pub id: i32,
    #[sea_orm(column_type = "Text")]
    pub name: String,
    pub logo: Option<i32>,
}

#[derive(Copy, Clone, Debug, EnumIter, DeriveRelation)]
pub enum Relation {
    #[sea_orm(
        belongs_to = "super::file::Entity",
        from = "Column::Logo",
        to = "super::file::Column::Id",
        on_update = "NoAction",
        on_delete = "Cascade"
    )]
    File,
    #[sea_orm(has_many = "super::menu_item::Entity")]
    MenuItem,
    #[sea_orm(has_many = "super::theme_menu::Entity")]
    ThemeMenu,
}

impl Related<super::file::Entity> for Entity {
    fn to() -> RelationDef {
        Relation::File.def()
    }
}

impl Related<super::menu_item::Entity> for Entity {
    fn to() -> RelationDef {
        Relation::MenuItem.def()
    }
}

impl Related<super::theme_menu::Entity> for Entity {
    fn to() -> RelationDef {
        Relation::ThemeMenu.def()
    }
}

impl Related<super::theme::Entity> for Entity {
    fn to() -> RelationDef {
        super::theme_menu::Relation::Theme.def()
    }
    fn via() -> Option<RelationDef> {
        Some(super::theme_menu::Relation::Menu.def().rev())
    }
}

impl ActiveModelBehavior for ActiveModel {}

impl From<Model> for crate::theming::Menu {
    fn from(value: Model) -> Self {
        Self {
            id: value.id,
            name: value.name,
            logo: value.logo,
        }
    }
}
