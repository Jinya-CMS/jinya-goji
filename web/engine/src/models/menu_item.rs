use sea_orm::entity::prelude::*;
use serde::{Deserialize, Serialize};

#[derive(Clone, Debug, PartialEq, DeriveEntityModel, Eq, Serialize, Deserialize)]
#[sea_orm(table_name = "menu_item")]
#[serde(rename_all = "camelCase")]
pub struct Model {
    #[sea_orm(primary_key)]
    pub id: i32,
    pub position: i32,
    pub menu_id: Option<i32>,
    pub parent_id: Option<i32>,
    #[sea_orm(column_type = "Text")]
    pub title: String,
    #[sea_orm(column_type = "Text", nullable)]
    pub route: Option<String>,
    pub highlighted: Option<bool>,
    pub artist_id: Option<i32>,
    pub gallery_id: Option<i32>,
    pub modern_page_id: Option<i32>,
    pub classic_page_id: Option<i32>,
    pub blog_category_id: Option<i32>,
    pub blog_home_page: Option<bool>,
}

#[derive(Copy, Clone, Debug, EnumIter, DeriveRelation)]
pub enum Relation {
    #[sea_orm(
        belongs_to = "super::blog_category::Entity",
        from = "Column::BlogCategoryId",
        to = "super::blog_category::Column::Id",
        on_update = "NoAction",
        on_delete = "Cascade"
    )]
    BlogCategory,
    #[sea_orm(
        belongs_to = "super::classic_page::Entity",
        from = "Column::ClassicPageId",
        to = "super::classic_page::Column::Id",
        on_update = "NoAction",
        on_delete = "Cascade"
    )]
    ClassicPage,
    #[sea_orm(
        belongs_to = "super::gallery::Entity",
        from = "Column::GalleryId",
        to = "super::gallery::Column::Id",
        on_update = "NoAction",
        on_delete = "Cascade"
    )]
    Gallery,
    #[sea_orm(
        belongs_to = "super::menu::Entity",
        from = "Column::MenuId",
        to = "super::menu::Column::Id",
        on_update = "NoAction",
        on_delete = "Cascade"
    )]
    Menu,
    #[sea_orm(
        belongs_to = "Entity",
        from = "Column::ParentId",
        to = "Column::Id",
        on_update = "NoAction",
        on_delete = "Cascade"
    )]
    Parent,
    #[sea_orm(
        belongs_to = "super::modern_page::Entity",
        from = "Column::ModernPageId",
        to = "super::modern_page::Column::Id",
        on_update = "NoAction",
        on_delete = "Cascade"
    )]
    ModernPage,
    #[sea_orm(
        belongs_to = "super::user::Entity",
        from = "Column::ArtistId",
        to = "super::user::Column::Id",
        on_update = "NoAction",
        on_delete = "Cascade"
    )]
    User,
}

impl Related<super::blog_category::Entity> for Entity {
    fn to() -> RelationDef {
        Relation::BlogCategory.def()
    }
}

impl Related<super::classic_page::Entity> for Entity {
    fn to() -> RelationDef {
        Relation::ClassicPage.def()
    }
}

impl Related<super::gallery::Entity> for Entity {
    fn to() -> RelationDef {
        Relation::Gallery.def()
    }
}

impl Related<super::menu::Entity> for Entity {
    fn to() -> RelationDef {
        Relation::Menu.def()
    }
}

impl Related<super::modern_page::Entity> for Entity {
    fn to() -> RelationDef {
        Relation::ModernPage.def()
    }
}

impl Related<super::user::Entity> for Entity {
    fn to() -> RelationDef {
        Relation::User.def()
    }
}

impl ActiveModelBehavior for ActiveModel {}

impl From<Model> for crate::theming::MenuItem {
    fn from(value: Model) -> Self {
        Self {
            id: value.id,
            menu_id: value.menu_id,
            parent_id: value.parent_id,
            position: value.position,
            title: value.title,
            route: value.route,
            highlighted: value.highlighted,
            artist_id: value.artist_id,
            gallery_id: value.gallery_id,
            modern_page_id: value.modern_page_id,
            classic_page_id: value.classic_page_id,
            blog_category_id: value.blog_category_id,
            blog_home_page: value.blog_home_page,
        }
    }
}

impl From<crate::theming::MenuItem> for Model {
    fn from(val: crate::theming::MenuItem) -> Self {
        let route = val.route.map(|route| {
            if !route.starts_with("/") {
                format!("/{route}")
            } else {
                route
            }
        });

        Model {
            id: val.id,
            position: val.position,
            menu_id: None,
            parent_id: None,
            title: val.title,
            route,
            highlighted: val.highlighted,
            artist_id: val.artist_id,
            gallery_id: val.gallery_id,
            modern_page_id: val.modern_page_id,
            classic_page_id: val.classic_page_id,
            blog_category_id: val.blog_category_id,
            blog_home_page: val.blog_home_page,
        }
    }
}
