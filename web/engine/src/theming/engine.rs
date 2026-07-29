use crate::models::blog_category::BlogCategoryPath;
use crate::models::gallery_file_position::GalleryFileWithFile;
use crate::models::theme::BoolOrString;
use crate::theming::engine::jinya::goji::{blog, db, http, image, request, theme, types};
use crate::theming::theme::Theme;
use crate::{dbal, models};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::path::Path;
use wasmtime::component::{HasData, ResourceTable, bindgen};
use wasmtime_wasi::{WasiCtx, WasiCtxView, WasiView};

bindgen!({
    path: "wit/goji.wit",
    imports: {
        default: async
    },
    exports: {
        default: async
    }
});

impl From<models::User> for User {
    fn from(user: models::User) -> Self {
        Self {
            email: user.email,
            artist_name: user.artist_name,
            profile_picture: user.profile_picture,
            about_me: user.about_me,
        }
    }
}

impl From<models::File> for File {
    fn from(file: models::File) -> Self {
        Self {
            id: file.id,
            path: file.path,
            name: file.name,
            file_type: file.r#type,
        }
    }
}

impl From<models::BlogCategory> for BlogCategory {
    fn from(category: models::BlogCategory) -> Self {
        Self {
            id: category.id,
            name: category.name,
            description: category.description,
        }
    }
}

impl From<models::BlogPost> for BlogPost {
    fn from(blog_post: models::BlogPost) -> Self {
        Self {
            id: blog_post.id,
            title: blog_post.title,
            link: format!(
                "/{}/{}",
                blog_post.created_at.format("%Y/%m/%d"),
                blog_post.slug
            ),
            slug: blog_post.slug,
            header_image_id: blog_post.header_image_id,
            public: blog_post.public,
            creator_id: blog_post.creator_id,
            created_at: blog_post.created_at.to_rfc3339(),
            updated_by_id: blog_post.updated_by_id,
            last_updated_at: blog_post.last_updated_at.to_rfc3339(),
            blog_category_id: blog_post.category_id,
        }
    }
}

impl From<models::BlogPostSection> for BlogPostSection {
    fn from(value: models::blog_post_section::Model) -> Self {
        Self {
            gallery_id: value.gallery_id,
            file_id: value.file_id,
            html: value.html,
            link: value.link,
            position: value.position,
        }
    }
}

impl From<models::Gallery> for Gallery {
    fn from(gallery: models::Gallery) -> Self {
        Self {
            id: gallery.id,
            creator_id: gallery.creator_id,
            created_at: gallery.created_at.to_rfc3339(),
            updated_by_id: gallery.updated_by_id,
            last_updated_at: gallery.last_updated_at.to_rfc3339(),
            name: gallery.name,
            description: gallery.description,
            type_: if gallery.r#type == "masonry" {
                GalleryType::Masonry
            } else {
                GalleryType::Sequence
            },
            orientation: if gallery.r#type == "horizontal" {
                GalleryOrientation::Horizontal
            } else {
                GalleryOrientation::Vertical
            },
        }
    }
}

impl From<GalleryFileWithFile> for GalleryFile {
    fn from(value: GalleryFileWithFile) -> Self {
        Self {
            position: value.position,
            file: File {
                id: value.id,
                path: value.path,
                name: value.name,
                file_type: value.r#type,
            },
        }
    }
}

impl From<models::ClassicPage> for ClassicPage {
    fn from(value: models::classic_page::Model) -> Self {
        Self {
            creator_id: value.creator_id,
            created_at: value.created_at.to_rfc3339(),
            updated_by_id: value.updated_by_id,
            last_updated_at: value.last_updated_at.to_rfc3339(),
            content: value.content,
            title: value.title,
        }
    }
}

impl From<models::ModernPage> for ModernPage {
    fn from(value: models::modern_page::Model) -> Self {
        Self {
            id: value.id,
            creator_id: value.creator_id,
            created_at: value.created_at.to_rfc3339(),
            updated_by_id: value.updated_by_id,
            last_updated_at: value.last_updated_at.to_rfc3339(),
            name: value.name,
        }
    }
}

impl From<models::ModernPageSection> for ModernPageSection {
    fn from(value: models::modern_page_section::Model) -> Self {
        Self {
            position: value.position,
            gallery_id: value.gallery_id,
            file_id: value.file_id,
            html: value.html,
            action: value.action,
            target: value.target,
        }
    }
}

impl From<models::Menu> for Menu {
    fn from(value: models::menu::Model) -> Self {
        Self {
            id: value.id,
            name: value.name,
            logo: value.logo,
        }
    }
}

