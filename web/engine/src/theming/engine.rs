use crate::dbal;
use crate::theming::jinya;
use crate::theming::theme::Theme;
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use wasmtime::component::{HasData, ResourceTable};
use wasmtime_wasi::{WasiCtx, WasiCtxView, WasiView};

pub struct GojiThemeStore {
    pub uri: String,
    pub path: String,
    pub theme: Theme,
    pub is_api_request: bool,
    pub wasi_ctx: WasiCtx,
    pub resource_table: ResourceTable,
}

impl WasiView for GojiThemeStore {
    fn ctx(&mut self) -> WasiCtxView<'_> {
        WasiCtxView {
            ctx: &mut self.wasi_ctx,
            table: &mut self.resource_table,
        }
    }
}

impl HasData for GojiThemeStore {
    type Data<'a> = GojiThemeStore;
}

impl jinya::goji::db::Host for GojiThemeStore {
    async fn find_gallery_files(&mut self, gallery_id: i32) -> Vec<jinya::goji::db::GalleryFile> {
        dbal::gallery::find_files(gallery_id)
            .await
            .unwrap_or_default()
            .into_iter()
            .map(Into::into)
            .collect()
    }

    async fn find_blog_post_sections(
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

    async fn find_blog_posts_by_category(
        &mut self,
        blog_category_id: i32,
        with_children: bool,
    ) -> Vec<jinya::goji::db::BlogPost> {
        dbal::blog_category::find_posts(blog_category_id, with_children)
            .await
            .unwrap_or_default()
            .into_iter()
            .map(Into::into)
            .collect()
    }

    async fn find_modern_page_sections(
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

    async fn find_menu_items_by_menu(&mut self, id: i32) -> Vec<jinya::goji::db::MenuItem> {
        dbal::menu::find_menu_items(id)
            .await
            .unwrap_or_default()
            .into_iter()
            .map(Into::into)
            .collect()
    }

    async fn find_menu_items_by_menu_item(&mut self, id: i32) -> Vec<jinya::goji::db::MenuItem> {
        dbal::menu_item::find_menu_items(id)
            .await
            .unwrap_or_default()
            .into_iter()
            .map(Into::into)
            .collect()
    }

    async fn find_menu_item_by_route(
        &mut self,
        route: String,
    ) -> Option<jinya::goji::db::MenuItem> {
        dbal::menu_item::find_menu_item_by_route(route)
            .await
            .map(Into::into)
    }

    async fn find_blog_category_by_id(&mut self, id: i32) -> Option<jinya::goji::db::BlogCategory> {
        dbal::blog_category::find_by_id(id).await.map(Into::into)
    }

    async fn find_blog_post_by_id(&mut self, id: i32) -> Option<jinya::goji::db::BlogPost> {
        dbal::blog_post::find_by_id(id).await.map(Into::into)
    }

    async fn find_classic_page_by_id(&mut self, id: i32) -> Option<jinya::goji::db::ClassicPage> {
        dbal::classic_page::find_by_id(id).await.map(Into::into)
    }

    async fn find_file_by_id(&mut self, id: i32) -> Option<jinya::goji::db::File> {
        dbal::file::find_by_id(id).await.map(Into::into)
    }

    async fn find_gallery_by_id(&mut self, id: i32) -> Option<jinya::goji::db::Gallery> {
        dbal::gallery::find_by_id(id).await.map(Into::into)
    }

    async fn find_menu_by_id(&mut self, id: i32) -> Option<jinya::goji::db::Menu> {
        dbal::menu::find_by_id(id).await.map(Into::into)
    }

    async fn find_modern_page_by_id(&mut self, id: i32) -> Option<jinya::goji::db::ModernPage> {
        dbal::modern_page::find_by_id(id).await.map(Into::into)
    }

    async fn find_user_by_id(&mut self, id: i32) -> Option<jinya::goji::db::User> {
        dbal::user::find_by_id(id).await.map(Into::into)
    }

    async fn find_all_blog_categories(&mut self) -> Vec<jinya::goji::db::BlogCategory> {
        dbal::blog_category::find_all()
            .await
            .unwrap_or_default()
            .into_iter()
            .map(Into::into)
            .collect()
    }

    async fn find_all_blog_posts(&mut self) -> Vec<jinya::goji::db::BlogPost> {
        dbal::blog_post::find_all()
            .await
            .unwrap_or_default()
            .into_iter()
            .map(Into::into)
            .collect()
    }

    async fn find_all_classic_pages(&mut self) -> Vec<jinya::goji::db::ClassicPage> {
        dbal::classic_page::find_all()
            .await
            .unwrap_or_default()
            .into_iter()
            .map(Into::into)
            .collect()
    }

    async fn find_all_files(&mut self) -> Vec<jinya::goji::db::File> {
        dbal::file::find_all()
            .await
            .unwrap_or_default()
            .into_iter()
            .map(Into::into)
            .collect()
    }

    async fn find_all_galleries(&mut self) -> Vec<jinya::goji::db::Gallery> {
        dbal::gallery::find_all()
            .await
            .unwrap_or_default()
            .into_iter()
            .map(Into::into)
            .collect()
    }

    async fn find_all_menus(&mut self) -> Vec<jinya::goji::db::Menu> {
        dbal::menu::find_all()
            .await
            .unwrap_or_default()
            .into_iter()
            .map(Into::into)
            .collect()
    }

    async fn find_all_modern_pages(&mut self) -> Vec<jinya::goji::db::ModernPage> {
        dbal::modern_page::find_all()
            .await
            .unwrap_or_default()
            .into_iter()
            .map(Into::into)
            .collect()
    }

    async fn find_all_users(&mut self) -> Vec<jinya::goji::db::User> {
        dbal::user::find_all()
            .await
            .unwrap_or_default()
            .into_iter()
            .map(Into::into)
            .collect()
    }

    async fn query_db(&mut self, query: String) -> Vec<HashMap<String, String>> {
        dbal::select(query).await
    }
}

#[derive(Debug, Deserialize)]
pub struct Styles {
    pub files: Vec<String>,
}

#[derive(Default, Serialize, Deserialize)]
struct CacheFile {
    styles: Vec<String>,
    scripts: Vec<String>,
}

impl jinya::goji::theme::Host for GojiThemeStore {
    async fn get_theme_config_value(
        &mut self,
        group: String,
        key: String,
    ) -> Option<jinya::goji::theme::BoolOrString> {
        self.theme.get_config(group, key).map(Into::into)
    }

    async fn get_active_menu_item(&mut self) -> Option<jinya::goji::theme::MenuItem> {
        dbal::menu_item::get_item_by_path(self.path.clone())
            .await
            .map(Into::into)
    }

    async fn is_menu_item_active(&mut self, item: jinya::goji::theme::MenuItem) -> bool {
        let route = item.route.clone().unwrap_or_default();
        let route = route.strip_prefix("/").unwrap_or(route.as_str());

        let path = self.path.strip_prefix("/").unwrap_or(self.path.as_str());

        route == path
    }

    async fn is_parent_of_active_menu_item(&mut self, item: jinya::goji::theme::MenuItem) -> bool {
        if self.path.is_empty() || self.path == "/" {
            return false;
        }
        dbal::menu_item::is_parent_of_active_menu_item(item.id, &self.path).await
    }

    async fn find_modern_page_by_theme(
        &mut self,
        name: String,
    ) -> Option<jinya::goji::theme::ModernPage> {
        dbal::modern_page::find_by_theme(name, self.theme.db_theme.id)
            .await
            .map(Into::into)
    }

    async fn find_classic_page_by_theme(
        &mut self,
        name: String,
    ) -> Option<jinya::goji::theme::ClassicPage> {
        dbal::classic_page::find_by_theme(name, self.theme.db_theme.id)
            .await
            .map(Into::into)
    }

    async fn find_file_by_theme(&mut self, name: String) -> Option<jinya::goji::theme::File> {
        dbal::file::find_by_theme(name, self.theme.db_theme.id)
            .await
            .map(Into::into)
    }

    async fn find_menu_by_theme(&mut self, name: String) -> Option<jinya::goji::theme::Menu> {
        dbal::menu::find_by_theme(name, self.theme.db_theme.id)
            .await
            .map(Into::into)
    }

    async fn find_gallery_by_theme(&mut self, name: String) -> Option<jinya::goji::theme::Gallery> {
        dbal::gallery::find_by_theme(name, self.theme.db_theme.id)
            .await
            .map(Into::into)
    }

    async fn find_blog_category_by_theme(
        &mut self,
        name: String,
    ) -> Option<jinya::goji::theme::BlogCategory> {
        dbal::blog_category::find_by_theme(name, self.theme.db_theme.id)
            .await
            .map(Into::into)
    }

    async fn scripts(&mut self) -> String {
        let cache_status = self.parse_cache_file().await;
        cache_status
            .scripts
            .into_iter()
            .map(|file| {
                format!(
                    r#"<script src="/jinya-content/themes/cache/{}/{file}"></script>"#,
                    self.theme.db_theme.name
                )
            })
            .collect::<Vec<_>>()
            .join("\n")
            .to_string()
    }

    async fn styles(&mut self) -> String {
        let cache_status = self.parse_cache_file().await;
        cache_status
            .styles
            .into_iter()
            .map(|file| format!(
                r#"<link rel="stylesheet" type="text/css" href="/jinya-content/themes/cache/{}/{file}">"#,
                self.theme.db_theme.name
            ))
            .collect::<Vec<_>>()
            .join("\n")
            .to_string()
    }

    async fn asset_path(&mut self, name: String) -> String {
        format!(
            "/jinya-content/themes/{}/asset/{name}",
            self.theme.db_theme.name
        )
    }
}

impl GojiThemeStore {
    async fn parse_cache_file(&mut self) -> CacheFile {
        yaml_serde::from_str(
            &tokio::fs::read_to_string(
                std::path::Path::new(&std::env::var("CACHE_PATH").unwrap_or_default())
                    .join(&self.theme.db_theme.name)
                    .join("cache-status.yaml"),
            )
            .await
            .unwrap_or_default(),
        )
        .unwrap_or_default()
    }
}

impl jinya::goji::types::Host for GojiThemeStore {}

impl jinya::goji::http::Host for GojiThemeStore {
    async fn get(&mut self, url: String) -> Result<String, String> {
        let res = reqwest::get(url).await;
        match res {
            Ok(res) => res.text().await.map_err(|err| err.to_string()),
            Err(err) => Err(err.to_string()),
        }
    }
}

impl jinya::goji::request::Host for GojiThemeStore {
    async fn get_path(&mut self) -> String {
        self.path.clone()
    }

    async fn get_url(&mut self) -> String {
        self.uri.clone()
    }

    async fn is_api_request(&mut self) -> bool {
        self.is_api_request
    }
}

impl jinya::goji::blog::Host for GojiThemeStore {
    async fn get_category_path_for_post(
        &mut self,
        post_id: i32,
    ) -> Vec<jinya::goji::blog::BlogCategoryPathEntry> {
        dbal::blog_post::get_category_path(post_id)
            .await
            .into_iter()
            .map(Into::into)
            .collect()
    }

    async fn get_category_path_for_category(
        &mut self,
        category_id: i32,
    ) -> Vec<jinya::goji::blog::BlogCategoryPathEntry> {
        dbal::blog_category::get_category_path(category_id)
            .await
            .into_iter()
            .map(Into::into)
            .collect()
    }
}

impl jinya::goji::image::Host for GojiThemeStore {
    async fn image_size(&mut self, file_id: i32) -> (i32, i32) {
        if let Some(file) = dbal::file::find_by_id(file_id).await {
            if let Ok(size) = imagesize::size(
                std::path::Path::new(&std::env::var("STORAGE_PATH").unwrap_or("data".to_string()))
                    .join(file.path),
            ) {
                (size.width as i32, size.height as i32)
            } else {
                (-1, -1)
            }
        } else {
            (-1, -1)
        }
    }
}
