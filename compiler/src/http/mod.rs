use std::fmt::Debug;
use warp::reject::Reject;

#[derive(Debug)]
pub struct Failure {}

impl Reject for Failure {}

pub async fn compile_theme(theme_id: i32) -> Result<impl warp::Reply, warp::Rejection> {
    crate::compiler::compile_theme(theme_id)
        .await
        .map_err(|_| Failure {})?;

    Ok(warp::reply::with_status(
        warp::reply::reply(),
        warp::http::StatusCode::NO_CONTENT,
    ))
}