impl From<models::MenuItem> for MenuItem {
    fn from(value: models::menu_item::Model) -> Self {
        Self {
            id: value.id,
            menu_id: value.menu_id,
            parent_id: value.parent_id,
            position: value.position,
            title: value.title,
            route: value.route,
            highlighted: value.highlighted,
            artist_id: value.artist_id,
            gallery_id: value.gallery_id,
            modern_page_id: value.modern_page_id,
            classic_page_id: value.classic_page_id,
            blog_category_id: value.blog_category_id,
            blog_home_page: value.blog_home_page,
        }
    }
}

impl From<MenuItem> for models::MenuItem {
    fn from(val: MenuItem) -> Self {
        let route = val.route.map(|route| {
            if !route.starts_with("/") {
                format!("/{route}")
            } else {
                route
            }
        });

        models::MenuItem {
            id: val.id,
            position: val.position,
            menu_id: None,
            parent_id: None,
            title: val.title,
            route,
            highlighted: val.highlighted,
            artist_id: val.artist_id,
            gallery_id: val.gallery_id,
            modern_page_id: val.modern_page_id,
            classic_page_id: val.classic_page_id,
            blog_category_id: val.blog_category_id,
            blog_home_page: val.blog_home_page,
        }
    }
}

impl From<BoolOrString> for theme::BoolOrString {
    fn from(value: BoolOrString) -> Self {
        match value {
            BoolOrString::Bool(val) => theme::BoolOrString::BoolValue(val),
            BoolOrString::String(val) => theme::BoolOrString::StringValue(val),
        }
    }
}

impl From<BlogCategoryPath> for blog::BlogCategoryPathEntry {
    fn from(value: BlogCategoryPath) -> Self {
        Self {
            route: value.route,
            category_id: value.category_id,
            category_title: value.category_name,
            menu_item_id: value.menu_item_id,
            menu_item_title: value.menu_item_title,
        }
    }
}

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

impl db::Host for GojiThemeStore {
    async fn find_gallery_files(&mut self, gallery_id: i32) -> Vec<db::GalleryFile> {
        dbal::gallery::find_files(gallery_id)
            .await
            .unwrap_or_default()
            .into_iter()
            .map(Into::into)
            .collect()
    }

    async fn find_blog_post_sections(&mut self, blog_post_id: i32) -> Vec<db::BlogPostSection> {
        dbal::blog_post::find_sections(blog_post_id)
            .await
            .unwrap_or_default()
            .into_iter()
            .map(Into::into)
            .collect()
    }

    async fn find_menu_items_by_menu(&mut self, id: i32) -> Vec<db::MenuItem> {
        dbal::menu::find_menu_items(id)
            .await
            .unwrap_or_default()
            .into_iter()
            .map(Into::into)
            .collect()
    }

    async fn find_menu_items_by_menu_item(&mut self, id: i32) -> Vec<db::MenuItem> {
        dbal::menu_item::find_menu_items(id)
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
    ) -> Vec<db::BlogPost> {
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
    ) -> Vec<db::ModernPageSection> {
        dbal::modern_page::find_sections(modern_page_id)
            .await
            .unwrap_or_default()
            .into_iter()
            .map(Into::into)
            .collect()
    }

    async fn find_menu_item_by_route(&mut self, route: String) -> Option<db::MenuItem> {
        dbal::menu_item::find_menu_item_by_route(route)
            .await
            .map(Into::into)
    }

    async fn find_blog_category_by_id(&mut self, id: i32) -> Option<db::BlogCategory> {
        dbal::blog_category::find_by_id(id).await.map(Into::into)
    }

    async fn find_blog_post_by_id(&mut self, id: i32) -> Option<db::BlogPost> {
        dbal::blog_post::find_by_id(id).await.map(Into::into)
    }

    async fn find_classic_page_by_id(&mut self, id: i32) -> Option<db::ClassicPage> {
        dbal::classic_page::find_by_id(id).await.map(Into::into)
    }

    async fn find_file_by_id(&mut self, id: i32) -> Option<db::File> {
        dbal::file::find_by_id(id).await.map(Into::into)
    }

    async fn find_gallery_by_id(&mut self, id: i32) -> Option<db::Gallery> {
        dbal::gallery::find_by_id(id).await.map(Into::into)
    }

    async fn find_menu_by_id(&mut self, id: i32) -> Option<db::Menu> {
        dbal::menu::find_by_id(id).await.map(Into::into)
    }

    async fn find_modern_page_by_id(&mut self, id: i32) -> Option<db::ModernPage> {
        dbal::modern_page::find_by_id(id).await.map(Into::into)
    }

    async fn find_user_by_id(&mut self, id: i32) -> Option<db::User> {
        dbal::user::find_by_id(id).await.map(Into::into)
    }

    async fn find_all_blog_categories(&mut self) -> Vec<db::BlogCategory> {
        dbal::blog_category::find_all()
            .await
            .unwrap_or_default()
            .into_iter()
            .map(Into::into)
            .collect()
    }

    async fn find_all_blog_posts(&mut self) -> Vec<db::BlogPost> {
        dbal::blog_post::find_all()
            .await
            .unwrap_or_default()
            .into_iter()
            .map(Into::into)
            .collect()
    }

