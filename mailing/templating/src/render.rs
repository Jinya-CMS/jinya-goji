use mrml::prelude::render::RenderOptions;

pub enum MrmlError {
    ParseError,
    RenderError,
}

pub fn render_template(tmpl: &str) -> Result<String, MrmlError> {
    mrml::parse(tmpl)
        .map_err(|err| {
            log::error!("Failed to parse template {err}");
            MrmlError::ParseError
        })?
        .element
        .render(&RenderOptions::default())
        .map_err(|err| {
            log::error!("Failed to render template {err}");
            MrmlError::RenderError
        })
}
