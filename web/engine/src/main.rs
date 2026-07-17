use crate::http::frontend_router;
use actix_web::{web, App, HttpServer};
use env_logger::Env;

pub mod dbal;
pub mod http;
pub mod models;
pub mod theming;

#[tokio::main]
async fn main() -> std::io::Result<()> {
    let env =
        Env::new().default_filter_or(std::env::var("LOG_LEVEL").unwrap_or("info".to_string()));

    env_logger::init_from_env(env);

    let uds_path = std::env::args().nth(1).expect("usage: engine <uds-path>");

    HttpServer::new(|| App::new().default_service(web::to(frontend_router)))
        .bind_uds(uds_path)?
        .run()
        .await
}
