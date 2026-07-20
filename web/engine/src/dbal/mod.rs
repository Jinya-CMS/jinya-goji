pub mod blog_category;
pub mod blog_post;
pub mod classic_page;
pub mod file;
pub mod gallery;
pub mod menu;
pub mod menu_item;
pub mod modern_page;
pub mod theme;
pub mod user;

use anyhow::anyhow;
use sea_orm::entity::prelude::*;
use sea_orm::{Database, DbBackend, Statement};
use sql_query_analyzer::query::{QueryType, SqlDialect};
use std::collections::HashMap;
use std::fmt::{Debug, Display, Formatter};

#[derive(Debug)]
pub struct Error {
    pub message: String,
}

impl Display for Error {
    fn fmt(&self, f: &mut Formatter<'_>) -> std::fmt::Result {
        write!(f, "{}", self.message)
    }
}

impl std::error::Error for Error {}

static DB: tokio::sync::OnceCell<DatabaseConnection> = tokio::sync::OnceCell::const_new();

pub async fn db() -> anyhow::Result<&'static DatabaseConnection> {
    DB.get_or_try_init(|| async {
        let database_url = std::env::var("DATABASE_URL")?;

        Database::connect(database_url).await.map_err(|err| {
            log::error!("Failed to connect to database: {err}");
            anyhow!(err)
        })
    })
    .await
}

#[macro_export]
macro_rules! find_entity {
    ($name:ident,$type:ident) => {
        pub async fn find_by_id(id: i32) -> Option<$crate::models::$type> {
            use sea_orm::EntityTrait;

            $crate::models::$name::Entity::find_by_id(id)
                .one($crate::dbal::db().await.ok()?)
                .await
                .ok()?
        }

        pub async fn find_all() -> Option<Vec<$crate::models::$type>> {
            use sea_orm::EntityTrait;

            $crate::models::$name::Entity::find()
                .all($crate::dbal::db().await.ok()?)
                .await
                .ok()
        }
    };
}

#[macro_export]
macro_rules! find_by_parent {
    ($name:ident,$type:ident,$parent_column:ident,$fn_name:ident) => {
        pub async fn $fn_name(id: i32) -> Option<Vec<$crate::models::$type>> {
            use sea_orm::{ColumnTrait, EntityTrait, QueryFilter};

            $crate::models::$name::Entity::find()
                .filter($crate::models::$name::Column::$parent_column.eq(id))
                .all($crate::dbal::db().await.ok()?)
                .await
                .ok()
        }
    };
}

#[macro_export]
macro_rules! find_by_theme {
    ($name:ident,$type:ident,$theme_type:ident,$fn_name:ident) => {
        pub async fn $fn_name(name: String) -> Option<$crate::models::$type> {
            use sea_orm::{ColumnTrait, EntityTrait, QueryFilter, QuerySelect, RelationTrait};

            $crate::models::$name::Entity::find()
                .join(
                    sea_orm::JoinType::InnerJoin,
                    $crate::models::$theme_type::Relation::$type.def(),
                )
                .join(
                    sea_orm::JoinType::InnerJoin,
                    $crate::models::$theme_type::Relation::Theme.def(),
                )
                .join(
                    sea_orm::JoinType::InnerJoin,
                    $crate::models::theme::Relation::Configuration.def(),
                )
                .filter($crate::models::$theme_type::Column::Name.eq(name))
                .one($crate::dbal::db().await.ok()?)
                .await
                .ok()?
        }
    };
}

pub async fn select(query: String) -> Vec<HashMap<String, String>> {
    let queries = sql_query_analyzer::query::parse_queries(&query, SqlDialect::PostgreSQL)
        .unwrap_or_default();
    if queries.is_empty() || queries.len() > 1 {
        return vec![];
    }

    let first = queries.first().unwrap();
    if first.query_type != QueryType::Select {
        return vec![];
    }

    if first.tables.iter().any(|table| {
        table == "api_key" || table == "known_device" || table == "jinya_configuration"
    }) {
        return vec![];
    }

    let stmt = Statement::from_string(DbBackend::Postgres, &query);

    if let Ok(db) = db().await
        && let Ok(rows) = db
            .query_all_raw(stmt)
            .await
        {
            return rows
                .into_iter()
                .map(|row| {
                    row.column_names()
                        .into_iter()
                        .map(|col| (col.clone(), row.try_get("", &col).unwrap_or_default()))
                        .collect()
                })
                .collect();
        }

    vec![]
}
