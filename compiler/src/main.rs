pub mod compiler;
pub mod dbal;
pub mod http;

use crate::http::compile_theme;
use std::convert::Infallible;
use std::env::args;
use tokio::net::UnixListener;
use warp::http::StatusCode;
use warp::{Filter, Rejection, Reply};

async fn handle_rejection(err: Rejection) -> Result<impl Reply, Infallible> {
    let code = if err.is_not_found() {
        StatusCode::NOT_FOUND
    } else if err.find::<warp::reject::MethodNotAllowed>().is_some() {
        StatusCode::METHOD_NOT_ALLOWED
    } else {
        StatusCode::INTERNAL_SERVER_ERROR
    };

    Ok(warp::reply::with_status(warp::reply(), code))
}

#[tokio::main]
async fn main() -> std::io::Result<()> {
    env_logger::init();

    let mut arguments = args();
    let socket = if let Some(socket) = arguments.nth(1) {
        socket
    } else {
        "goji-compiler.sock".to_string()
    };

    log::info!("Start listener on {socket}");
    let listener = UnixListener::bind(&socket).map_err(std::io::Error::other)?;

    let compile = warp::post()
        .and(warp::path::path("compile"))
        .and(warp::path::param())
        .and(warp::path::end())
        .and_then(compile_theme)
        .recover(handle_rejection);

    warp::serve(compile)
        .incoming(listener)
        .graceful(async { tokio::signal::ctrl_c().await.unwrap() })
        .run()
        .await;

    tokio::fs::remove_file(socket).await
}
