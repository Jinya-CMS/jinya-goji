use sea_orm::entity::prelude::*;

#[derive(Clone, Debug, PartialEq, DeriveEntityModel, Eq)]
#[sea_orm(table_name = "user")]
pub struct Model {
    #[sea_orm(primary_key)]
    pub id: i32,
    #[sea_orm(column_type = "Text", nullable, unique)]
    pub email: Option<String>,
    #[sea_orm(column_type = "Text", nullable)]
    pub artist_name: Option<String>,
    #[sea_orm(column_type = "Text", nullable)]
    pub profile_picture: Option<String>,
    #[sea_orm(column_type = "Text", nullable)]
    pub about_me: Option<String>,
}

#[derive(Copy, Clone, Debug, EnumIter, DeriveRelation)]
pub enum Relation {
    #[sea_orm(has_many = "super::menu_item::Entity")]
    MenuItem,
}

impl Related<super::menu_item::Entity> for Entity {
    fn to() -> RelationDef {
        Relation::MenuItem.def()
    }
}

impl ActiveModelBehavior for ActiveModel {}
