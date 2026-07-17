use sea_orm::entity::prelude::*;

#[derive(Clone, Debug, PartialEq, DeriveEntityModel, Eq)]
#[sea_orm(table_name = "theme_menu")]
pub struct Model {
    #[sea_orm(primary_key, auto_increment = false)]
    pub theme_id: i32,
    #[sea_orm(primary_key, auto_increment = false, column_type = "Text")]
    pub name: String,
    pub menu_id: Option<i32>,
}

#[derive(Copy, Clone, Debug, EnumIter, DeriveRelation)]
pub enum Relation {
    #[sea_orm(
        belongs_to = "super::menu::Entity",
        from = "Column::MenuId",
        to = "super::menu::Column::Id",
        on_update = "NoAction",
        on_delete = "Cascade"
    )]
    Menu,
    #[sea_orm(
        belongs_to = "super::theme::Entity",
        from = "Column::ThemeId",
        to = "super::theme::Column::Id",
        on_update = "NoAction",
        on_delete = "Cascade"
    )]
    Theme,
}

impl Related<super::menu::Entity> for Entity {
    fn to() -> RelationDef {
        Relation::Menu.def()
    }
}

impl Related<super::theme::Entity> for Entity {
    fn to() -> RelationDef {
        Relation::Theme.def()
    }
}

impl ActiveModelBehavior for ActiveModel {}
