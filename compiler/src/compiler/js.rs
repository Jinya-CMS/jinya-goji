use crate::compiler::Config;
use std::collections::HashMap;
use std::ffi::OsStr;
use std::path::Path;
use std::sync::Arc;
use swc::config::JsMinifyOptions;
use swc::{BoolOrDataConfig, Compiler, JsMinifyExtras};
use swc_common::errors::{ColorConfig, Handler};
use swc_common::{FileName, Globals, SourceMap, GLOBALS};

async fn minify_script(source: &str) -> anyhow::Result<String> {
    let cm = Arc::new(SourceMap::default());
    let compiler = Compiler::new(cm.clone());

    let fm = cm.new_source_file(FileName::Anon.into(), source.to_string());

    let options = JsMinifyOptions {
        compress: BoolOrDataConfig::from_bool(true),
        mangle: BoolOrDataConfig::from_bool(true),
        ..Default::default()
    };

    let handler = Handler::with_tty_emitter(ColorConfig::Auto, true, false, Some(cm.clone()));

    let globals = Globals::new();
    GLOBALS.set(&globals, || {
        compiler
            .minify(fm, &handler, &options, JsMinifyExtras::default())
            .map(|output| output.code)
    })
}

pub async fn minify_js(
    theme_name: &str,
    config: &Config,
) -> anyhow::Result<HashMap<String, String>> {
    let mut result = HashMap::new();
    for file in &config.scripts {
        let full_path = Path::new(&std::env::var("THEME_PATH").unwrap())
            .join(theme_name)
            .join(file);
        let script = tokio::fs::read_to_string(full_path).await?;
        let path = Path::new(file);
        let extension = path.extension().unwrap_or(OsStr::new("js"));
        let filename = path.file_stem().unwrap_or(OsStr::new(file));
        let hash = md5::compute(&script);

        let cached_name = format!(
            "{}.{hash:x}.{}",
            filename.to_string_lossy(),
            extension.to_string_lossy()
        );
        result.insert(cached_name, minify_script(&script).await?);
    }

    Ok(result)
}
