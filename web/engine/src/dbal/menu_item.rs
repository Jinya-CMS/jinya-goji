use crate::dbal::db;
use crate::models::MenuItem;
use crate::{find_by_parent, find_entity, models};
use sea_orm::{ColumnTrait, Condition, EntityTrait, QueryFilter};

find_entity!(menu, Menu);

find_by_parent!(menu_item, MenuItem, ParentId, find_menu_items);

pub async fn get_item_by_path(path: String) -> Option<MenuItem> {
    models::menu_item::Entity::find()
        .filter(
            Condition::any()
                .add(models::menu_item::Column::Route.eq(path.strip_prefix("/")))
                .add(models::menu_item::Column::Route.eq(path)),
        )
        .one(db().await.ok()?)
        .await
        .ok()?
}

pub async fn is_parent_of_active_menu_item(item_id: i32, item_path: &str, path: &str) -> bool {
    let normalized_path = path.strip_prefix("/").unwrap_or(path);

    if item_path.strip_prefix("/").unwrap_or(item_path) == normalized_path {
        return true;
    }

    let mut stack = vec![item_id];

    while let Some(parent_id) = stack.pop() {
        let children = find_menu_items(parent_id).await.unwrap_or_default();

        for item in children {
            let item_route = item.route.unwrap_or_default();

            if item_route.strip_prefix("/").unwrap_or(&item_route) == normalized_path {
                return true;
            }

            stack.push(item.id);
        }
    }

    false
}

pub async fn find_menu_item_by_route(route: String) -> Option<MenuItem> {
    models::menu_item::Entity::find()
        .filter(
            Condition::any()
                .add(models::menu_item::Column::Route.eq(&route))
                .add(
                    models::menu_item::Column::Route.eq(route.strip_prefix("/").unwrap_or(&route)),
                ),
        )
        .one(db().await.ok()?)
        .await
        .unwrap_or_default()
}
