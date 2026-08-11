use crate::http::{error_route, frontend_router};
use actix_web::{App, HttpServer, web};
use env_logger::Env;

pub mod dbal;
pub mod http;
pub mod models;
pub mod theming;
pub mod preview;

#[tokio::main]
async fn main() -> std::io::Result<()> {
    let env =
        Env::new().default_filter_or(std::env::var("LOG_LEVEL").unwrap_or("info".to_string()));

    env_logger::init_from_env(env);

    let uds_path = std::env::args().nth(1).expect("usage: engine <uds-path>");

    HttpServer::new(|| {
        App::new()
            .route("/internal/error-handler", web::to(error_route))
            .default_service(web::to(frontend_router))
    })
    .bind_uds(uds_path)?
    .run()
    .await
}
