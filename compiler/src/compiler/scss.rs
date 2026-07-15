use crate::compiler::Config;
use std::collections::HashMap;

pub async fn compile_scss(
    _config: &Config,
    _variables: HashMap<String, String>,
) -> anyhow::Result<HashMap<String, String>> {
    Ok(HashMap::new())
}