    async fn find_all_classic_pages(&mut self) -> Vec<db::ClassicPage> {
        dbal::classic_page::find_all()
            .await
            .unwrap_or_default()
            .into_iter()
            .map(Into::into)
            .collect()
    }

    async fn find_all_files(&mut self) -> Vec<db::File> {
        dbal::file::find_all()
            .await
            .unwrap_or_default()
            .into_iter()
            .map(Into::into)
            .collect()
    }

    async fn find_all_galleries(&mut self) -> Vec<db::Gallery> {
        dbal::gallery::find_all()
            .await
            .unwrap_or_default()
            .into_iter()
            .map(Into::into)
            .collect()
    }

    async fn find_all_menus(&mut self) -> Vec<db::Menu> {
        dbal::menu::find_all()
            .await
            .unwrap_or_default()
            .into_iter()
            .map(Into::into)
            .collect()
    }

    async fn find_all_modern_pages(&mut self) -> Vec<db::ModernPage> {
        dbal::modern_page::find_all()
            .await
            .unwrap_or_default()
            .into_iter()
            .map(Into::into)
            .collect()
    }

    async fn find_all_users(&mut self) -> Vec<db::User> {
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

impl theme::Host for GojiThemeStore {
    async fn get_theme_config_value(
        &mut self,
        group: String,
        key: String,
    ) -> Option<theme::BoolOrString> {
        self.theme.get_config(group, key).map(Into::into)
    }

    async fn get_active_menu_item(&mut self) -> Option<theme::MenuItem> {
        dbal::menu_item::get_item_by_path(self.path.clone())
            .await
            .map(Into::into)
    }

    async fn is_menu_item_active(&mut self, item: theme::MenuItem) -> bool {
        let route = item.route.clone().unwrap_or_default();
        let route = route.strip_prefix("/").unwrap_or(route.as_str());

        let path = self.path.strip_prefix("/").unwrap_or(self.path.as_str());

        route == path
    }

    async fn is_parent_of_active_menu_item(&mut self, item: theme::MenuItem) -> bool {
        dbal::menu_item::is_parent_of_active_menu_item(item.id, &self.path).await
    }

    async fn find_modern_page_by_theme(&mut self, name: String) -> Option<theme::ModernPage> {
        dbal::modern_page::find_by_theme(name, self.theme.db_theme.id)
            .await
            .map(Into::into)
    }

    async fn find_classic_page_by_theme(&mut self, name: String) -> Option<theme::ClassicPage> {
        dbal::classic_page::find_by_theme(name, self.theme.db_theme.id)
            .await
            .map(Into::into)
    }

    async fn find_file_by_theme(&mut self, name: String) -> Option<theme::File> {
        dbal::file::find_by_theme(name, self.theme.db_theme.id)
            .await
            .map(Into::into)
    }

    async fn find_menu_by_theme(&mut self, name: String) -> Option<theme::Menu> {
        dbal::menu::find_by_theme(name, self.theme.db_theme.id)
            .await
            .map(Into::into)
    }

    async fn find_gallery_by_theme(&mut self, name: String) -> Option<theme::Gallery> {
        dbal::gallery::find_by_theme(name, self.theme.db_theme.id)
            .await
            .map(Into::into)
    }

    async fn find_blog_category_by_theme(&mut self, name: String) -> Option<theme::BlogCategory> {
        dbal::blog_category::find_by_theme(name, self.theme.db_theme.id)
            .await
            .map(Into::into)
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

impl types::Host for GojiThemeStore {}

impl http::Host for GojiThemeStore {
    async fn get(&mut self, url: String) -> Result<String, String> {
        let res = reqwest::get(url).await;
        match res {
            Ok(res) => res.text().await.map_err(|err| err.to_string()),
            Err(err) => Err(err.to_string()),
        }
    }
}

impl request::Host for GojiThemeStore {
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

impl blog::Host for GojiThemeStore {
    async fn get_category_path_for_post(
        &mut self,
        post_id: i32,
    ) -> Vec<blog::BlogCategoryPathEntry> {
        dbal::blog_post::get_category_path(post_id)
            .await
            .into_iter()
            .map(Into::into)
            .collect()
    }

    async fn get_category_path_for_category(
        &mut self,
        category_id: i32,
    ) -> Vec<blog::BlogCategoryPathEntry> {
        dbal::blog_category::get_category_path(category_id)
            .await
            .into_iter()
            .map(Into::into)
            .collect()
    }
}

impl image::Host for GojiThemeStore {
    async fn image_size(&mut self, file_id: i32) -> (i32, i32) {
        if let Some(file) = dbal::file::find_by_id(file_id).await {
            if let Ok(size) = imagesize::size(
                Path::new(&std::env::var("STORAGE_PATH").unwrap_or("data".to_string()))
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
