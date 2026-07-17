use crate::dbal::db;
use crate::models::{BlogCategory, BlogPost, ClassicPage, Gallery, ModernPage, User};
use crate::theming::theme::Theme;
use crate::{dbal, models};
use actix_web::http::StatusCode;
use actix_web::{HttpRequest, HttpResponse};
use sea_orm::{ColumnTrait, Condition, EntityTrait, QueryFilter};

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

pub async fn frontend_router(req: HttpRequest) -> HttpResponse {
    let theme = if let Ok(theme) = dbal::theme::get_active_theme().await {
        theme
    } else {
        return HttpResponse::InternalServerError().finish();
    };
    let theme = if let Ok(theme) = Theme::new(theme).await {
        theme
    } else {
        return HttpResponse::InternalServerError().finish();
    };

    let is_json_request = req.headers().get("Accept").is_some_and(|accept| {
        accept
            .to_str()
            .is_ok_and(|accept| accept.starts_with("application/json"))
    });

    let request_entity = if req.path() == "" || req.path() == "/" {
        RequestEntity::HomePage
    } else if let Some(post) = get_blog_post(req.path()).await {
        RequestEntity::BlogPost(post)
    } else if let Some(request_entity) = get_request_entity(req.path()).await {
        request_entity
    } else {
        RequestEntity::Error(StatusCode::NOT_FOUND)
    };

    if theme.has_api && is_json_request {
        handle_json_request(request_entity, theme).await
    } else {
        handle_standard_request(request_entity, theme).await
    }
}

async fn get_request_entity(path: &str) -> Option<RequestEntity> {
    let menu_item = models::menu_item::Entity::find()
        .filter(
            Condition::any()
                .add(models::menu_item::Column::Route.eq(path))
                .add(models::menu_item::Column::Route.eq(path.strip_prefix("/"))),
        )
        .one(db().await.ok()?)
        .await
        .ok()??;

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
    } else if let Some(is_homepage) = menu_item.blog_home_page && is_homepage {
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

async fn handle_json_request(request_entity: RequestEntity, theme: &Theme) -> HttpResponse {
    HttpResponse::Ok().finish()
}

async fn handle_standard_request(request_entity: RequestEntity, theme: &Theme) -> HttpResponse {
    HttpResponse::Ok().finish()
}
