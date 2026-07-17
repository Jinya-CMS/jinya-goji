pub mod blog_category;
pub mod blog_post;
pub mod classic_page;
pub mod gallery;
pub mod modern_page;
pub mod theme;
pub mod user;

use anyhow::anyhow;
use sea_orm::entity::prelude::*;
use sea_orm::Database;
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
