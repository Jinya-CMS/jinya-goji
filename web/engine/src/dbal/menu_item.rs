use crate::dbal::db;
use crate::models::MenuItem;
use crate::{find_by_parent, find_entity, models};
use sea_orm::{
    ColumnTrait, Condition, ConnectionTrait, DbBackend, EntityTrait, QueryFilter, Statement,
};

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

pub async fn is_parent_of_active_menu_item(item_id: i32, path: &str) -> bool {
    //language=postgresql
    let query = r#"
with recursive menu_item_path as (select mi.*, 0 as depth
                                  from menu_item mi
                                  where mi.route = $1
                                     or mi.route = $2

                                  union all

                                  -- recursive step: walk up to the parent
                                  select mi2.*, mip.depth + 1
                                  from menu_item_path mip
                                           inner join menu_item mi2 on mi2.id = mip.parent_id)
select *
from menu_item_path
where id = $3
limit 1
    "#;
    if let Ok(db) = db().await {
        db.query_one_raw(Statement::from_sql_and_values(
            DbBackend::Postgres,
            query,
            vec![
                sea_orm::Value::String(Some(path.to_string())),
                sea_orm::Value::String(Some(path.strip_prefix("/").unwrap_or(path).to_string())),
                sea_orm::Value::Int(Some(item_id)),
            ],
        ))
        .await
        .ok()
        .unwrap_or_default()
        .is_some()
    } else {
        false
    }
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
