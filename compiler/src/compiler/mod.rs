use std::collections::HashMap;
use serde::Deserialize;

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
