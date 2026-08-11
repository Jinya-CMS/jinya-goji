use crate::models::{BlogPostSection, ClassicPage, ModernPageSection};
use crate::theming::engine::GojiEngine;
use crate::theming::jinya;
use crate::{
    dbal, goji_blog, goji_db, goji_engine, goji_http, goji_image, goji_request, goji_theme,
};

#[derive(Clone)]
pub enum PreviewEntity {
    ClassicPage(ClassicPage),
    ModernPage((i32, Vec<ModernPageSection>)),
    BlogPost((i32, Vec<BlogPostSection>)),
}

goji_engine!(
    GojiPreviewEngineStore,
    pub preview_entity: PreviewEntity,
);

goji_db!(GojiPreviewEngineStore);

goji_theme!(GojiPreviewEngineStore);

goji_request!(GojiPreviewEngineStore);

goji_http!(GojiPreviewEngineStore);

goji_image!(GojiPreviewEngineStore);

goji_blog!(GojiPreviewEngineStore);

impl GojiEngine for GojiPreviewEngineStore {
    async fn engine_find_blog_post_sections(
        &mut self,
        blog_post_id: i32,
    ) -> Vec<jinya::goji::db::BlogPostSection> {
        if let PreviewEntity::BlogPost((id, sections)) = self.preview_entity.clone()
            && blog_post_id == id
        {
            sections
        } else {
            dbal::blog_post::find_sections(blog_post_id)
                .await
                .unwrap_or_default()
        }
        .into_iter()
        .map(Into::into)
        .collect()
    }

    async fn engine_find_modern_page_sections(
        &mut self,
        modern_page_id: i32,
    ) -> Vec<jinya::goji::db::ModernPageSection> {
        if let PreviewEntity::ModernPage((id, sections)) = self.preview_entity.clone()
            && modern_page_id == id
        {
            sections
        } else {
            dbal::modern_page::find_sections(modern_page_id)
                .await
                .unwrap_or_default()
        }
        .into_iter()
        .map(Into::into)
        .collect()
    }

    async fn engine_find_classic_page_by_id(
        &mut self,
        id: i32,
    ) -> Option<jinya::goji::db::ClassicPage> {
        if let PreviewEntity::ClassicPage(classic_page) = self.preview_entity.clone() {
            Some(classic_page.clone().into())
        } else {
            dbal::classic_page::find_by_id(id).await.map(Into::into)
        }
    }

    async fn engine_find_all_classic_pages(&mut self) -> Vec<jinya::goji::db::ClassicPage> {
        let mut pages = dbal::classic_page::find_all().await.unwrap_or_default();

        if let PreviewEntity::ClassicPage(classic_page) = self.preview_entity.clone() {
            if let Some(idx) = pages.iter().position(|page| page.id == classic_page.id) {
                pages[idx] = classic_page.clone();
            }
        }

        pages.into_iter().map(Into::into).collect()
    }

    async fn engine_find_classic_page_by_theme(
        &mut self,
        name: String,
    ) -> Option<jinya::goji::theme::ClassicPage> {
        let page = dbal::classic_page::find_by_theme(name, self.theme.db_theme.id).await;
        if let Some(page) = page {
            if let PreviewEntity::ClassicPage(classic_page) = self.preview_entity.clone()
                && page.id == classic_page.id
            {
                Some(classic_page.into())
            } else {
                Some(page.into())
            }
        } else {
            None
        }
    }

    fn engine_is_api_request(&mut self) -> bool {
        false
    }
}
