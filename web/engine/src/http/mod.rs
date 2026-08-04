use crate::dbal;
use crate::models::{BlogCategory, BlogPost, ClassicPage, Gallery, ModernPage, User};
use crate::theming::engine::{GojiTheme, GojiThemeStore};
use crate::theming::theme::Theme;
use actix_web::body::BoxBody;
use actix_web::http::StatusCode;
use actix_web::web::Query;
use actix_web::{HttpRequest, HttpResponse};
use serde::Deserialize;
use std::path::Path;
use url::Url;
use wasmtime::component::{Component, HasSelf, Linker, ResourceTable};
use wasmtime::{Engine, Store};
use wasmtime_wasi::WasiCtx;

#[derive(Debug)]
enum RequestEntity {
    BlogCategory(BlogCategory),
    BlogPost(BlogPost),
    BlogHomePage,
    HomePage,
    Gallery(Gallery),
    ModernPage(ModernPage),
    ClassicPage(ClassicPage),
    Profile(User),
    Error(StatusCode),
}

#[derive(Deserialize)]
pub struct ErrorQuery {
    pub error: u16,
}

pub async fn error_route(req: HttpRequest, query: Query<ErrorQuery>) -> HttpResponse {
    let is_json_request = req.headers().get("Accept").is_some_and(|accept| {
        accept
            .to_str()
            .is_ok_and(|accept| accept.starts_with("application/json"))
    });

    let (status_code, content) =
        handle_error(&req.uri().to_string(), is_json_request, query.error).await;

    HttpResponse::new(status_code).set_body(BoxBody::new(content))
}

pub async fn frontend_router(req: HttpRequest) -> HttpResponse {
    let is_json_request = req.headers().get("Accept").is_some_and(|accept| {
        accept
            .to_str()
            .is_ok_and(|accept| accept.starts_with("application/json"))
    });

    let (status_code, content) = handle_request(&req.uri().to_string(), is_json_request).await;

    HttpResponse::new(status_code).set_body(BoxBody::new(content))
}

async fn get_request_entity(path: &str) -> Option<RequestEntity> {
    let menu_item = dbal::menu_item::get_item_by_path(path.to_string()).await?;

    if let Some(id) = menu_item.artist_id
        && let Some(user) = dbal::user::find_by_id(id).await
    {
        Some(RequestEntity::Profile(user))
    } else if let Some(id) = menu_item.gallery_id
        && let Some(gallery) = dbal::gallery::find_by_id(id).await
    {
        Some(RequestEntity::Gallery(gallery))
    } else if let Some(id) = menu_item.modern_page_id
        && let Some(modern_page) = dbal::modern_page::find_by_id(id).await
    {
        Some(RequestEntity::ModernPage(modern_page))
    } else if let Some(id) = menu_item.classic_page_id
        && let Some(classic_page) = dbal::classic_page::find_by_id(id).await
    {
        Some(RequestEntity::ClassicPage(classic_page))
    } else if let Some(id) = menu_item.blog_category_id
        && let Some(blog_category) = dbal::blog_category::find_by_id(id).await
    {
        Some(RequestEntity::BlogCategory(blog_category))
    } else if let Some(is_homepage) = menu_item.blog_home_page
        && is_homepage
    {
        Some(RequestEntity::BlogHomePage)
    } else {
        None
    }
}

fn parse_blog_path(path: &str) -> Option<String> {
    let mut parts = path.trim_matches('/').split('/');

    let year = parts.next()?.parse::<u16>().ok()?;
    let month = parts.next()?.parse::<u8>().ok()?;
    let day = parts.next()?.parse::<u8>().ok()?;
    let slug = parts.next()?.to_string();

    if parts.next().is_some() {
        return None;
    }

    if !(1970..=9999).contains(&year) {
        return None;
    }

    if !(1..=12).contains(&month) {
        return None;
    }

    if !(1..=31).contains(&day) {
        return None;
    }

    Some(slug)
}

async fn get_blog_post(path: &str) -> Option<BlogPost> {
    let slug = parse_blog_path(path)?;

    dbal::blog_post::find_by_slug(&slug).await
}

pub async fn handle_request(uri: &str, is_api: bool) -> (StatusCode, String) {
    let theme = match dbal::theme::get_active_theme().await {
        Ok(theme) => theme,
        Err(err) => {
            log::error!("{err}");
            return (StatusCode::INTERNAL_SERVER_ERROR, "".to_string());
        }
    };
    let theme = if let Ok(theme) = Theme::new(theme).await {
        theme
    } else {
        return (StatusCode::INTERNAL_SERVER_ERROR, "".to_string());
    };

    let uri = if std::env::var("ENV").unwrap_or("prod".to_string()) == "dev" {
        format!("http://localhost/{uri}")
    } else {
        uri.to_string()
    };

    let url = Url::parse(&uri);
    if url.is_err() {
        log::error!("Failed to parse url");
        return (StatusCode::INTERNAL_SERVER_ERROR, "".to_string());
    }

    let url = url.unwrap();
    let request_entity = if url.path() == "" || url.path() == "/" {
        RequestEntity::HomePage
    } else if let Some(post) = get_blog_post(url.path()).await {
        RequestEntity::BlogPost(post)
    } else if let Some(request_entity) = get_request_entity(url.path()).await {
        request_entity
    } else {
        RequestEntity::Error(StatusCode::NOT_FOUND)
    };

    call_theme(&uri, is_api, theme, url, request_entity).await
}

