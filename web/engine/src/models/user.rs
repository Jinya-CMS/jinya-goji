use crate::theming;
use sea_orm::entity::prelude::*;

#[derive(Clone, Debug, PartialEq, DeriveEntityModel, Eq)]
#[sea_orm(table_name = "users")]
pub struct Model {
    #[sea_orm(primary_key)]
    pub id: i32,
    #[sea_orm(column_type = "Text", unique)]
    pub email: String,
    #[sea_orm(column_type = "Text")]
    pub artist_name: String,
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

impl From<Model> for theming::User {
    fn from(user: Model) -> Self {
        Self {
            email: user.email,
            artist_name: user.artist_name,
            profile_picture: user.profile_picture,
            about_me: user.about_me,
        }
    }
}
