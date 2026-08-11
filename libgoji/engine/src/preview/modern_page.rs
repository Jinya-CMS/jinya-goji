use dbal::models::ModernPageSection;
use crate::preview::get_goji_theme;
use crate::theming::preview_engine::PreviewEntity;
use crate::theming::theme::Theme;
use anyhow::anyhow;

pub async fn modern_page_preview(
    theme_id: i32,
    modern_page_id: i32,
    modern_page_sections_data: String,
) -> anyhow::Result<String> {
    let theme = dbal::theme::find_theme_by_id(theme_id).await?;

    let theme = Theme::new(theme).await?;

    let route = if let Some(item) =
        dbal::menu_item::find_first_by_modern_page(modern_page_id, theme_id).await
    {
        item.route.unwrap_or_default()
    } else {
        "/".to_string()
    };

    let modern_page_sections =
        serde_json::from_str::<Vec<ModernPageSection>>(&modern_page_sections_data)?;

    let (goji_theme, store) = get_goji_theme(
        &theme,
        PreviewEntity::ModernPage((modern_page_id, modern_page_sections.clone())),
        route,
    )
    .await?;

    let modern_page = dbal::modern_page::find_by_id(modern_page_id)
        .await
        .ok_or_else(|| anyhow!("Modern page not found"))?;

    goji_theme
        .func_handle_modern_page()
        .call_async(store, (modern_page.into(),))
        .await
        .map(|(response, ..)| response.content)
        .map_err(|err| err.into())
}
