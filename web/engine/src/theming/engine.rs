use serde::{Deserialize, Serialize};
use crate::theming::jinya;
use crate::theming::theme::Theme;

#[derive(Default, Serialize, Deserialize)]
pub struct CacheFile {
    pub styles: Vec<String>,
    pub scripts: Vec<String>,
}

pub async fn parse_cache_file(theme: &Theme) -> CacheFile {
    yaml_serde::from_str(
        &tokio::fs::read_to_string(
            std::path::Path::new(&std::env::var("CACHE_PATH").unwrap_or_default())
                .join(&theme.db_theme.name)
                .join("cache-status.yaml"),
        )
            .await
            .unwrap_or_default(),
    )
        .unwrap_or_default()
}

#[macro_export]
macro_rules! goji_engine {
    ($engine:ident $(, $($additional_props:tt)*)?) => {
        pub struct $engine {
            pub uri: String,
            pub path: String,
            pub theme: crate::theming::theme::Theme,
            pub wasi_ctx: wasmtime_wasi::WasiCtx,
            pub resource_table: wasmtime::component::ResourceTable,
            $($($additional_props)*)?
        }

        impl crate::theming::jinya::goji::types::Host for $engine {}

        impl wasmtime_wasi::WasiView for $engine {
            fn ctx(&mut self) -> wasmtime_wasi::WasiCtxView<'_> {
                wasmtime_wasi::WasiCtxView {
                    ctx: &mut self.wasi_ctx,
                    table: &mut self.resource_table,
                }
            }
        }

        impl wasmtime::component::HasData for $engine {
            type Data<'a> = $engine;
        }
    };
}

#[macro_export]
macro_rules! goji_http {
    ($engine:ident) => {
        impl crate::theming::jinya::goji::http::Host for $engine {
            async fn get(&mut self, url: String) -> Result<String, String> {
                let res = reqwest::get(url).await;
                match res {
                    Ok(res) => res.text().await.map_err(|err| err.to_string()),
                    Err(err) => Err(err.to_string()),
                }
            }
        }
    };
}

