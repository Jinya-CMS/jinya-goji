use sea_orm::FromQueryResult;
use sea_orm::entity::prelude::*;
use serde::{Deserialize, Serialize};
use crate::theming::jinya::goji::blog;

#[derive(Clone, Debug, PartialEq, DeriveEntityModel, Eq, Serialize, Deserialize)]
#[sea_orm(table_name = "blog_category")]
#[serde(rename_all = "camelCase")]
pub struct Model {
    #[sea_orm(primary_key)]
    pub id: i32,
    #[sea_orm(column_type = "Text")]
    pub name: String,
    #[sea_orm(column_type = "Text")]
    pub description: String,
    pub parent_id: Option<i32>,
}

#[derive(Copy, Clone, Debug, EnumIter, DeriveRelation)]
pub enum Relation {
    #[sea_orm(
        belongs_to = "Entity",
        from = "Column::ParentId",
        to = "Column::Id",
        on_update = "NoAction",
        on_delete = "Cascade"
    )]
    SelfRef,
    #[sea_orm(has_many = "super::blog_post::Entity")]
    BlogPost,
    #[sea_orm(has_many = "super::menu_item::Entity")]
    MenuItem,
    #[sea_orm(has_many = "super::theme_blog_category::Entity")]
    ThemeBlogCategory,
}

impl Related<super::blog_post::Entity> for Entity {
    fn to() -> RelationDef {
        Relation::BlogPost.def()
    }
}

impl Related<super::menu_item::Entity> for Entity {
    fn to() -> RelationDef {
        Relation::MenuItem.def()
    }
}

impl Related<super::theme_blog_category::Entity> for Entity {
    fn to() -> RelationDef {
        Relation::ThemeBlogCategory.def()
    }
}

impl Related<super::theme::Entity> for Entity {
    fn to() -> RelationDef {
        super::theme_blog_category::Relation::Theme.def()
    }
    fn via() -> Option<RelationDef> {
        Some(
            super::theme_blog_category::Relation::BlogCategory
                .def()
                .rev(),
        )
    }
}

impl ActiveModelBehavior for ActiveModel {}

#[derive(FromQueryResult)]
pub struct BlogCategoryPath {
    pub category_id: i32,
    pub category_name: String,
    pub menu_item_id: i32,
    pub route: String,
    pub menu_item_title: String,
}

impl From<Model> for crate::theming::BlogCategory {
    fn from(category: Model) -> Self {
        Self {
            id: category.id,
            name: category.name,
            description: category.description,
        }
    }
}

impl From<BlogCategoryPath> for blog::BlogCategoryPathEntry {
    fn from(value: BlogCategoryPath) -> Self {
        Self {
            route: value.route,
            category_id: value.category_id,
            category_title: value.category_name,
            menu_item_id: value.menu_item_id,
            menu_item_title: value.menu_item_title,
        }
    }
}
