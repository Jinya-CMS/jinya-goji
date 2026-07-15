use crate::compiler::js;
use crate::compiler::scss;
use crate::compiler::Config;
use crate::dbal::get_theme;
use serde::Serialize;
use std::fmt::Debug;
use warp::reject::Reject;

#[derive(Debug)]
pub struct Failure {}

impl Reject for Failure {}

#[derive(Serialize)]
struct CacheFile {
    styles: Vec<String>,
    scripts: Vec<String>,
}

pub async fn compile_theme(theme_id: i32) -> Result<impl warp::Reply, warp::Rejection> {
    log::info!("Handle compile request for theme {theme_id}");
    let db_theme = get_theme(theme_id)
        .await
        .map_err(|_| warp::reject::not_found())?;

    let full_path =
        std::path::Path::new(&std::env::var("THEME_PATH").unwrap()).join(&db_theme.name);
    log::info!("Read theme.yaml from {}", full_path.to_string_lossy());
    let yaml = tokio::fs::read_to_string(full_path.join("theme.yaml"))
        .await
        .map_err(|_| warp::reject::not_found())?;

    let theme_config =
        yaml_serde::from_str::<Config>(&yaml).map_err(|_| warp::reject::not_found())?;

    let minified_scripts = js::minify_js(&db_theme.name, &theme_config)
        .await
        .map_err(|_| Failure {})?;
    let compiled_styles = scss::compile_scss(&theme_config, db_theme.scss_variables)
        .await
        .map_err(|_| Failure {})?;

    let cache_path =
        std::path::Path::new(&std::env::var("CACHE_PATH").unwrap()).join(&db_theme.name);
    tokio::fs::create_dir_all(&cache_path)
        .await
        .map_err(|_| Failure {})?;

    for (file, content) in &minified_scripts {
        let path = cache_path.clone().join(file);
        tokio::fs::write(path, content)
            .await
            .map_err(|_| Failure {})?;
    }
    for (file, content) in &compiled_styles {
        let path = cache_path.clone().join(file);
        tokio::fs::write(path, content)
            .await
            .map_err(|_| Failure {})?;
    }

    let minified_scripts = minified_scripts.keys().cloned().collect();
    let compiled_styles = compiled_styles.keys().cloned().collect();

    let cache_config = CacheFile {
        scripts: minified_scripts,
        styles: compiled_styles,
    };

    let cache_file = cache_path.join("cache-status.yaml");
    let cache_file_content = yaml_serde::to_string(&cache_config).map_err(|_| Failure {})?;

    tokio::fs::write(cache_file, cache_file_content)
        .await
        .map_err(|_| Failure {})?;

    Ok(warp::reply::with_status(
        warp::reply::reply(),
        warp::http::StatusCode::NO_CONTENT,
    ))
}