#[macro_export]
macro_rules! goji_image {
    ($engine:ident) => {
        impl crate::theming::jinya::goji::image::Host for $engine {
            async fn image_size(&mut self, file_id: i32) -> (i32, i32) {
                if let Some(file) = crate::dbal::file::find_by_id(file_id).await {
                    if let Ok(size) = imagesize::size(
                        std::path::Path::new(
                            &std::env::var("STORAGE_PATH").unwrap_or("data".to_string()),
                        )
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
    };
}

#[macro_export]
macro_rules! goji_blog {
    ($engine:ident) => {
        impl crate::theming::jinya::goji::blog::Host for $engine {
            async fn get_category_path_for_post(
                &mut self,
                post_id: i32,
            ) -> Vec<crate::theming::jinya::goji::blog::BlogCategoryPathEntry> {
                crate::dbal::blog_post::get_category_path(post_id)
                    .await
                    .into_iter()
                    .map(Into::into)
                    .collect()
            }

            async fn get_category_path_for_category(
                &mut self,
                category_id: i32,
            ) -> Vec<crate::theming::jinya::goji::blog::BlogCategoryPathEntry> {
                crate::dbal::blog_category::get_category_path(category_id)
                    .await
                    .into_iter()
                    .map(Into::into)
                    .collect()
            }
        }
    };
}

#[macro_export]
macro_rules! goji_theme {
    ($engine:ident) => {
        impl crate::theming::jinya::goji::theme::Host for $engine {
            async fn get_theme_config_value(
                &mut self,
                group: String,
                key: String,
            ) -> Option<crate::theming::jinya::goji::theme::BoolOrString> {
                self.theme.get_config(group, key).map(Into::into)
            }

            async fn get_active_menu_item(&mut self) -> Option<crate::theming::jinya::goji::theme::MenuItem> {
                crate::dbal::menu_item::get_item_by_path(self.path.clone())
                    .await
                    .map(Into::into)
            }

            async fn is_menu_item_active(&mut self, item: crate::theming::jinya::goji::theme::MenuItem) -> bool {
                let route = item.route.clone().unwrap_or_default();
                let route = route.strip_prefix("/").unwrap_or(route.as_str());

                let path = self.path.strip_prefix("/").unwrap_or(self.path.as_str());

                route == path
            }

            async fn is_parent_of_active_menu_item(&mut self, item: crate::theming::jinya::goji::theme::MenuItem) -> bool {
                if self.path.is_empty() || self.path == "/" {
                    return false;
                }

                crate::dbal::menu_item::is_parent_of_active_menu_item(item.id, &self.path).await
            }

            async fn find_modern_page_by_theme(
                &mut self,
                name: String,
            ) -> Option<crate::theming::jinya::goji::theme::ModernPage> {
                crate::dbal::modern_page::find_by_theme(name, self.theme.db_theme.id)
                    .await
                    .map(Into::into)
            }

            async fn find_classic_page_by_theme(
                &mut self,
                name: String,
            ) -> Option<crate::theming::jinya::goji::theme::ClassicPage> {
                self.engine_find_classic_page_by_theme(name).await
            }

            async fn find_file_by_theme(&mut self, name: String) -> Option<crate::theming::jinya::goji::theme::File> {
                crate::dbal::file::find_by_theme(name, self.theme.db_theme.id)
                    .await
                    .map(Into::into)
            }

            async fn find_menu_by_theme(&mut self, name: String) -> Option<crate::theming::jinya::goji::theme::Menu> {
                crate::dbal::menu::find_by_theme(name, self.theme.db_theme.id)
                    .await
                    .map(Into::into)
            }

            async fn find_gallery_by_theme(&mut self, name: String) -> Option<crate::theming::jinya::goji::theme::Gallery> {
                crate::dbal::gallery::find_by_theme(name, self.theme.db_theme.id)
                    .await
                    .map(Into::into)
            }

            async fn find_blog_category_by_theme(
                &mut self,
                name: String,
            ) -> Option<crate::theming::jinya::goji::theme::BlogCategory> {
                crate::dbal::blog_category::find_by_theme(name, self.theme.db_theme.id)
                    .await
                    .map(Into::into)
            }

            async fn scripts(&mut self) -> String {
                let cache_status = crate::theming::engine::parse_cache_file(&self.theme).await;
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
                let cache_status = crate::theming::engine::parse_cache_file(&self.theme).await;
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
                format!("/jinya-content/themes/{}/asset/{name}", self.theme.db_theme.name)
            }
        }
    };
}

#[macro_export]
macro_rules! goji_db {
    ($engine:ident) => {
        impl crate::theming::jinya::goji::db::Host for $engine {
            async fn find_gallery_files(&mut self, gallery_id: i32) -> Vec<crate::theming::jinya::goji::db::GalleryFile> {
                crate::dbal::gallery::find_files(gallery_id)
                    .await
                    .unwrap_or_default()
                    .into_iter()
                    .map(Into::into)
                    .collect()
            }

            async fn find_blog_post_sections(
                &mut self,
                blog_post_id: i32,
            ) -> Vec<crate::theming::jinya::goji::db::BlogPostSection> {
                self.engine_find_blog_post_sections(blog_post_id).await
            }

            async fn find_blog_posts_by_category(
                &mut self,
                blog_category_id: i32,
                with_children: bool,
            ) -> Vec<crate::theming::jinya::goji::db::BlogPost> {
                crate::dbal::blog_category::find_posts(blog_category_id, with_children)
                    .await
                    .unwrap_or_default()
                    .into_iter()
                    .map(Into::into)
                    .collect()
            }

            async fn find_modern_page_sections(
                &mut self,
                modern_page_id: i32,
            ) -> Vec<crate::theming::jinya::goji::db::ModernPageSection> {
                self.engine_find_modern_page_sections(modern_page_id).await
            }

            async fn find_menu_items_by_menu(&mut self, id: i32) -> Vec<crate::theming::jinya::goji::db::MenuItem> {
                crate::dbal::menu::find_menu_items(id)
                    .await
                    .unwrap_or_default()
                    .into_iter()
                    .map(Into::into)
                    .collect()
            }

            async fn find_menu_items_by_menu_item(&mut self, id: i32) -> Vec<crate::theming::jinya::goji::db::MenuItem> {
                crate::dbal::menu_item::find_menu_items(id)
                    .await
                    .unwrap_or_default()
                    .into_iter()
                    .map(Into::into)
                    .collect()
            }

            async fn find_menu_item_by_route(
                &mut self,
                route: String,
            ) -> Option<crate::theming::jinya::goji::db::MenuItem> {
                crate::dbal::menu_item::find_menu_item_by_route(route)
                    .await
                    .map(Into::into)
            }

            async fn find_blog_category_by_id(&mut self, id: i32) -> Option<crate::theming::jinya::goji::db::BlogCategory> {
                crate::dbal::blog_category::find_by_id(id).await.map(Into::into)
            }

            async fn find_blog_post_by_id(&mut self, id: i32) -> Option<crate::theming::jinya::goji::db::BlogPost> {
                crate::dbal::blog_post::find_by_id(id).await.map(Into::into)
            }

            async fn find_classic_page_by_id(&mut self, id: i32) -> Option<crate::theming::jinya::goji::db::ClassicPage> {
                self.engine_find_classic_page_by_id(id).await
            }

            async fn find_file_by_id(&mut self, id: i32) -> Option<crate::theming::jinya::goji::db::File> {
                crate::dbal::file::find_by_id(id).await.map(Into::into)
            }

            async fn find_gallery_by_id(&mut self, id: i32) -> Option<crate::theming::jinya::goji::db::Gallery> {
                crate::dbal::gallery::find_by_id(id).await.map(Into::into)
            }

            async fn find_menu_by_id(&mut self, id: i32) -> Option<crate::theming::jinya::goji::db::Menu> {
                crate::dbal::menu::find_by_id(id).await.map(Into::into)
            }

            async fn find_modern_page_by_id(&mut self, id: i32) -> Option<crate::theming::jinya::goji::db::ModernPage> {
                crate::dbal::modern_page::find_by_id(id).await.map(Into::into)
            }

            async fn find_user_by_id(&mut self, id: i32) -> Option<crate::theming::jinya::goji::db::User> {
                crate::dbal::user::find_by_id(id).await.map(Into::into)
            }

            async fn find_all_blog_categories(&mut self) -> Vec<crate::theming::jinya::goji::db::BlogCategory> {
                crate::dbal::blog_category::find_all()
                    .await
                    .unwrap_or_default()
                    .into_iter()
                    .map(Into::into)
                    .collect()
            }

            async fn find_all_blog_posts(&mut self) -> Vec<crate::theming::jinya::goji::db::BlogPost> {
                crate::dbal::blog_post::find_all()
                    .await
                    .unwrap_or_default()
                    .into_iter()
                    .map(Into::into)
                    .collect()
            }

            async fn find_all_classic_pages(&mut self) -> Vec<crate::theming::jinya::goji::db::ClassicPage> {
                self.engine_find_all_classic_pages().await
            }

            async fn find_all_files(&mut self) -> Vec<crate::theming::jinya::goji::db::File> {
                crate::dbal::file::find_all()
                    .await
                    .unwrap_or_default()
                    .into_iter()
                    .map(Into::into)
                    .collect()
            }

            async fn find_all_galleries(&mut self) -> Vec<crate::theming::jinya::goji::db::Gallery> {
                crate::dbal::gallery::find_all()
                    .await
                    .unwrap_or_default()
                    .into_iter()
                    .map(Into::into)
                    .collect()
            }

            async fn find_all_menus(&mut self) -> Vec<crate::theming::jinya::goji::db::Menu> {
                crate::dbal::menu::find_all()
                    .await
                    .unwrap_or_default()
                    .into_iter()
                    .map(Into::into)
                    .collect()
            }

            async fn find_all_modern_pages(&mut self) -> Vec<crate::theming::jinya::goji::db::ModernPage> {
                crate::dbal::modern_page::find_all()
                    .await
                    .unwrap_or_default()
                    .into_iter()
                    .map(Into::into)
                    .collect()
            }

            async fn find_all_users(&mut self) -> Vec<crate::theming::jinya::goji::db::User> {
                crate::dbal::user::find_all()
                    .await
                    .unwrap_or_default()
                    .into_iter()
                    .map(Into::into)
                    .collect()
            }

            async fn query_db(&mut self, query: String) -> Vec<std::collections::HashMap<String, String>> {
                crate::dbal::select(query).await
            }
        }
    };
}

#[macro_export]
macro_rules! goji_request {
    ($engine:ident) => {
        impl crate::theming::jinya::goji::request::Host for $engine {
            async fn get_path(&mut self) -> String {
                self.path.clone()
            }

            async fn get_url(&mut self) -> String {
                self.uri.clone()
            }

            async fn is_api_request(&mut self) -> bool {
                self.engine_is_api_request()
            }
        }
    };
}

pub trait GojiEngine {
    fn engine_find_classic_page_by_theme(
        &mut self,
        name: String,
    ) -> impl Future<Output = Option<jinya::goji::theme::ClassicPage>>;

    fn engine_find_blog_post_sections(
        &mut self,
        blog_post_id: i32,
    ) -> impl Future<Output = Vec<crate::theming::jinya::goji::db::BlogPostSection>>;

    fn engine_find_modern_page_sections(
        &mut self,
        modern_page_id: i32,
    ) -> impl Future<Output = Vec<crate::theming::jinya::goji::db::ModernPageSection>>;

    fn engine_find_classic_page_by_id(
        &mut self,
        id: i32,
    ) -> impl Future<Output = Option<crate::theming::jinya::goji::db::ClassicPage>>;

    fn engine_find_all_classic_pages(&mut self) -> impl Future<Output = Vec<crate::theming::jinya::goji::db::ClassicPage>>;

    fn engine_is_api_request(&mut self) -> bool {
        false
    }
}
