mod theme;

use sea_orm::entity::prelude::*;
use std::fmt::{Debug, Display, Formatter};

pub use theme::Theme;
pub use theme::get_theme;

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

async fn get_database() -> Result<DatabaseConnection, DbErr> {
    sea_orm::Database::connect(std::env::var("DATABASE_URL").unwrap()).await
}
