use crate::preview::get_goji_theme;
use crate::theming::preview_engine::PreviewEntity;
use crate::theming::theme::Theme;
use anyhow::anyhow;
use dbal::models::ModernPageSection;
use wasi_bridge::jinya;

pub async fn blog_post_preview(
    theme_id: i32,
    blog_post_id: i32,
    blog_post_sections_data: String,
) -> anyhow::Result<String> {
    let theme = dbal::theme::find_theme_by_id(theme_id).await?;

    let theme = Theme::new(theme).await?;

    let blog_post_sections =
        serde_json::from_str::<Vec<ModernPageSection>>(&blog_post_sections_data)?;

    let blog_post = jinya::goji::types::BlogPost::from(
        dbal::blog_post::find_by_id(blog_post_id)
            .await
            .ok_or_else(|| anyhow!("Blog post not found"))?,
    );

    let (goji_theme, store) = get_goji_theme(
        &theme,
        PreviewEntity::ModernPage((blog_post_id, blog_post_sections.clone())),
        blog_post.link.clone(),
    )
    .await?;

    goji_theme
        .func_handle_blog_post()
        .call_async(store, (blog_post,))
        .await
        .map(|(response, ..)| response.content)
        .map_err(|err| err.into())
}
