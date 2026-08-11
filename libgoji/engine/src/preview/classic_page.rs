use dbal::models::ClassicPage;
use crate::preview::get_goji_theme;
use crate::theming::preview_engine::PreviewEntity;
use crate::theming::theme::Theme;

pub async fn classic_page_preview(
    theme_id: i32,
    classic_page_id: i32,
    classic_page_data: String,
) -> anyhow::Result<String> {
    let theme = dbal::theme::find_theme_by_id(theme_id).await?;

    let theme = Theme::new(theme).await?;

    let route = if let Some(item) =
        dbal::menu_item::find_first_by_classic_page(classic_page_id, theme_id).await
    {
        item.route.unwrap_or_default()
    } else {
        "/".to_string()
    };

    let classic_page = serde_json::from_str::<ClassicPage>(&classic_page_data)?;

    let (goji_theme, store) = get_goji_theme(
        &theme,
        PreviewEntity::ClassicPage(classic_page.clone()),
        route,
    )
    .await?;

    goji_theme
        .func_handle_classic_page()
        .call_async(store, (classic_page.into(),))
        .await
        .map(|(response, ..)| response.content)
        .map_err(|err| err.into())
}
