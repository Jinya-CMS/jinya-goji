use crate::compiler::Config;
use grass_compiler::sass_value::QuoteKind;
use grass_compiler::{
    sass_value::{ArgumentResult, Value}, Builtin, Options, OutputStyle, Result as SassResult,
    Visitor,
};
use std::cell::RefCell;
use std::collections::{HashMap, HashSet};
use std::ffi::OsStr;
use std::path::{Path, PathBuf};

#[derive(Clone)]
struct SassCompileContext {
    theme_name: String,
    assets: HashSet<String>,
}

thread_local! {
    static SASS_COMPILE_CONTEXT: RefCell<Option<SassCompileContext>> = const { RefCell::new(None) };
}

fn jinya_asset(mut args: ArgumentResult, _: &mut Visitor) -> SassResult<Value> {
    args.max_args(1)?;

    let context = SASS_COMPILE_CONTEXT
        .with(|context| context.borrow().clone())
        .ok_or_else(|| std::io::Error::other("Sass compile context not set"))?;

    if let Value::String(asset_name, _) = args.get_err(0, "asset")? {
        if !context.assets.contains(&asset_name) {
            return Err(std::io::Error::other(format!("Asset {asset_name} does not exist")).into());
        }

        Ok(Value::String(
            format!(
                "/jinya-content/theme/{}/asset/{asset_name}",
                context.theme_name
            ),
            QuoteKind::Quoted,
        ))
    } else {
        Err(std::io::Error::other("Asset name not given").into())
    }
}

async fn compile_file(
    theme_name: &str,
    theme_path: &PathBuf,
    config: &Config,
    file: &str,
    variables: &HashMap<String, String>,
) -> anyhow::Result<(String, String)> {
    let stylesheet_path = theme_path.join(file);
    let stylesheet = tokio::fs::read_to_string(&stylesheet_path).await?;
    let variables_style = variables
        .into_iter()
        .map(|(key, value)| format!("{key}: {value};"))
        .collect::<Vec<_>>()
        .join("\n");

    let styles_directory = stylesheet_path.parent().unwrap_or(theme_path);

    let options = Options::default()
        .add_custom_fn("jinya-asset", Builtin::new(jinya_asset))
        .load_path(styles_directory)
        .style(OutputStyle::Compressed);

    let context = SassCompileContext {
        theme_name: theme_name.to_owned(),
        assets: config.assets.keys().cloned().collect(),
    };

    SASS_COMPILE_CONTEXT.with(|compile_context| {
        *compile_context.borrow_mut() = Some(context);
    });

    let stylesheet = format!("{variables_style}\n\n{stylesheet}");
    let compiled = grass::from_string(&stylesheet, &options);
    SASS_COMPILE_CONTEXT.with(|compile_context| {
        *compile_context.borrow_mut() = None;
    });

    if let Err(err) = &compiled {
        log::error!("Failed to compile {file}: {err}");
    }
    let compiled = compiled?;

    let filename = stylesheet_path.file_stem().unwrap_or(OsStr::new(file));
    let hash = md5::compute(&compiled);

    let cached_name = format!("{}.{hash:x}.css", filename.to_string_lossy(),);

    Ok((cached_name, compiled))
}

pub async fn compile_scss(
    theme_name: &str,
    config: &Config,
    variables: HashMap<String, String>,
) -> anyhow::Result<HashMap<String, String>> {
    let base_path = std::env::var("THEME_PATH")?;
    let theme_path = Path::new(&base_path).join(theme_name);

    let mut result = HashMap::new();

    for file in &config.styles.files {
        let (name, compiled) =
            compile_file(&theme_name, &theme_path, &config, file, &variables).await?;
        result.insert(name, compiled);
    }

    Ok(result)
}
