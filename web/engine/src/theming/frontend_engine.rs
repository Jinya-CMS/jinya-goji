use crate::theming::engine::GojiEngine;
use crate::theming::jinya;
use crate::{
    dbal, goji_blog, goji_db, goji_engine, goji_http, goji_image, goji_request, goji_theme,
};

goji_engine!(
    GojiFrontendEngineStore,
    pub is_api_request: bool,
);

goji_db!(GojiFrontendEngineStore);

goji_theme!(GojiFrontendEngineStore);

goji_request!(GojiFrontendEngineStore);

goji_http!(GojiFrontendEngineStore);

goji_image!(GojiFrontendEngineStore);

goji_blog!(GojiFrontendEngineStore);

impl GojiEngine for GojiFrontendEngineStore {
    async fn engine_find_blog_post_sections(
        &mut self,
        blog_post_id: i32,
    ) -> Vec<jinya::goji::db::BlogPostSection> {
        dbal::blog_post::find_sections(blog_post_id)
            .await
            .unwrap_or_default()
            .into_iter()
            .map(Into::into)
            .collect()
    }

    async fn engine_find_modern_page_sections(
        &mut self,
        modern_page_id: i32,
    ) -> Vec<jinya::goji::db::ModernPageSection> {
        dbal::modern_page::find_sections(modern_page_id)
            .await
            .unwrap_or_default()
            .into_iter()
            .map(Into::into)
            .collect()
    }

    async fn engine_find_classic_page_by_id(
        &mut self,
        id: i32,
    ) -> Option<jinya::goji::db::ClassicPage> {
        dbal::classic_page::find_by_id(id).await.map(Into::into)
    }

    async fn engine_find_all_classic_pages(&mut self) -> Vec<jinya::goji::db::ClassicPage> {
        dbal::classic_page::find_all()
            .await
            .unwrap_or_default()
            .into_iter()
            .map(Into::into)
            .collect()
    }

    async fn engine_find_classic_page_by_theme(
        &mut self,
        name: String,
    ) -> Option<jinya::goji::theme::ClassicPage> {
        dbal::classic_page::find_by_theme(name, self.theme.db_theme.id)
            .await
            .map(Into::into)
    }

    fn engine_is_api_request(&mut self) -> bool {
        self.is_api_request
    }
}
