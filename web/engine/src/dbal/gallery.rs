use crate::dbal::db;
use crate::models::gallery_file_position::GalleryFileWithFile;
use crate::{find_by_theme, find_entity, models};
use sea_orm::{ColumnTrait, EntityTrait, JoinType, QueryFilter, QuerySelect, RelationTrait};

find_entity!(gallery, Gallery);

find_by_theme!(gallery, Gallery, theme_gallery, ThemeGallery, find_by_theme);

pub async fn find_files(gallery_id: i32) -> Option<Vec<GalleryFileWithFile>> {
    models::gallery_file_position::Entity::find()
        .join(
            JoinType::InnerJoin,
            models::gallery_file_position::Relation::File.def(),
        )
        .filter(models::gallery_file_position::Column::GalleryId.eq(gallery_id))
        .select_only()
        .column(models::gallery_file_position::Column::Position)
        .column(models::file::Column::Path)
        .column(models::file::Column::Name)
        .column(models::file::Column::Id)
        .into_model::<GalleryFileWithFile>()
        .all(db().await.ok()?)
        .await
        .ok()
}
