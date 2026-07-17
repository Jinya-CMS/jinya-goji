use crate::models;
use anyhow::anyhow;
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::path::Path;

fn is_false() -> bool {
    false
}

#[derive(Serialize, Deserialize, Debug)]
pub struct Theme {
    #[serde(skip)]
    pub db_theme: models::theme::Theme,
    #[serde(rename = "configuration")]
    pub default_configuration: HashMap<String, HashMap<String, String>>,
    #[serde(rename = "errorBehavior")]
    pub error_behavior: Option<String>,
    pub assets: HashMap<String, String>,
    #[serde(rename = "hasApi", default = "is_false")]
    pub has_api: bool,
}

static THEME: tokio::sync::OnceCell<Theme> = tokio::sync::OnceCell::const_new();
static ACTIVE_THEME_ID: tokio::sync::OnceCell<i32> = tokio::sync::OnceCell::const_new();

impl Theme {
    pub async fn new<'a>(db_theme: models::theme::Theme) -> anyhow::Result<&'a Self> {
        if let Some(active_theme_id) = ACTIVE_THEME_ID.get()
            && *active_theme_id == db_theme.id
        {
            THEME
                .get_or_try_init(move || Theme::from_yaml(db_theme))
                .await
        } else {
            let theme_id = db_theme.id;

            let theme = Theme::from_yaml(db_theme).await?;

            THEME.set(theme)?;
            ACTIVE_THEME_ID.set(theme_id)?;

            Ok(THEME.get().unwrap())
        }
    }

    async fn from_yaml(db_theme: models::theme::Theme) -> anyhow::Result<Theme> {
        log::debug!("Parse the theme.yaml file");
        let theme_yaml_file = Path::new(&db_theme.path).join("theme.yaml");

        if !theme_yaml_file.exists() {
            log::error!("Could not find theme.yaml: {}", theme_yaml_file.display());

            return Err(anyhow!(
                "theme.yaml is not found, expecting it here: {}",
                theme_yaml_file.display()
            ));
        }

        log::debug!("Load theme.yaml");

        let data = tokio::fs::read_to_string(&theme_yaml_file).await?;
        let mut theme = serde_yaml::from_str::<Theme>(&data)?;
        theme.db_theme = db_theme;

        Ok(theme)
    }
}
