use crate::models::gallery_file_position::GalleryFileWithFile;
use crate::models::theme::BoolOrString;
use crate::theming::engine::jinya::goji::{db, http, request, theme, types};
use crate::theming::theme::Theme;
use crate::{dbal, models};
use std::collections::HashMap;
use wasmtime::component::{bindgen, HasData, ResourceTable};
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
        }
    }
}

impl From<models::BlogCategory> for BlogCategory {
    fn from(category: models::BlogCategory) -> Self {
        Self {
            name: category.name,
            description: category.description,
        }
    }
}

impl From<models::BlogPost> for BlogPost {
    fn from(blog_post: models::BlogPost) -> Self {
        Self {
            title: blog_post.title,
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
            orientation: if gallery.r#type == "masonry" {
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
        models::MenuItem {
            id: val.id,
            position: val.position,
            menu_id: None,
            parent_id: None,
            title: val.title,
            route: val.route,
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
        item.route
            .unwrap_or_default()
            .strip_prefix("/")
            .unwrap_or_default()
            == self.path.strip_prefix("/").unwrap_or_default()
    }

    async fn is_parent_of_active_menu_item(&mut self, item: theme::MenuItem) -> bool {
        dbal::menu_item::is_parent_of_active_menu_item(
            item.id,
            &item.route.unwrap_or_default(),
            &self.path,
        )
        .await
    }

    async fn find_modern_page_by_theme(&mut self, name: String) -> Option<theme::ModernPage> {
        dbal::modern_page::find_by_theme(name).await.map(Into::into)
    }

    async fn find_classic_page_by_theme(&mut self, name: String) -> Option<theme::ClassicPage> {
        dbal::classic_page::find_by_theme(name)
            .await
            .map(Into::into)
    }

    async fn find_file_by_theme(&mut self, name: String) -> Option<theme::File> {
        dbal::file::find_by_theme(name).await.map(Into::into)
    }

    async fn find_menu_by_theme(&mut self, name: String) -> Option<theme::Menu> {
        dbal::menu::find_by_theme(name).await.map(Into::into)
    }

    async fn find_gallery_by_theme(&mut self, name: String) -> Option<theme::Gallery> {
        dbal::gallery::find_by_theme(name).await.map(Into::into)
    }

    async fn find_blog_category_by_theme(&mut self, name: String) -> Option<theme::BlogCategory> {
        dbal::blog_category::find_by_theme(name)
            .await
            .map(Into::into)
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
