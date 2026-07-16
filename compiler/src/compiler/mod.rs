use crate::dbal::get_theme;
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::path::{Path, PathBuf};

pub mod js;
pub mod scss;

#[derive(Debug, Deserialize)]
pub struct Config {
    pub styles: Styles,
    pub assets: HashMap<String, String>,
    pub scripts: Vec<String>,
}

#[derive(Debug, Deserialize)]
pub struct Styles {
    pub files: Vec<String>,
}

#[derive(Serialize, Deserialize)]
struct CacheFile {
    styles: Vec<String>,
    scripts: Vec<String>,
}

async fn clear_cache_directory(cache_directory: &PathBuf) -> anyhow::Result<()> {
    let cache_file = Path::new(cache_directory).join("cache-status.yaml");
    let cache_file_content = tokio::fs::read_to_string(cache_file).await?;
    let cache_file = yaml_serde::from_str::<CacheFile>(&cache_file_content)?;
    for file in cache_file.styles {
        let _ = tokio::fs::remove_file(Path::new(cache_directory).join(file)).await;
    }
    for file in cache_file.scripts {
        let _ = tokio::fs::remove_file(Path::new(cache_directory).join(file)).await;
    }

    Ok(())
}

pub async fn compile_theme(theme_id: i32) -> anyhow::Result<()> {
    log::info!("Handle compile request for theme {theme_id}");
    let db_theme = get_theme(theme_id).await?;

    let full_path = Path::new(&std::env::var("THEME_PATH")?).join(&db_theme.name);
    log::info!("Read theme.yaml from {}", full_path.to_string_lossy());
    let yaml = tokio::fs::read_to_string(full_path.join("theme.yaml")).await?;

    let theme_config = yaml_serde::from_str::<Config>(&yaml)?;

    let minified_scripts = js::minify_js(&db_theme.name, &theme_config).await?;
    let compiled_styles =
        scss::compile_scss(&db_theme.name, &theme_config, db_theme.scss_variables).await?;

    let cache_path = Path::new(&std::env::var("CACHE_PATH")?).join(&db_theme.name);
    tokio::fs::create_dir_all(&cache_path).await?;

    if let Err(err) = clear_cache_directory(&cache_path).await {
        log::warn!("Failed to clear the cache for theme {theme_id}: {err}");
    }

    for (file, content) in &minified_scripts {
        let path = cache_path.clone().join(file);
        tokio::fs::write(path, content).await?;
    }
    for (file, content) in &compiled_styles {
        let path = cache_path.clone().join(file);
        tokio::fs::write(path, content).await?;
    }

    let minified_scripts = minified_scripts.keys().cloned().collect();
    let compiled_styles = compiled_styles.keys().cloned().collect();

    let cache_config = CacheFile {
        scripts: minified_scripts,
        styles: compiled_styles,
    };

    let cache_file = cache_path.join("cache-status.yaml");
    let cache_file_content = yaml_serde::to_string(&cache_config)?;

    tokio::fs::write(cache_file, cache_file_content).await?;

    Ok(())
}