pub async fn handle_error(uri: &str, is_api: bool, status_code: u16) -> (StatusCode, String) {
    let theme = if let Ok(theme) = dbal::theme::get_active_theme().await {
        theme
    } else {
        return (StatusCode::INTERNAL_SERVER_ERROR, "".to_string());
    };
    let theme = if let Ok(theme) = Theme::new(theme).await {
        theme
    } else {
        return (StatusCode::INTERNAL_SERVER_ERROR, "".to_string());
    };

    let uri = if std::env::var("ENV").unwrap_or("prod".to_string()) == "dev" {
        format!("http://localhost/{uri}")
    } else {
        uri.to_string()
    };

    let url = Url::parse(&uri);
    if url.is_err() {
        log::error!("Failed to parse url");
        return (StatusCode::INTERNAL_SERVER_ERROR, "".to_string());
    }

    let url = url.unwrap();

    call_theme(
        &uri,
        is_api,
        theme,
        url,
        RequestEntity::Error(
            StatusCode::from_u16(status_code).unwrap_or(StatusCode::INTERNAL_SERVER_ERROR),
        ),
    )
    .await
}

async fn call_theme(
    uri: &str,
    is_api: bool,
    theme: &Theme,
    url: Url,
    request_entity: RequestEntity,
) -> (StatusCode, String) {
    let engine = Engine::default();
    let mut linker = Linker::new(&engine);
    if let Err(err) = wasmtime_wasi::p2::add_to_linker_async(&mut linker) {
        log::error!(
            "Failed to append wasip2 for theme {}: {err}",
            theme.db_theme.name
        );
        return (StatusCode::INTERNAL_SERVER_ERROR, "".to_string());
    }

    let wasi_ctx = WasiCtx::builder().inherit_stdio().inherit_args().build();

    let mut store = Store::new(
        &engine,
        GojiThemeStore {
            theme: theme.clone(),
            uri: uri.to_string(),
            path: url.path().to_string(),
            is_api_request: is_api && theme.has_api,
            wasi_ctx,
            resource_table: ResourceTable::new(),
        },
    );

    let theme_path = Path::new(&theme.db_theme.path)
        .join("theme.wasm")
        .to_string_lossy()
        .to_string();

    let component = match Component::from_file(&engine, &theme_path) {
        Ok(component) => component,
        Err(err) => {
            log::error!(
                "Failed to build component for theme {}: {err:#?}",
                theme.db_theme.name
            );
            return (StatusCode::INTERNAL_SERVER_ERROR, "".to_string());
        }
    };

    if let Err(err) =
        GojiTheme::add_to_linker::<_, HasSelf<GojiThemeStore>>(&mut linker, |state| state)
    {
        log::error!(
            "Failed to add goji theme {} to linker: {err}",
            theme.db_theme.name
        );
        return (StatusCode::INTERNAL_SERVER_ERROR, "".to_string());
    }

    let goji_theme = match GojiTheme::instantiate_async(&mut store, &component, &linker).await {
        Ok(theme) => theme,
        Err(err) => {
            log::error!(
                "Failed to instantiate goji theme {}: {err}",
                theme.db_theme.name
            );
            return (StatusCode::INTERNAL_SERVER_ERROR, "".to_string());
        }
    };

    let response = match request_entity {
        RequestEntity::BlogCategory(cat) => {
            goji_theme
                .func_handle_blog_category()
                .call_async(&mut store, (cat.into(),))
                .await
        }
        RequestEntity::BlogPost(post) => {
            goji_theme
                .func_handle_blog_post()
                .call_async(&mut store, (post.into(),))
                .await
        }
        RequestEntity::BlogHomePage => {
            goji_theme
                .func_handle_blog_homepage()
                .call_async(&mut store, ())
                .await
        }
        RequestEntity::HomePage => {
            goji_theme
                .func_handle_homepage()
                .call_async(&mut store, ())
                .await
        }
        RequestEntity::Gallery(gallery) => {
            goji_theme
                .func_handle_gallery()
                .call_async(&mut store, (gallery.into(),))
                .await
        }
        RequestEntity::ModernPage(page) => {
            goji_theme
                .func_handle_modern_page()
                .call_async(&mut store, (page.into(),))
                .await
        }
        RequestEntity::ClassicPage(page) => {
            goji_theme
                .func_handle_classic_page()
                .call_async(&mut store, (page.into(),))
                .await
        }
        RequestEntity::Profile(user) => {
            goji_theme
                .func_handle_profile()
                .call_async(&mut store, (user.into(),))
                .await
        }
        RequestEntity::Error(status) => {
            goji_theme
                .func_handle_error()
                .call_async(&mut store, (status.as_u16(),))
                .await
        }
    };
    match response {
        Ok((response, ..)) => (
            StatusCode::from_u16(response.status_code).unwrap_or(StatusCode::INTERNAL_SERVER_ERROR),
            response.content,
        ),
        Err(err) => {
            log::error!(
                "Failed to rendering goji theme {}: {err}",
                theme.db_theme.name
            );
            (StatusCode::INTERNAL_SERVER_ERROR, "".to_string())
        }
    }
}
